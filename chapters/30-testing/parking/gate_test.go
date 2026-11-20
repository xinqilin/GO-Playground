package parking

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

// fakePayments 是測試用的金流：記下每一筆收費；err 不是 nil 時，每一筆都失敗。
type fakePayments struct {
	mu      sync.Mutex
	charges []string
	err     error
}

func (f *fakePayments) Charge(ctx context.Context, plate string, amount int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.err != nil {
		return f.err
	}
	f.charges = append(f.charges, fmt.Sprintf("%s %d", plate, amount))
	return nil
}

// fakeClock 是可以手動撥快的時鐘。
type fakeClock struct{ now time.Time }

func (c *fakeClock) Now() time.Time          { return c.now }
func (c *fakeClock) Advance(d time.Duration) { c.now = c.now.Add(d) }

// newTestGate 建立一個用假時鐘、假金流的柵欄機；測試失敗時順便印出出場紀錄。
func newTestGate(t *testing.T) (*Gate, *fakeClock, *fakePayments) {
	t.Helper()
	clock := &fakeClock{now: at(20, 9, 0)}
	pay := &fakePayments{}
	g := NewGate(clock.Now, pay)
	t.Cleanup(func() {
		if t.Failed() {
			t.Logf("出場紀錄：%+v", g.Records())
		}
	})
	return g, clock, pay
}

// mustEnter 進場失敗就停止這個測試。t.Helper 讓失敗訊息指向呼叫它的那一行。
func mustEnter(t *testing.T, g *Gate, plate string) {
	t.Helper()
	if err := g.Enter(plate); err != nil {
		t.Fatalf("Enter(%q)：%v", plate, err)
	}
}

func TestGateExit(t *testing.T) {
	g, clock, pay := newTestGate(t)
	mustEnter(t, g, "abc 1234")
	clock.Advance(95 * time.Minute)

	rec, err := g.Exit(t.Context(), "ABC-1234")
	if err != nil {
		t.Fatalf("Exit：%v", err)
	}
	if rec.Fee != 80 {
		t.Errorf("Fee = %d，want 80", rec.Fee)
	}
	if want := []string{"ABC-1234 80"}; !slices.Equal(pay.charges, want) {
		t.Errorf("收費紀錄 = %q，want %q", pay.charges, want)
	}
}

func TestGateFreeExitDoesNotCharge(t *testing.T) {
	g, clock, pay := newTestGate(t)
	mustEnter(t, g, "1234-AB")
	clock.Advance(5 * time.Minute)
	if _, err := g.Exit(t.Context(), "1234-AB"); err != nil {
		t.Fatalf("Exit：%v", err)
	}
	if len(pay.charges) != 0 {
		t.Errorf("免費時段不應該收費，收了 %q", pay.charges)
	}
}

func TestGatePaymentFails(t *testing.T) {
	g, clock, pay := newTestGate(t)
	pay.err = errors.New("刷卡機沒有回應")
	mustEnter(t, g, "ABC-1234")
	clock.Advance(time.Hour)

	if _, err := g.Exit(t.Context(), "ABC-1234"); !errors.Is(err, pay.err) {
		t.Fatalf("err = %v，want %v", err, pay.err)
	}
	pay.err = nil // 刷卡機好了，車還在場內，可以再試一次
	if rec, err := g.Exit(t.Context(), "ABC-1234"); err != nil || rec.Fee != 40 {
		t.Errorf("再試一次 = %+v, %v，want 40 元", rec, err)
	}
}

// slowPayments 一直等到 ctx 結束才回應，模擬卡住的金流。
type slowPayments struct{}

func (slowPayments) Charge(ctx context.Context, plate string, amount int) error {
	<-ctx.Done()
	return ctx.Err()
}

// TestGateChargeTimeout：在 synctest 的泡泡裡，時間是假的，等 1 小時加 5 秒也是一瞬間。
func TestGateChargeTimeout(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		g := NewGate(time.Now, slowPayments{})
		mustEnter(t, g, "ABC-1234")
		time.Sleep(time.Hour)

		start := time.Now()
		_, err := g.Exit(t.Context(), "ABC-1234")
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("err = %v，want DeadlineExceeded", err)
		}
		if waited := time.Since(start); waited != ChargeTimeout {
			t.Errorf("等了 %v，want %v", waited, ChargeTimeout)
		}
		if err := g.Enter("ABC-1234"); !errors.Is(err, ErrInside) {
			t.Errorf("收費逾時後車應該還在場內，Enter 回傳 %v", err)
		}
	})
}

// reenterPayments 在收費途中讓同一個車牌再進場一次，然後收費失敗。
type reenterPayments struct {
	g     *Gate
	clock *fakeClock
}

func (r *reenterPayments) Charge(ctx context.Context, plate string, amount int) error {
	r.clock.Advance(time.Minute)
	if err := r.g.Enter(plate); err != nil {
		return err
	}
	return errors.New("刷卡失敗")
}

// TestGatePaymentFailsAfterReenter：收費失敗時把車放回場內，不能蓋掉等金流的這段時間裡新的進場時間。
func TestGatePaymentFailsAfterReenter(t *testing.T) {
	clock := &fakeClock{now: at(20, 9, 0)}
	pay := &reenterPayments{clock: clock}
	g := NewGate(clock.Now, pay)
	pay.g = g
	mustEnter(t, g, "ABC-1234")
	clock.Advance(time.Hour)

	if _, err := g.Exit(t.Context(), "ABC-1234"); err == nil {
		t.Fatal("收費失敗，Exit 應該回傳錯誤")
	}
	if in, want := g.inside["ABC-1234"], at(20, 10, 1); !in.Equal(want) {
		t.Errorf("場內的進場時間 = %v，want 新的進場時間 %v", in, want)
	}
}

// TestGateConcurrentExit：同一台車同時出場兩次，只能收一次錢。用 go test -race 跑。
func TestGateConcurrentExit(t *testing.T) {
	g, clock, pay := newTestGate(t)
	mustEnter(t, g, "ABC-1234")
	clock.Advance(time.Hour)

	var wg sync.WaitGroup
	errs := make([]error, 2)
	for i := range errs {
		wg.Go(func() {
			_, errs[i] = g.Exit(t.Context(), "ABC-1234")
		})
	}
	wg.Wait()
	if n := len(pay.charges); n != 1 {
		t.Errorf("收了 %d 次錢，want 1", n)
	}
	if !errors.Is(errs[0], ErrNotInside) && !errors.Is(errs[1], ErrNotInside) {
		t.Errorf("其中一次應該是 ErrNotInside，拿到 %v", errs)
	}
}

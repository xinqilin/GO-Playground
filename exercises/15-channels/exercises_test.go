package relay

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"testing/synctest"
	"time"
)

// recv 收 ch 直到它被 close。等了一秒還沒 close，就回傳 ok 為 false，避免測試永遠卡住。
func recv[T any](ch <-chan T) (got []T, ok bool) {
	timeout := time.After(time.Second)
	for {
		select {
		case v, open := <-ch:
			if !open {
				return got, true
			}
			got = append(got, v)
		case <-timeout:
			return got, false
		}
	}
}

// feed 回傳一個已經放好 vals、而且 close 的 channel。
func feed[T any](vals ...T) <-chan T {
	ch := make(chan T, len(vals))
	for _, v := range vals {
		ch <- v
	}
	close(ch)
	return ch
}

func TestTickets(t *testing.T) {
	got, ok := recv(Tickets(17, 20))
	if !ok {
		t.Fatalf("Tickets(17, 20) 的 channel 一直沒有 close，收到 %v", got)
	}
	if !slices.Equal(got, []int{17, 18, 19, 20}) {
		t.Errorf("Tickets(17, 20) = %v, want [17 18 19 20]", got)
	}
	if got, ok := recv(Tickets(5, 5)); !ok || !slices.Equal(got, []int{5}) {
		t.Errorf("Tickets(5, 5) = %v（close 了：%v）, want [5]", got, ok)
	}
	if got, ok := recv(Tickets(3, 1)); !ok || len(got) != 0 {
		t.Errorf("Tickets(3, 1) = %v（close 了：%v）, want 空的而且 close", got, ok)
	}
}

func TestCollect(t *testing.T) {
	got := Collect(feed("排骨飯", "雞腿飯", "魚排飯"))
	if want := []string{"排骨飯", "雞腿飯", "魚排飯"}; !slices.Equal(got, want) {
		t.Errorf("Collect = %q, want %q", got, want)
	}
	if got := Collect(feed[int]()); len(got) != 0 {
		t.Errorf("Collect(空的) = %v", got)
	}
}

func TestStage(t *testing.T) {
	got, ok := recv(Stage(feed(85, 95, 30), func(p int) string {
		return strings.Repeat("$", p/30)
	}))
	if !ok {
		t.Fatalf("Stage 的 channel 一直沒有 close，收到 %q", got)
	}
	if want := []string{"$$", "$$$", "$"}; !slices.Equal(got, want) {
		t.Errorf("Stage = %q, want %q", got, want)
	}
}

func TestFanIn(t *testing.T) {
	got, ok := recv(FanIn(feed(1, 2, 3), feed[int](), feed(4, 5)))
	if !ok {
		t.Fatalf("FanIn 的 channel 一直沒有 close，收到 %v", got)
	}
	slices.Sort(got)
	if !slices.Equal(got, []int{1, 2, 3, 4, 5}) {
		t.Errorf("FanIn 收到 %v（排序後）, want [1 2 3 4 5]", got)
	}
	if got, ok := recv(FanIn[int]()); !ok || len(got) != 0 {
		t.Errorf("沒有輸入時 FanIn = %v（close 了：%v）, want 空的而且 close", got, ok)
	}
}

var errBusy = errors.New("忙線中")

// TestFirstOK 在 synctest 的泡泡裡執行，time.Sleep 用的是泡泡裡的假時間，所以誰先完成是固定的；
// 泡泡結束時如果還有 goroutine 卡著，就代表有洩漏。
func TestFirstOK(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		slow := func(d time.Duration, v string, err error) func() (string, error) {
			return func() (string, error) {
				time.Sleep(d)
				return v, err
			}
		}

		start := time.Now()
		got, err := FirstOK(
			slow(3*time.Second, "熊貓", nil),
			slow(1*time.Second, "", errBusy),
			slow(2*time.Second, "優食", nil),
		)
		if got != "優食" || err != nil {
			t.Errorf("FirstOK = %q, %v; want 優食, nil", got, err)
		}
		if waited := time.Since(start); waited != 2*time.Second {
			t.Errorf("FirstOK 等了 %v，應該在第一個成功時（2s）就回傳", waited)
		}

		_, err = FirstOK(
			slow(2*time.Second, "", errors.New("熊貓：不在外送範圍")),
			slow(1*time.Second, "", errBusy),
		)
		if err == nil || err.Error() != "熊貓：不在外送範圍\n忙線中" {
			t.Errorf("全部失敗時 err = %q, want 兩個錯誤依照 tries 的順序", err)
		}
		if !errors.Is(err, errBusy) {
			t.Errorf("errors.Is(err, errBusy) = false")
		}
	})
}

func TestReceipts(t *testing.T) {
	got, ok := recv(Receipts([]int{85, 95, 30}))
	if want := []string{"NT$85", "NT$95", "NT$30"}; !slices.Equal(got, want) {
		t.Errorf("Receipts 收到 %q, want %q", got, want)
	}
	if !ok {
		t.Errorf("Receipts 的 channel 一直沒有 close：range 它的人會永遠等下去")
	}
}

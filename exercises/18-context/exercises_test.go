package courier

import (
	"context"
	"errors"
	"slices"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

func TestSleep(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		start := time.Now()
		if err := Sleep(context.Background(), 3*time.Second); err != nil || time.Since(start) != 3*time.Second {
			t.Errorf("沒有取消：Sleep = %v，等了 %v; want nil，3s", err, time.Since(start))
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		start = time.Now()
		if err := Sleep(ctx, 3*time.Second); !errors.Is(err, context.DeadlineExceeded) || time.Since(start) != time.Second {
			t.Errorf("1 秒後逾時：Sleep = %v，等了 %v; want DeadlineExceeded，1s", err, time.Since(start))
		}

		done, cancel2 := context.WithCancel(context.Background())
		cancel2()
		start = time.Now()
		if err := Sleep(done, 3*time.Second); !errors.Is(err, context.Canceled) || time.Since(start) != 0 {
			t.Errorf("已經取消：Sleep = %v，等了 %v; want Canceled，0s", err, time.Since(start))
		}
	})
}

func TestShouldStart(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		if !ShouldStart(context.Background(), time.Hour) {
			t.Errorf("沒有期限的 ctx 應該回傳 true")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if !ShouldStart(ctx, 3*time.Second) {
			t.Errorf("剩 5 秒、需要 3 秒，應該回傳 true")
		}
		time.Sleep(3 * time.Second)
		if ShouldStart(ctx, 3*time.Second) {
			t.Errorf("剩 2 秒、需要 3 秒，應該回傳 false")
		}
		if !ShouldStart(ctx, 2*time.Second) {
			t.Errorf("剩 2 秒、需要 2 秒，應該回傳 true")
		}

		canceled, cancel2 := context.WithCancel(context.Background())
		cancel2()
		if ShouldStart(canceled, 0) {
			t.Errorf("已經取消的 ctx 應該回傳 false")
		}
	})
}

func TestRoute(t *testing.T) {
	bg := context.Background()
	if r := Route(bg); r != nil {
		t.Fatalf("沒有任何一站：Route = %q, want nil", r)
	}
	shop := WithStop(bg, "總店")
	mrt := WithStop(shop, "捷運站")
	park := WithStop(mrt, "公園")
	school := WithStop(park, "學校")
	hospital := WithStop(park, "醫院")

	for _, tc := range []struct {
		ctx  context.Context
		want []string
	}{
		{shop, []string{"總店"}},
		{park, []string{"總店", "捷運站", "公園"}},
		{school, []string{"總店", "捷運站", "公園", "學校"}},
		{hospital, []string{"總店", "捷運站", "公園", "醫院"}},
	} {
		if got := Route(tc.ctx); !slices.Equal(got, tc.want) {
			t.Errorf("Route = %q, want %q", got, tc.want)
		}
	}

	if r := Route(park); len(r) > 0 {
		r[0] = "改掉了"
		if got := Route(park); len(got) == 0 || got[0] != "總店" {
			t.Errorf("修改 Route 回傳的 slice 之後，ctx 裡的路線變成 %q", got)
		}
	}
}

func TestPump(t *testing.T) {
	in := make(chan int, 3)
	in <- 1
	in <- 2
	in <- 3
	close(in)
	out := make(chan int, 3)
	if err := Pump(context.Background(), in, out); err != nil {
		t.Fatalf("in 被 close：Pump = %v, want nil", err)
	}
	close(out)
	var got []int
	for v := range out {
		got = append(got, v)
	}
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Errorf("轉送了 %v, want [1 2 3]", got)
	}

	synctest.Test(t, func(t *testing.T) {
		// 等收的時候被取消：in 一直沒有東西。
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		start := time.Now()
		err := Pump(ctx, make(chan int), make(chan int))
		if !errors.Is(err, context.DeadlineExceeded) || time.Since(start) != 2*time.Second {
			t.Errorf("等收時逾時：Pump = %v，花了 %v; want DeadlineExceeded，2s", err, time.Since(start))
		}

		// 等送的時候被取消：out 一直沒有人收。
		in := make(chan int, 1)
		in <- 7
		ctx2, cancel2 := context.WithCancel(context.Background())
		time.AfterFunc(time.Second, cancel2)
		start = time.Now()
		err = Pump(ctx2, in, make(chan int))
		if !errors.Is(err, context.Canceled) || time.Since(start) != time.Second {
			t.Errorf("等送時取消：Pump = %v，花了 %v; want Canceled，1s", err, time.Since(start))
		}
	})
}

// kitchen 回傳一間 d 之後做好 dish 的廚房；被取消時把 stopped 加一並回傳 ctx.Err()。
func kitchen(d time.Duration, dish string, stopped *atomic.Int64) func(context.Context) (string, error) {
	return func(ctx context.Context) (string, error) {
		select {
		case <-time.After(d):
			return dish, nil
		case <-ctx.Done():
			stopped.Add(1)
			return "", ctx.Err()
		}
	}
}

func TestFastest(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var stopped atomic.Int64
		errBroken := errors.New("瓦斯沒了")
		broken := func(ctx context.Context) (string, error) {
			time.Sleep(time.Second)
			return "", errBroken
		}

		start := time.Now()
		got, err := Fastest(context.Background(),
			kitchen(5*time.Second, "西區 排骨飯", &stopped),
			broken,
			kitchen(3*time.Second, "東區 排骨飯", &stopped),
			kitchen(time.Hour, "北區 排骨飯", &stopped),
		)
		if got != "東區 排骨飯" || err != nil || time.Since(start) != 3*time.Second {
			t.Fatalf("Fastest = %q, %v，花了 %v; want 東區 排骨飯, nil，3s", got, err, time.Since(start))
		}
		if n := stopped.Load(); n != 2 {
			t.Errorf("Fastest 回傳時，有 %d 間廚房收到取消, want 2（西區、北區）", n)
		}

		errClosed := errors.New("今天公休")
		closed := func(ctx context.Context) (string, error) { return "", errClosed }
		_, err = Fastest(context.Background(), broken, closed)
		if !errors.Is(err, errBroken) || !errors.Is(err, errClosed) {
			t.Errorf("全部失敗：錯誤 = %v, want 同時包含 %v 和 %v", err, errBroken, errClosed)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, err = Fastest(ctx, kitchen(time.Hour, "西區 排骨飯", &stopped))
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("呼叫的人的 ctx 逾時：錯誤 = %v, want DeadlineExceeded", err)
		}
	})
}

// ask 是一間要 10 秒才回答的廚房，會聽 ctx。
func ask(ctx context.Context) (time.Duration, error) {
	select {
	case <-time.After(10 * time.Second):
		return 15 * time.Minute, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func TestQuote(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		start := time.Now()
		_, err := Quote(context.Background(), ask)
		if !errors.Is(err, context.DeadlineExceeded) || time.Since(start) != time.Second {
			t.Errorf("廚房太慢：Quote = %v，花了 %v; want DeadlineExceeded，1s", err, time.Since(start))
		}

		ctx, cancel := context.WithCancel(context.Background())
		time.AfterFunc(200*time.Millisecond, cancel)
		start = time.Now()
		_, err = Quote(ctx, ask)
		if !errors.Is(err, context.Canceled) || time.Since(start) != 200*time.Millisecond {
			t.Errorf("呼叫的人 0.2 秒後取消：Quote = %v，花了 %v; want Canceled，0.2s", err, time.Since(start))
		}

		fast := func(ctx context.Context) (time.Duration, error) { return 5 * time.Minute, nil }
		if d, err := Quote(context.Background(), fast); d != 5*time.Minute || err != nil {
			t.Errorf("廚房馬上回答：Quote = %v, %v; want 5m0s, nil", d, err)
		}
	})
}

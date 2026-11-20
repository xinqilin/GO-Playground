package desk

import (
	"slices"
	"testing"
	"testing/synctest"
	"time"
)

// feed 回傳一個已經放好 vals、而且 close 的 channel。
func feed[T any](vals ...T) <-chan T {
	ch := make(chan T, len(vals))
	for _, v := range vals {
		ch <- v
	}
	close(ch)
	return ch
}

func TestTryTake(t *testing.T) {
	ch := make(chan string, 1)
	if v, ok := TryTake(ch); ok || v != "" {
		t.Errorf("空的 channel：TryTake = %q, %v; want \"\", false", v, ok)
	}
	ch <- "排骨飯"
	if v, ok := TryTake(ch); !ok || v != "排骨飯" {
		t.Errorf("有值時 TryTake = %q, %v; want 排骨飯, true", v, ok)
	}
	if v, ok := TryTake(feed[int]()); ok || v != 0 {
		t.Errorf("已 close 的 channel：TryTake = %d, %v; want 0, false", v, ok)
	}
	if v, ok := TryTake(feed(17)); !ok || v != 17 {
		t.Errorf("close 前放進去的值：TryTake = %d, %v; want 17, true", v, ok)
	}
}

func TestWithin(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		start := time.Now()
		v, ok := Within(make(chan string), 5*time.Second)
		if ok || v != "" || time.Since(start) != 5*time.Second {
			t.Errorf("一直沒有值：Within = %q, %v，等了 %v; want \"\", false，5s", v, ok, time.Since(start))
		}

		// ch 有緩衝、最後等 sent：不管 Within 寫得對不對，送的 goroutine 都能結束，不會讓泡泡 panic。
		ch := make(chan string, 1)
		sent := make(chan struct{})
		defer func() { <-sent }()
		go func() {
			defer close(sent)
			time.Sleep(2 * time.Second)
			ch <- "雞腿飯"
		}()
		start = time.Now()
		v, ok = Within(ch, 5*time.Second)
		if !ok || v != "雞腿飯" || time.Since(start) != 2*time.Second {
			t.Errorf("2 秒後有值：Within = %q, %v，等了 %v; want 雞腿飯, true，2s", v, ok, time.Since(start))
		}

		if _, ok := Within(feed[string](), 5*time.Second); ok {
			t.Errorf("已 close 的 channel：Within 的 ok 應該是 false")
		}
	})
}

// pending 數一數 ch 現在有幾個值可以不等待地收。
func pending(ch <-chan struct{}) int {
	n := 0
	for {
		select {
		case <-ch:
			n++
		default:
			return n
		}
	}
}

func TestBell(t *testing.T) {
	b := NewBell()
	for range 3 {
		b.Ring()
	}
	if n := pending(b.C()); n != 1 {
		t.Errorf("按了 3 次，店員看到 %d 個訊號, want 1", n)
	}
	if n := pending(b.C()); n != 0 {
		t.Errorf("看過之後還剩 %d 個訊號, want 0", n)
	}
	b.Ring()
	if n := pending(b.C()); n != 1 {
		t.Errorf("看過之後再按一次，看到 %d 個訊號, want 1", n)
	}
}

func TestTally(t *testing.T) {
	if got := Tally(feed(85, 95, 30), feed(30)); got != 180 {
		t.Errorf("Tally = %d, want 180", got)
	}
	if got := Tally(feed(85, 95), nil); got != 180 {
		t.Errorf("refunds 是 nil：Tally = %d, want 180", got)
	}
	if got := Tally(nil, nil); got != 0 {
		t.Errorf("兩個都是 nil：Tally = %d, want 0", got)
	}

	sales := make(chan int)
	refunds := make(chan int)
	go func() {
		defer close(sales)
		for _, v := range []int{100, 100, 100} {
			sales <- v
		}
	}()
	go func() {
		defer close(refunds)
		refunds <- 50
	}()
	if got := Tally(sales, refunds); got != 250 {
		t.Errorf("兩邊同時送：Tally = %d, want 250", got)
	}
}

func TestCounter(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		done := make(chan struct{})
		out := Counter(done)
		got := []int{<-out, <-out, <-out}
		if !slices.Equal(got, []int{1, 2, 3}) {
			t.Errorf("前三個 = %v, want [1 2 3]", got)
		}
		close(done)

		timeout := time.After(time.Minute)
		for {
			select {
			case _, ok := <-out:
				if !ok {
					return
				}
			case <-timeout:
				t.Errorf("done 被 close 之後，回傳的 channel 一直沒有 close")
				return
			}
		}
	})
}

func TestDebounce(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// in 有緩衝、最後等 sent：Debounce 沒有在收時，送的 goroutine 也不會卡住。
		in := make(chan string, 4)
		sent := make(chan struct{})
		defer func() { <-sent }()
		start := time.Now()
		go func() {
			defer close(sent)
			in <- "a"
			time.Sleep(500 * time.Millisecond)
			in <- "b"
			time.Sleep(300 * time.Millisecond)
			in <- "c" // 0.8s
			time.Sleep(2 * time.Second)
			in <- "d" // 2.8s
			time.Sleep(200 * time.Millisecond)
			close(in) // 3.0s：d 還沒送出，要立刻送
		}()

		var got []string
		var at []time.Duration
		for v := range Debounce(in, time.Second) {
			got = append(got, v)
			at = append(at, time.Since(start))
		}
		if !slices.Equal(got, []string{"c", "d"}) {
			t.Fatalf("Debounce 送出 %q, want [c d]", got)
		}
		if want := []time.Duration{1800 * time.Millisecond, 3 * time.Second}; !slices.Equal(at, want) {
			t.Errorf("送出的時間 = %v, want %v", at, want)
		}
	})
}

func TestCollectFor(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ch := make(chan int, 10)
		go func() {
			for i := 1; i <= 10; i++ {
				time.Sleep(time.Second)
				ch <- i
			}
			close(ch)
		}()
		start := time.Now()
		got := CollectFor(ch, 3500*time.Millisecond)
		if !slices.Equal(got, []int{1, 2, 3}) {
			t.Errorf("每秒來一個、期限 3.5 秒：CollectFor = %v，花了 %v; want [1 2 3]", got, time.Since(start))
		}
		for range ch {
		}
	})
}

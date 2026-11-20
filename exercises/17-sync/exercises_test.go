package stall

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestJar(t *testing.T) {
	var j Jar
	stop := make(chan struct{})
	var wg sync.WaitGroup
	// 一邊投錢，一邊不停地看：每次看到的金額都必須剛好是次數的 10 倍。
	wg.Go(func() {
		for {
			select {
			case <-stop:
				return
			default:
			}
			if times, total := j.Count(); total != times*10 {
				t.Errorf("Count 讀到 %d 次、%d 元，兩個數字對不起來", times, total)
				return
			}
		}
	})
	var puts sync.WaitGroup
	for range 100 {
		puts.Go(func() { j.Put(10) })
	}
	puts.Wait()
	close(stop)
	wg.Wait()

	if times, total := j.Count(); times != 100 || total != 1000 {
		t.Errorf("100 個人各投 10 元：Count = %d 次、%d 元; want 100 次、1000 元", times, total)
	}
}

func TestShutter(t *testing.T) {
	var s Shutter
	if s.Closed() {
		t.Fatalf("零值的鐵捲門應該是開著的")
	}
	var won atomic.Int64
	var wg sync.WaitGroup
	for range 50 {
		wg.Go(func() {
			if s.Close() {
				won.Add(1)
			}
		})
	}
	wg.Wait()
	if n := won.Load(); n != 1 {
		t.Errorf("50 個人同時拉門，有 %d 個人的 Close 回傳 true, want 1", n)
	}
	if !s.Closed() {
		t.Errorf("Close 之後 Closed() 應該是 true")
	}
	if s.Close() {
		t.Errorf("門已經拉下來了，再 Close 應該回傳 false")
	}
}

func TestReservations(t *testing.T) {
	r := NewReservations()
	if _, ok := r.Who(5); ok {
		t.Fatalf("還沒有人訂，Who(5) 的 ok 應該是 false")
	}

	names := []string{"王先生", "林小姐", "陳同學", "張太太", "李老師"}
	errs := make([]error, 50)
	var wg sync.WaitGroup
	for i := range errs {
		wg.Go(func() { errs[i] = r.Book(5, names[i%len(names)]) })
	}
	wg.Wait()

	winner := -1
	for i, err := range errs {
		switch {
		case err == nil && winner == -1:
			winner = i
		case err == nil:
			t.Fatalf("5 號桌被訂了不只一次（第 %d 和第 %d 個人都成功）", winner, i)
		case !errors.Is(err, ErrTaken):
			t.Fatalf("第 %d 個人拿到 %v, want ErrTaken", i, err)
		}
	}
	if winner == -1 {
		t.Fatalf("50 個人都沒訂到 5 號桌")
	}
	if got, ok := r.Who(5); !ok || got != names[winner%len(names)] {
		t.Errorf("Who(5) = %q, %v; want %q, true", got, ok, names[winner%len(names)])
	}

	if err := r.Book(6, "王先生"); err != nil {
		t.Errorf("6 號桌還沒有人訂，Book 回傳 %v", err)
	}
}

func TestCurrent(t *testing.T) {
	var c Current
	if got := c.Get(); got != (Settings{}) {
		t.Fatalf("零值的 Get() = %+v, want 零值", got)
	}
	c.Update(func(s Settings) Settings {
		s.Open = true
		s.Discount = 9
		return s
	})

	var wg sync.WaitGroup
	for range 100 {
		wg.Go(func() {
			c.Update(func(s Settings) Settings {
				s.Served++
				return s
			})
		})
	}
	wg.Wait()
	if got, want := c.Get(), (Settings{Open: true, Discount: 9, Served: 100}); got != want {
		t.Errorf("Get() = %+v, want %+v", got, want)
	}
}

func TestOnceOK(t *testing.T) {
	var calls, running, overlap atomic.Int64
	f := func() (string, error) {
		n := calls.Add(1)
		if running.Add(1) > 1 {
			overlap.Add(1)
		}
		defer running.Add(-1)
		time.Sleep(time.Millisecond)
		if n <= 3 {
			return "", errors.New("價目表讀取失敗")
		}
		return "排骨飯 90", nil
	}
	get := OnceOK(f)

	for i := range 3 {
		if _, err := get(); err == nil {
			t.Fatalf("第 %d 次呼叫：f 失敗了，應該回傳錯誤", i+1)
		}
	}
	if v, err := get(); err != nil || v != "排骨飯 90" {
		t.Fatalf("第 4 次呼叫 = %q, %v; want 排骨飯 90, nil", v, err)
	}

	var wg sync.WaitGroup
	for range 20 {
		wg.Go(func() {
			if v, err := get(); err != nil || v != "排骨飯 90" {
				t.Errorf("成功之後再呼叫 = %q, %v; want 排骨飯 90, nil", v, err)
			}
		})
	}
	wg.Wait()
	if n := calls.Load(); n != 4 {
		t.Errorf("f 被呼叫了 %d 次, want 4（失敗 3 次、成功 1 次，之後不再呼叫）", n)
	}

	// 一開始就有很多人同時呼叫：f 一次就成功時，只能被執行一次。
	calls.Store(3)
	get = OnceOK(f)
	for range 20 {
		wg.Go(func() { get() })
	}
	wg.Wait()
	if n := calls.Load() - 3; n != 1 {
		t.Errorf("20 個人同時呼叫，f 被執行了 %d 次, want 1", n)
	}
	if n := overlap.Load(); n != 0 {
		t.Errorf("有 %d 次是兩個 f 同時在執行", n)
	}
}

// summary 在另一個 goroutine 呼叫 Summary，一秒內沒有回傳就當作死結。
func summary(t *testing.T, q *Queue) string {
	t.Helper()
	got := make(chan string, 1)
	go func() { got <- q.Summary() }()
	select {
	case s := <-got:
		return s
	case <-time.After(time.Second):
		t.Fatalf("Summary 一秒內沒有回傳，可能死結了")
		return ""
	}
}

func TestQueueSummary(t *testing.T) {
	var q Queue
	if got := summary(t, &q); got != "目前不用等" {
		t.Errorf("沒有人時 Summary = %q, want 目前不用等", got)
	}
	q.Join("王先生")
	q.Join("林小姐")
	q.Join("陳同學")
	if got, want := summary(t, &q), "3 組候位，下一組：王先生"; got != want {
		t.Errorf("Summary = %q, want %q", got, want)
	}
}

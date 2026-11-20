package rush

import (
	"context"
	"errors"
	"maps"
	"slices"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

// gauge 記錄同時有幾個工作在跑，以及最多的時候是幾個。
type gauge struct {
	running, peak atomic.Int64
}

func (g *gauge) enter() {
	now := g.running.Add(1)
	for {
		p := g.peak.Load()
		if now <= p || g.peak.CompareAndSwap(p, now) {
			return
		}
	}
}

func (g *gauge) leave() { g.running.Add(-1) }

func TestLimit(t *testing.T) {
	for _, tc := range []struct {
		n, tasks int
		took     time.Duration
		peak     int64
	}{
		{2, 6, 3 * time.Second, 2},
		{3, 7, 3 * time.Second, 3},
		{10, 4, time.Second, 4},
	} {
		synctest.Test(t, func(t *testing.T) {
			var g gauge
			var ran atomic.Int64
			tasks := make([]func(), tc.tasks)
			for i := range tasks {
				tasks[i] = func() {
					g.enter()
					defer g.leave()
					time.Sleep(time.Second)
					ran.Add(1)
				}
			}
			start := time.Now()
			Limit(tc.n, tasks...)
			if ran.Load() != int64(tc.tasks) || time.Since(start) != tc.took || g.peak.Load() != tc.peak {
				t.Errorf("Limit(%d, %d 個工作)：做了 %d 個、花 %v、同時最多 %d 個；want %d 個、%v、%d 個",
					tc.n, tc.tasks, ran.Load(), time.Since(start), g.peak.Load(), tc.tasks, tc.took, tc.peak)
			}
		})
	}
}

func TestMapOrdered(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var g gauge
		in := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		square := func(n int) int {
			g.enter()
			defer g.leave()
			time.Sleep(time.Second)
			return n * n
		}
		start := time.Now()
		got := MapOrdered(3, in, square)
		want := []int{1, 4, 9, 16, 25, 36, 49, 64, 81, 100}
		if !slices.Equal(got, want) {
			t.Fatalf("MapOrdered = %v, want %v", got, want)
		}
		if d, p := time.Since(start), g.peak.Load(); d != 4*time.Second || p != 3 {
			t.Errorf("10 個值、3 個 worker：花了 %v、同時最多 %d 個；want 4s、3 個", d, p)
		}
		if got := MapOrdered(8, []string{"a"}, func(s string) string { return s + s }); !slices.Equal(got, []string{"aa"}) {
			t.Errorf("worker 比值多：MapOrdered = %q, want [aa]", got)
		}
	})
}

func TestBatches(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// in 有緩衝、最後等 sent：Batches 沒有在收時，送的 goroutine 也不會卡住。
		in := make(chan int, 10)
		sent := make(chan struct{})
		defer func() { <-sent }()
		start := time.Now()
		go func() {
			defer close(sent)
			in <- 1
			in <- 2
			in <- 3 // 0s：湊滿 3 個
			in <- 4
			time.Sleep(500 * time.Millisecond)
			in <- 5 // 0.5s；4 是 0s 收到的，1s 時要送出
			time.Sleep(2500 * time.Millisecond)
			in <- 6 // 3s
			time.Sleep(200 * time.Millisecond)
			close(in) // 3.2s：6 還沒送出，要馬上送
		}()

		var got [][]int
		var at []time.Duration
		for b := range Batches(t.Context(), in, 3, time.Second) {
			got = append(got, b)
			at = append(at, time.Since(start))
		}
		want := [][]int{{1, 2, 3}, {4, 5}, {6}}
		if !slices.EqualFunc(got, want, slices.Equal) {
			t.Fatalf("Batches 送出 %v, want %v", got, want)
		}
		if wantAt := []time.Duration{0, time.Second, 3200 * time.Millisecond}; !slices.Equal(at, wantAt) {
			t.Errorf("送出的時間 = %v, want %v", at, wantAt)
		}
	})

	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
		defer cancel()
		in := make(chan int, 1)
		in <- 1
		start := time.Now()
		var got [][]int
		for b := range Batches(ctx, in, 3, time.Hour) {
			got = append(got, b)
		}
		if len(got) != 0 || time.Since(start) != 2*time.Second {
			t.Errorf("2 秒後取消：送出 %v、在 %v 結束；want 沒有送出、2s", got, time.Since(start))
		}
	})
}

func TestTeam(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var g gauge
		var ran atomic.Int64
		team, ctx := NewTeam(t.Context(), 2)
		start := time.Now()
		for range 5 {
			team.Go(func() error {
				g.enter()
				defer g.leave()
				time.Sleep(time.Second)
				ran.Add(1)
				return nil
			})
		}
		if err := team.Wait(); err != nil {
			t.Fatalf("全部成功：Wait = %v, want nil", err)
		}
		if ran.Load() != 5 || time.Since(start) != 3*time.Second || g.peak.Load() != 2 {
			t.Errorf("5 個工作、同時 2 個：做了 %d 個、花 %v、同時最多 %d 個；want 5 個、3s、2 個",
				ran.Load(), time.Since(start), g.peak.Load())
		}
		if ctx.Err() == nil {
			t.Errorf("Wait 回傳之後，ctx 應該被取消")
		}
	})

	synctest.Test(t, func(t *testing.T) {
		errBroken := errors.New("西區 瓦斯沒了")
		var stopped atomic.Int64
		team, ctx := NewTeam(t.Context(), 0)
		start := time.Now()
		team.Go(func() error {
			time.Sleep(time.Second)
			return errBroken
		})
		for range 3 {
			team.Go(func() error {
				select {
				case <-time.After(time.Hour):
					return nil
				case <-ctx.Done():
					stopped.Add(1)
					return errors.New("被取消了")
				}
			})
		}
		err := team.Wait()
		if err != errBroken || time.Since(start) != time.Second || stopped.Load() != 3 {
			t.Errorf("Wait = %v、花 %v、%d 個收到取消；want %v、1s、3 個", err, time.Since(start), stopped.Load(), errBroken)
		}
	})
}

// prices 在另一個 goroutine 呼叫 Prices，一秒內沒有回傳就當作死結。
func prices(t *testing.T, workers int, dishes []string) map[string]int {
	t.Helper()
	menu := map[string]int{"排骨飯": 90, "雞腿飯": 100, "魚排飯": 95, "滷肉飯": 50, "控肉飯": 80}
	got := make(chan map[string]int, 1)
	go func() { got <- Prices(workers, dishes, func(d string) int { return menu[d] }) }()
	select {
	case m := <-got:
		return m
	case <-time.After(time.Second):
		t.Fatalf("Prices(%d, %d 道菜) 一秒內沒有回傳，可能死結了", workers, len(dishes))
		return nil
	}
}

func TestPrices(t *testing.T) {
	want := map[string]int{"排骨飯": 90, "雞腿飯": 100, "魚排飯": 95, "滷肉飯": 50, "控肉飯": 80}
	dishes := slices.Sorted(maps.Keys(want))
	if got := prices(t, 2, dishes); !maps.Equal(got, want) {
		t.Errorf("Prices = %v, want %v", got, want)
	}
	if got := prices(t, 8, dishes[:1]); !maps.Equal(got, map[string]int{dishes[0]: want[dishes[0]]}) {
		t.Errorf("一道菜：Prices = %v", got)
	}
	if got := prices(t, 3, nil); len(got) != 0 {
		t.Errorf("沒有菜：Prices = %v, want 空的 map", got)
	}
}

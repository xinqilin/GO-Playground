package main

import (
	"context"
	"errors"
	"math/rand/v2"
	"slices"
	"strconv"
	"strings"
	"testing"
	"testing/synctest"
	"time"
)

func Example_poolDemo() {
	poolDemo()
	// Output:
	// 排骨飯 好了
	// 雞腿飯 好了
	// 魚排飯 好了
	// 滷肉飯 好了
	// 控肉飯 好了
}

func Example_semaphoreDemo() {
	semaphoreDemo()
	// Output:
	// [12 0 7 30 5 18 9]
	// 同時在查的分店沒有超過 2 家：true
}

func Example_pipelineDemo() {
	pipelineDemo()
	// Output:
	// 停下來： 發票金額「九十」不是數字
	// 含稅合計 215
}

func Example_orderDemo() {
	orderDemo()
	// Output:
	// 排骨飯
	// 雞腿飯
	// 魚排飯
	// 滷肉飯
}

func Example_groupDemo() {
	groupDemo()
	// Output:
	// 錯誤： 西區 排骨缺貨
	// 東區 停下來：context canceled
	// 西區 缺貨
	// 北區 停下來：context canceled
}

func Example_tickerDemo() {
	tickerDemo()
	// Output:
	// 已通知 0912-000-111
	// 已通知 0922-000-222
	// 已通知 0933-000-333
	// 錯誤： <nil>
}

// TestFetchAllLimit：7 家分店、每家查 1 秒、同時最多 2 家，用假時間剛好 4 秒，同時最多剛好 2 家。
func TestFetchAllLimit(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var g gauge
		start := time.Now()
		fetchAll(make([]string, 7), 2, func(string) int {
			g.enter()
			defer g.leave()
			time.Sleep(time.Second)
			return 0
		})
		if d, p := time.Since(start), g.peak.Load(); d != 4*time.Second || p != 2 {
			t.Errorf("花了 %v、同時最多 %d 個；want 4s、2 個", d, p)
		}
	})
}

// TestCookAllCancel：一個廚師、每張單 1 秒，2.5 秒時取消：手上那張做完，後面的不再派。
func TestCookAllCancel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		time.AfterFunc(2500*time.Millisecond, cancel)
		start := time.Now()
		got := cookAll(ctx, []string{"a", "b", "c", "d", "e"}, 1, func(o string) string {
			time.Sleep(time.Second)
			return strings.ToUpper(o)
		})
		if want := []string{"A", "B", "C", "", ""}; !slices.Equal(got, want) || time.Since(start) != 3*time.Second {
			t.Errorf("cookAll = %q，花了 %v；want %q，3s", got, time.Since(start), want)
		}
	})
}

// TestPipelineStops：消費的一方提早離開並取消，管線的每一站都要結束，泡泡裡不留 goroutine。
func TestPipelineStops(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		out := stage(ctx, stage(ctx, emit(ctx, "1", "2", "3", "4"), parse), addTax)
		if a := <-out; a.v != 1 || a.err != nil {
			t.Errorf("第一個 = %+v, want {v:1}", a)
		}
		cancel()
	})
}

func TestInOrder(t *testing.T) {
	const n = 50
	in := make(chan finished, n)
	for _, i := range rand.Perm(n) {
		in <- finished{i, strconv.Itoa(i)}
	}
	close(in)
	next := 0
	for dish := range inOrder(in) {
		if dish != strconv.Itoa(next) {
			t.Fatalf("第 %d 個送出 %q", next, dish)
		}
		next++
	}
	if next != n {
		t.Errorf("送出了 %d 個, want %d", next, n)
	}
}

func TestGroup(t *testing.T) {
	g, ctx := NewGroup(t.Context())
	for range 5 {
		g.Go(func() error { return nil })
	}
	if err := g.Wait(); err != nil || ctx.Err() == nil {
		t.Errorf("全部成功：Wait = %v、ctx.Err = %v；want nil、Wait 之後 ctx 被取消", err, ctx.Err())
	}

	errFirst := errors.New("第一個錯誤")
	g, ctx = NewGroup(t.Context())
	g.Go(func() error { return errFirst })
	for range 5 {
		g.Go(func() error {
			<-ctx.Done()
			return errors.New("後來的錯誤")
		})
	}
	if err := g.Wait(); err != errFirst || context.Cause(ctx) != errFirst {
		t.Errorf("Wait = %v、Cause = %v；want 都是 %v", err, context.Cause(ctx), errFirst)
	}
}

// TestNotifyAllPacing：每隔 1 秒通知一位，第 1、2、3 秒各一位；中途取消就停。
func TestNotifyAllPacing(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		start := time.Now()
		var at []time.Duration
		send := func(string) { at = append(at, time.Since(start)) }
		if err := notifyAll(t.Context(), []string{"a", "b", "c"}, time.Second, send); err != nil {
			t.Fatal(err)
		}
		if want := []time.Duration{time.Second, 2 * time.Second, 3 * time.Second}; !slices.Equal(at, want) {
			t.Errorf("通知的時間 = %v, want %v", at, want)
		}

		ctx, cancel := context.WithTimeout(t.Context(), 1500*time.Millisecond)
		defer cancel()
		at, start = nil, time.Now()
		err := notifyAll(ctx, []string{"a", "b", "c"}, time.Second, send)
		if !errors.Is(err, context.DeadlineExceeded) || len(at) != 1 {
			t.Errorf("1.5 秒逾時：錯誤 = %v、通知了 %d 位；want DeadlineExceeded、1 位", err, len(at))
		}
	})
}

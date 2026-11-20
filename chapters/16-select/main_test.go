package main

import (
	"slices"
	"testing"
	"testing/synctest"
	"time"
)

func Example_readyDemo() {
	readyDemo()
	// Output:
	// 處理 外送單 A-07
}

func Example_fairDemo() {
	fairDemo()
	// Output:
	// true true 1000
}

func Example_timeoutDemo() {
	timeoutDemo()
	// Output:
	// 拿到 排骨飯
	// 逾時
}

func Example_deadlineDemo() {
	deadlineDemo()
	// Output:
	// [排骨飯 雞腿飯]
}

func Example_tryReceiveDemo() {
	tryReceiveDemo()
	// Output:
	// 現在沒有
	// 收到 17
}

func Example_dropDemo() {
	dropDemo()
	// Output:
	// 2 3
}

func Example_loopDemo() {
	loopDemo()
	// Output:
	// 內用 排骨飯
	// 內用 雞腿飯
	// 外帶 魚排飯
}

func Example_findDemo() {
	findDemo()
	// Output:
	// 1
}

func Example_priorityDemo() {
	priorityDemo()
	// Output:
	// 急件 1
	// 急件 2
	// 普通 1
	// 普通 2
}

// TestFindNoLeak：find 提早 return 時，close(done) 讓 fromUntil 的 goroutine 結束，泡泡裡不會留下卡住的 goroutine。
func TestFindNoLeak(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		if got := find([]string{"排骨飯", "雞腿飯", "魚排飯"}, "排骨飯"); got != 0 {
			t.Errorf("find = %d, want 0", got)
		}
		if got := find([]string{"排骨飯"}, "紅茶"); got != -1 {
			t.Errorf("find(紅茶) = %d, want -1", got)
		}
	})
}

// TestCollectUntilIsOverallDeadline：每秒來一份，期限 3.5 秒，只收得到前三份。
// 泡泡裡的 time.Sleep 用的是假時間，測試不必真的等。最後把剩下的收完，讓送的 goroutine 能結束。
func TestCollectUntilIsOverallDeadline(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ch := make(chan string, 10)
		go func() {
			for _, o := range []string{"a", "b", "c", "d", "e"} {
				time.Sleep(time.Second)
				ch <- o
			}
			close(ch)
		}()
		start := time.Now()
		got := collectUntil(ch, 3500*time.Millisecond)
		if !slices.Equal(got, []string{"a", "b", "c"}) || time.Since(start) != 3500*time.Millisecond {
			t.Errorf("collectUntil = %q，花了 %v；want [a b c]，3.5s", got, time.Since(start))
		}
		for range ch {
		}
	})
}

func TestNextPrefersUrgent(t *testing.T) {
	urgent := make(chan string, 1)
	normal := make(chan string, 1)
	for range 100 {
		urgent <- "急"
		normal <- "普"
		if got := next(urgent, normal); got != "急" {
			t.Fatalf("兩邊都有單時 next = %q, want 急", got)
		}
		<-normal
	}
}

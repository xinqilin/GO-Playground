package main

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"
)

func Example_basicsDemo() {
	basicsDemo()
	// Output:
	// 滷蛋 好了
	// 佛跳牆 放棄： context deadline exceeded
}

func Example_treeDemo() {
	treeDemo()
	// Output:
	// 總店 <nil>
	// 東區 context canceled
	// 東區的訂單 context canceled
	// 西區 <nil>
	// 總店關門之後，西區 context canceled
}

func Example_errDemo() {
	errDemo()
	// Output:
	// context canceled：Canceled=true DeadlineExceeded=false
	// context deadline exceeded：Canceled=false DeadlineExceeded=true
}

func Example_deadlineDemo() {
	deadlineDemo()
	// Output:
	// 要 1 小時，實際期限和上層一樣：true
	// 要 10 毫秒，期限比上層早：true
	// Background 有期限：false
}

func Example_causeDemo() {
	causeDemo()
	// Output:
	// context canceled / 排骨賣完了
	// 衍生的 ctx 也拿得到原因： 排骨賣完了
	// context deadline exceeded / 廚房出餐太慢
	// 沒有指定原因： context canceled
}

func Example_listenDemo() {
	listenDemo()
	// Output:
	// 排骨飯
	// 雞腿飯
}

func Example_tallyDemo() {
	tallyDemo()
	// Output:
	// 12502500 <nil>
	// 0 context canceled
}

func Example_valueDemo() {
	valueDemo()
	// Output:
	// [A-07] 開始備餐
	// [-] 沒有訂單的紀錄
}

func Example_detachDemo() {
	detachDemo()
	// Output:
	// 請求： context canceled
	// 對帳紀錄： <nil> A-07
}

func Example_afterDemo() {
	afterDemo()
	// Output:
	// 通知外送員：訂單取消
	// 撤回通知： true
}

// TestCookStopsOnCancel：泡泡裡用假時間，2 秒後取消，cook 在 2 秒時回來，不會等滿 1 小時。
func TestCookStopsOnCancel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		time.AfterFunc(2*time.Second, cancel)
		start := time.Now()
		err := cook(ctx, time.Hour)
		if !errors.Is(err, context.Canceled) || time.Since(start) != 2*time.Second {
			t.Errorf("cook = %v，花了 %v；want context.Canceled，2s", err, time.Since(start))
		}
	})
}

// TestPrepareNoLeak：提早不收、cancel 之後，prepare 的 goroutine 要結束，泡泡裡不留卡住的 goroutine。
func TestPrepareNoLeak(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		ch := prepare(ctx, "排骨飯", "雞腿飯", "魚排飯")
		if got := <-ch; got != "排骨飯" {
			t.Errorf("第一份 = %q, want 排骨飯", got)
		}
		cancel()
	})
}

func TestTally(t *testing.T) {
	amounts := make([]int, 2500)
	for i := range amounts {
		amounts[i] = 2
	}
	if got, err := tally(t.Context(), amounts); got != 5000 || err != nil {
		t.Errorf("tally = %d, %v; want 5000, nil", got, err)
	}
	ctx, cancel := context.WithTimeout(t.Context(), 0)
	defer cancel()
	if _, err := tally(ctx, amounts); !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("期限已過：tally 的錯誤 = %v, want DeadlineExceeded", err)
	}
}

package main

import (
	"slices"
	"testing"
	"testing/synctest"
)

func Example_handoffDemo() {
	handoffDemo()
	// Output:
	// 拿到 排骨飯
	// 拿到 雞腿飯
}

func Example_doneDemo() {
	doneDemo()
	// Output:
	// 排骨飯 x2，NT$170
}

func Example_bufferDemo() {
	bufferDemo()
	// Output:
	// 3 3
	// 叫號 17
	// 3 3
	// 叫號 18 19 20
	// 0 3
}

func Example_rangeDemo() {
	rangeDemo()
	// Output:
	// 出餐 排骨飯
	// 出餐 雞腿飯
	// 出餐 魚排飯
}

func Example_closedDemo() {
	closedDemo()
	// Output:
	// 17 true
	// 18 true
	// 0 false
	// 0 false
}

func Example_closePanicDemo() {
	closePanicDemo()
	// Output:
	// panic: send on closed channel
	// panic: close of closed channel
	// panic: close of nil channel
}

func Example_directionDemo() {
	directionDemo()
	// Output:
	// 請 17 號
	// 請 18 號
	// 請 19 號
}

func Example_pipelineDemo() {
	pipelineDemo()
	// Output:
	// 外帶袋[排骨飯（熱）]
	// 外帶袋[雞腿飯（熱）]
}

func Example_fanDemo() {
	fanDemo()
	// Output:
	// 5 [排骨飯好了 控肉飯好了 滷肉飯好了 雞腿飯好了 魚排飯好了]
}

func Example_findDemo() {
	findDemo()
	// Output:
	// 1
}

// TestFindNoLeak 在 synctest 的泡泡裡執行 find：
// 泡泡結束時如果還有 goroutine 卡著，synctest 會讓測試失敗。
func TestFindNoLeak(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		if got := find([]string{"排骨飯", "雞腿飯", "魚排飯"}, "雞腿飯"); got != 1 {
			t.Errorf("find = %d, want 1", got)
		}
	})
}

func TestFind(t *testing.T) {
	items := []string{"排骨飯", "雞腿飯"}
	if got := find(items, "排骨飯"); got != 0 {
		t.Errorf("find(排骨飯) = %d, want 0", got)
	}
	if got := find(items, "紅茶"); got != -1 {
		t.Errorf("find(紅茶) = %d, want -1", got)
	}
}

func TestKitchenOneCook(t *testing.T) {
	var got []string
	for m := range kitchen(orders("a", "b", "c"), 1) {
		got = append(got, m)
	}
	if want := []string{"a好了", "b好了", "c好了"}; !slices.Equal(got, want) {
		t.Errorf("一位廚師時 kitchen = %q, want %q", got, want)
	}
}

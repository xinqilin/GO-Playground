package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

func Example_lostUpdateDemo() {
	lostUpdateDemo()
	// Output:
	// 兩個店員各記一份，結果是 6
}

func Example_stockDemo() {
	stockDemo()
	// Output:
	// 100 個人搶 30 份：賣出 30 份，剩 0 份
}

func Example_oversellDemo() {
	oversellDemo()
	// Output:
	// A 買到了
	// B 買到了
	// 只有 1 份，卻賣出 2 份，剩 0
}

func Example_menuDemo() {
	menuDemo()
	// Output:
	// 50 個人查到的不是 90 就是 95：true
	// 現在的價格 95
}

func Example_onceDemo() {
	onceDemo()
	// Output:
	// 讀取菜單
	// [排骨飯 雞腿飯] 讀取次數 1
}

func Example_onceErrDemo() {
	onceErrDemo()
	// Output:
	// 價目表讀取失敗
	// 價目表讀取失敗
	// 價目表讀取失敗
	// 實際嘗試 1 次
}

func Example_ticketDemo() {
	ticketDemo()
	// Output:
	// 發出 200 張，號碼不重複：true，最大 200
}

func Example_maxDemo() {
	maxDemo()
	// Output:
	// 今天最大的一筆： 1260
}

func Example_noticeDemo() {
	noticeDemo()
	// Output:
	// 換下來的： 正常營業 21:00
	// 現在的： 颱風提早打烊 18:00
}

func Example_seenDemo() {
	seenDemo()
	// Output:
	// A-01 接單
	// A-02 接單
	// A-01 重複，略過
}

func Example_poolDemo() {
	poolDemo()
	// Output:
	// 排骨飯 NT$90
	// 雞腿飯 NT$100
}

func Example_orderDemo() {
	orderDemo()
	// Output:
	// [排骨飯 雞腿飯]
}

// TestStockTake：用 go test -race 執行時，Take 少了鎖就會被抓到；
// 少了「檢查和減一在同一段鎖裡」，賣出的份數就可能超過庫存。
func TestStockTake(t *testing.T) {
	for _, tc := range []struct{ left, buyers, want int }{
		{30, 100, 30},
		{100, 30, 30},
		{0, 10, 0},
	} {
		s := &Stock{left: tc.left}
		var sold atomic.Int64
		var wg sync.WaitGroup
		for range tc.buyers {
			wg.Go(func() {
				if s.Take() {
					sold.Add(1)
				}
			})
		}
		wg.Wait()
		if got := int(sold.Load()); got != tc.want || s.Left() != tc.left-tc.want {
			t.Errorf("庫存 %d、%d 人搶：賣出 %d、剩 %d；want 賣出 %d、剩 %d",
				tc.left, tc.buyers, got, s.Left(), tc.want, tc.left-tc.want)
		}
	}
}

func TestRaiseTo(t *testing.T) {
	var best atomic.Int64
	var wg sync.WaitGroup
	for v := range int64(1000) {
		wg.Go(func() { raiseTo(&best, v) })
	}
	wg.Wait()
	if got := best.Load(); got != 999 {
		t.Errorf("best = %d, want 999", got)
	}
}

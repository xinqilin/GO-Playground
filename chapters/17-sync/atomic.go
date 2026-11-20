package main

import (
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
)

// ticketDemo：200 個人同時抽號碼牌。Add 一次完成「加一並拿到加完的值」，每個人拿到的號碼都不一樣。
func ticketDemo() {
	var next atomic.Int64
	got := make([]int64, 200)
	var wg sync.WaitGroup
	for i := range got {
		wg.Go(func() { got[i] = next.Add(1) })
	}
	wg.Wait()

	slices.Sort(got)
	last := got[len(got)-1]
	unique := len(slices.Compact(got)) == len(got)
	fmt.Printf("發出 %d 張，號碼不重複：%t，最大 %d\n", next.Load(), unique, last)
}

// raiseTo 把 best 更新成 best 和 v 比較大的那一個。
// Load 和 Store 分開做，中間會被插隊，所以用 CompareAndSwap：
// 只有 best 還是剛剛讀到的那個值時才寫入，否則重讀再試一次。
func raiseTo(best *atomic.Int64, v int64) {
	for {
		cur := best.Load()
		if v <= cur || best.CompareAndSwap(cur, v) {
			return
		}
	}
}

func maxDemo() {
	var best atomic.Int64
	var wg sync.WaitGroup
	for _, amount := range []int64{320, 85, 1260, 470, 990} {
		wg.Go(func() { raiseTo(&best, amount) })
	}
	wg.Wait()
	fmt.Println("今天最大的一筆：", best.Load())
}

// Notice 是店門口的公告。整份公告一起換，讀的人不會看到換到一半的內容。
type Notice struct {
	Text  string
	Until string
}

func noticeDemo() {
	var board atomic.Pointer[Notice]
	board.Store(&Notice{Text: "正常營業", Until: "21:00"})

	old := board.Swap(&Notice{Text: "颱風提早打烊", Until: "18:00"})
	now := board.Load()
	fmt.Println("換下來的：", old.Text, old.Until)
	fmt.Println("現在的：", now.Text, now.Until)
}

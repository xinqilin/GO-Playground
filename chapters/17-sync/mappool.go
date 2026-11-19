package main

import (
	"bytes"
	"fmt"
	"sync"
)

// seenDemo：外送平台可能把同一張單通知好幾次。
// LoadOrStore 一次完成「沒有就放進去」，loaded 告訴你之前是不是已經有了。
func seenDemo() {
	var seen sync.Map
	for _, id := range []string{"A-01", "A-02", "A-01"} {
		if _, loaded := seen.LoadOrStore(id, true); loaded {
			fmt.Println(id, "重複，略過")
			continue
		}
		fmt.Println(id, "接單")
	}
}

// slips 是收據用的 Buffer 回收箱：要用時先拿舊的，沒有才 New 一個新的。
var slips = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}

// receipt 借一個 Buffer 來組收據，用完 Reset 再還回去。
// b.String() 會複製一份內容，所以 Buffer 還回去之後，回傳的字串不受影響。
func receipt(item string, price int) string {
	b := slips.Get().(*bytes.Buffer)
	defer func() {
		b.Reset()
		slips.Put(b)
	}()
	fmt.Fprintf(b, "%s NT$%d", item, price)
	return b.String()
}

func poolDemo() {
	fmt.Println(receipt("排骨飯", 90))
	fmt.Println(receipt("雞腿飯", 100))
}

package main

import (
	"bytes"
	"fmt"
	"strconv"
	"sync"
)

// bufPool 存放用過的緩衝區。拿不到時 New 建一個新的。
var bufPool = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}

// maxPooled 以上的緩衝區不放回去，免得偶爾一張大收據讓 pool 一直佔著大塊記憶體。
const maxPooled = 4 << 10

// receipt 把訂單排成收據，每次都用新的緩衝區。
func receipt(orders []Order) string {
	var b bytes.Buffer
	writeReceipt(&b, orders)
	return b.String()
}

// receiptPooled 和 receipt 一樣，但緩衝區從 pool 借、用完還回去。
func receiptPooled(orders []Order) string {
	b := bufPool.Get().(*bytes.Buffer)
	b.Reset() // 上一個人留下的內容要清掉
	writeReceipt(b, orders)
	s := b.String() // 先複製出字串，緩衝區還回去之後就可能被別人改掉
	if b.Cap() <= maxPooled {
		bufPool.Put(b)
	}
	return s
}

func writeReceipt(b *bytes.Buffer, orders []Order) {
	total := 0
	for _, o := range orders {
		b.WriteString(o.Time)
		b.WriteByte(' ')
		b.WriteString(o.Dish)
		b.WriteString(" x")
		b.WriteString(strconv.Itoa(o.Qty))
		b.WriteByte('\n')
		total += o.Qty * o.Price
	}
	b.WriteString("合計 ")
	b.WriteString(strconv.Itoa(total))
	b.WriteByte('\n')
}

// poolDemo：兩種寫法印出來的收據一樣。
func poolDemo() {
	orders := []Order{{"12:01", "排骨飯", 2, 90}, {"12:03", "滷肉飯", 1, 50}}
	a, b := receipt(orders), receiptPooled(orders)
	fmt.Print(b)
	fmt.Println("相同：", a == b)
}

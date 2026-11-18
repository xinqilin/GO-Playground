package main

import "fmt"

// Page 是分頁查詢的結果：這一頁的資料，以及下一頁從哪個號碼開始（0 表示沒有下一頁）。
type Page[T any] struct {
	Items []T
	Next  int
}

// Map 是 generic method：除了 receiver 帶來的 T，方法自己再宣告一個型別參數 U。
func (p Page[T]) Map[U any](f func(T) U) Page[U] {
	out := Page[U]{Items: make([]U, len(p.Items)), Next: p.Next}
	for i, v := range p.Items {
		out.Items[i] = f(v)
	}
	return out
}

// Order 是一張訂單。
type Order struct {
	No   int
	Item string
}

func methodDemo() {
	p := Page[Order]{Items: []Order{{17, "排骨飯"}, {18, "雞腿飯"}}, Next: 19}
	labels := p.Map(func(o Order) string {
		return fmt.Sprintf("%d 號 %s", o.No, o.Item)
	})
	fmt.Printf("%T\n", labels)
	fmt.Println(labels.Items, labels.Next)
}

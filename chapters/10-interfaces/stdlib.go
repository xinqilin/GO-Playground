package main

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"sort"
)

// ErrOutOfPaper 表示收據機的紙用完了。
var ErrOutOfPaper = errors.New("紙用完了")

// PaperRoll 是收據機的紙捲，實作 io.Writer：最多只能印 Cap 個 byte，印不下的部分回傳錯誤。
type PaperRoll struct {
	Cap int
	buf []byte
}

func (r *PaperRoll) Write(p []byte) (int, error) {
	room := r.Cap - len(r.buf)
	if len(p) <= room {
		r.buf = append(r.buf, p...)
		return len(p), nil
	}
	r.buf = append(r.buf, p[:room]...)
	return room, ErrOutOfPaper
}

// writerDemo：fmt.Fprintf 只要求一個 io.Writer，並把 Write 的回傳值原樣交還給你。
func writerDemo() {
	roll := &PaperRoll{Cap: 19}
	n, err := fmt.Fprintf(roll, "排骨飯 x%d\n", 2)
	fmt.Println(n, err)
	n, err = fmt.Fprintf(roll, "滷蛋 x%d\n", 3)
	fmt.Println(n, err)
	fmt.Printf("%q\n", roll.buf)
}

// Order 是一張外帶單。
type Order struct {
	No  int
	VIP bool
}

// byPriority 實作 sort.Interface：VIP 排前面，同樣是或不是 VIP 時號碼小的先。
type byPriority []Order

func (o byPriority) Len() int      { return len(o) }
func (o byPriority) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
func (o byPriority) Less(i, j int) bool {
	if o[i].VIP != o[j].VIP {
		return o[i].VIP
	}
	return o[i].No < o[j].No
}

// sortDemo：sort.Sort 要一個有 Len、Less、Swap 的值；slices.SortFunc 只要一個比較函式。
func sortDemo() {
	orders := []Order{{No: 3}, {No: 1, VIP: true}, {No: 2}, {No: 4, VIP: true}}
	sort.Sort(byPriority(orders))
	fmt.Println(orders)

	slices.SortFunc(orders, func(a, b Order) int {
		return cmp.Compare(a.No, b.No)
	})
	fmt.Println(orders)
}

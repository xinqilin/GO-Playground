package main

import "fmt"

type counter struct{ n int }

// newCounter 回傳區域變數的位址。這在 Go 是安全的：
// 編譯器發現變數離開函式後還會被用到，就把它放到 heap 上。
func newCounter() *counter {
	c := counter{n: 1}
	return &c
}

// returnLocalPointer 示範使用 newCounter 回傳的指標，以及兩次呼叫得到不同的變數。
func returnLocalPointer() {
	a := newCounter()
	b := newCounter()
	a.n++
	fmt.Println(a.n, b.n, a == b)
}

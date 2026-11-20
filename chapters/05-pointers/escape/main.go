// escape 是給逃逸分析看的小程式。用下面的指令看編譯器把哪些變數放到 heap：
//
//	go build -gcflags=-m ./chapters/05-pointers/escape
package main

import "os"

type counter struct{ n int }

func newCounter() *counter {
	c := counter{n: 1} // 位址被回傳出去，離開函式後還會被用到：放到 heap
	return &c
}

func localOnly() int {
	c := counter{n: 1} // 位址只在函式裡用：留在 stack
	p := &c
	p.n++
	return c.n
}

func main() {
	c := newCounter()
	c.n += localOnly()
	os.Exit(c.n - 3) // 結果是 3，所以結束碼是 0
}

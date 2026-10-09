// 第 05 章：指標。
//
// 指標是一張寫著「地址」的紙條：拿著它就能找到、也能改動原本的變數。
// 執行：go run ./chapters/05-pointers
package main

import "fmt"

func main() {
	fmt.Println("== 取地址與解參考 ==")
	addressAndDeref()

	fmt.Println("== 透過指標修改 ==")
	modifyThroughPointer()

	fmt.Println("== 傳值 vs 傳指標 ==")
	passByValue()

	fmt.Println("== nil 與 new ==")
	nilAndNew()

	fmt.Println("== 回傳區域變數的指標 ==")
	returnLocalPointer()
}

// addressAndDeref 示範 & 取得地址、* 照地址取值。
func addressAndDeref() {
	x := 5
	p := &x // p 的型別是 *int：「指向 int 的指標」
	fmt.Println("x =", x)
	fmt.Println("*p =", *p)
	// 地址的實際數值每次執行都不同，所以這裡只比較它是不是 x 的地址。
	fmt.Println("p == &x:", p == &x)
}

// modifyThroughPointer 示範透過指標改動原本的變數。
func modifyThroughPointer() {
	x := 5
	p := &x
	*p = 10 // 改的是 p 指向的那一格，也就是 x
	fmt.Println("x =", x)
}

// addOneValue 收到的是複本，改了也不影響呼叫端。
func addOneValue(n int) { n++ }

// addOnePointer 收到的是地址，改的是原本那一格。
func addOnePointer(n *int) { *n++ }

// passByValue 示範 Go 一律傳值：想讓函式改到呼叫端的變數，就傳指標。
func passByValue() {
	a, b := 1, 1
	addOneValue(a)
	addOnePointer(&b)
	fmt.Println("a =", a, "b =", b)
}

// nilAndNew 示範指標的零值，以及用 new 配置一個新變數。
func nilAndNew() {
	var p *int // 指標的零值是 nil：一張空白紙條
	fmt.Println("p == nil:", p == nil)
	// 對 nil 解參考（*p）會 panic，使用前要先確認不是 nil。

	q := new(int) // 配置一格 int（值為 0），回傳它的地址
	*q = 42
	fmt.Println("*q =", *q)

	r := new(7) // Go 1.26 起，new 可以直接給初始值
	fmt.Println("*r =", *r)
}

type counter struct{ n int }

// newCounter 回傳區域變數的地址。這在 Go 是安全的：
// 編譯器發現變數離開函式後還會被用到，就把它放到 heap 上。
func newCounter() *counter {
	c := counter{n: 1}
	return &c
}

// returnLocalPointer 示範使用 newCounter 回傳的指標。
func returnLocalPointer() {
	c := newCounter()
	c.n++ // 指向 struct 的指標可以直接用 . 存取欄位，不必寫 (*c).n
	fmt.Println("c.n =", c.n)
}

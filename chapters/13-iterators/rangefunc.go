package main

import (
	"fmt"
	"iter"
)

// callNumbers 回傳一個疊代器：依序交出 from 到 to 的號碼。
func callNumbers(from, to int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for n := from; n <= to; n++ {
			if !yield(n) {
				return
			}
		}
	}
}

func rangeFuncDemo() {
	for n := range callNumbers(17, 19) {
		fmt.Println("請", n, "號")
	}
}

// byHand：range 背後做的事，就是把迴圈本體包成 yield 函式，交給疊代器去呼叫。
func byHand() {
	seq := callNumbers(17, 19)
	seq(func(n int) bool {
		fmt.Println("請", n, "號")
		return true
	})
}

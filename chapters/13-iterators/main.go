// 第 13 章：疊代器。
//
// 疊代器是一台「一次吐一個」的發球機：range 每要一個，它就丟一個出來；你喊停（break），它就停。
// 執行：go run ./chapters/13-iterators
package main

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"strings"
)

// Countdown 回傳一個從 n 數到 1 的疊代器。
// iter.Seq[int] 其實就是 func(yield func(int) bool)：每呼叫一次 yield，就交出一個值。
func Countdown(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := n; i >= 1; i-- {
			if !yield(i) { // 使用端 break 時 yield 回傳 false，要立刻停
				return
			}
		}
	}
}

// Fib 回傳無限長的費氏數列，靠使用端 break 停下來。
func Fib() iter.Seq[int] {
	return func(yield func(int) bool) {
		a, b := 0, 1
		for {
			if !yield(a) {
				return
			}
			a, b = b, a+b
		}
	}
}

// Filter 把一個疊代器包成另一個：只交出 keep 回傳 true 的值。
func Filter[V any](seq iter.Seq[V], keep func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if keep(v) && !yield(v) {
				return
			}
		}
	}
}

// Take 只取前 n 個值，之後就停止上游的疊代器。
func Take[V any](seq iter.Seq[V], n int) iter.Seq[V] {
	return func(yield func(V) bool) {
		if n <= 0 {
			return
		}
		count := 0
		for v := range seq {
			if !yield(v) {
				return
			}
			count++
			if count == n {
				return
			}
		}
	}
}

func main() {
	fmt.Println("== 自己寫疊代器 ==")
	for n := range Countdown(3) {
		fmt.Println(n)
	}

	fmt.Println("== 無限數列，break 就停 ==")
	var small []int
	for v := range Fib() {
		if v > 50 {
			break
		}
		small = append(small, v)
	}
	fmt.Println(small)

	fmt.Println("== 把疊代器接起來 ==")
	even := Filter(Fib(), func(n int) bool { return n%2 == 0 })
	fmt.Println(slices.Collect(Take(even, 5)))

	fmt.Println("== 標準庫的疊代器 ==")
	m := map[string]int{"b": 2, "a": 1, "c": 3}
	fmt.Println(slices.Sorted(maps.Keys(m)))
	for i, v := range slices.All([]string{"x", "y"}) {
		fmt.Println(i, v)
	}
	for part := range strings.SplitSeq("a,b,c", ",") {
		fmt.Print(part)
	}
	fmt.Println()
}

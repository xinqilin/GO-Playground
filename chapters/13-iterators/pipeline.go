package main

import (
	"fmt"
	"iter"
	"slices"
)

// Filter 只交出 keep 回傳 true 的值。
func Filter[T any](seq iter.Seq[T], keep func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			if keep(v) && !yield(v) {
				return
			}
		}
	}
}

// Map 把每個值換成 f(v) 再交出去。
func Map[T, U any](seq iter.Seq[T], f func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range seq {
			if !yield(f(v)) {
				return
			}
		}
	}
}

// Take 最多交出 n 個值，然後停止，上游的疊代器也跟著停。
func Take[T any](seq iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
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

// naturals 從 1 開始無止盡地數下去，要靠下游喊停。
func naturals() iter.Seq[int] {
	return func(yield func(int) bool) {
		for n := 1; ; n++ {
			if !yield(n) {
				return
			}
		}
	}
}

// pipelineDemo：跳過尾數是 4 的號碼，格式化成三位數，取前 5 個。
func pipelineDemo() {
	nums := Filter(naturals(), func(n int) bool { return n%10 != 4 })
	labels := Map(nums, func(n int) string { return fmt.Sprintf("%03d", n) })
	fmt.Println(slices.Collect(Take(labels, 5)))
}

// lazyDemo：值是一個一個流過整條管線的，用不到的值根本不會被產生。
func lazyDemo() {
	src := func(yield func(int) bool) {
		for n := 1; n <= 3; n++ {
			fmt.Println("產生", n)
			if !yield(n) {
				return
			}
		}
	}
	doubled := Map(src, func(n int) int {
		fmt.Println("加倍", n)
		return n * 2
	})
	for v := range doubled {
		fmt.Println("拿到", v)
		if v == 4 {
			break
		}
	}
}

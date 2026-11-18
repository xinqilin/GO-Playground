package main

import (
	"fmt"
	"iter"
	"slices"
)

// pullBasics：iter.Pull 把疊代器換成「要一個、拿一個」的 next 函式。
func pullBasics() {
	next, stop := iter.Pull(slices.Values([]string{"A-001", "A-002"}))
	defer stop()

	v1, ok1 := next()
	v2, ok2 := next()
	v3, ok3 := next()
	fmt.Println(v1, ok1, v2, ok2)
	fmt.Printf("%q %v\n", v3, ok3)
}

// assign 把等待中的號碼和空出來的櫃台一一配對，任何一邊用完就停。
func assign(nums iter.Seq[int], counters iter.Seq[string]) iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		next, stop := iter.Pull(counters)
		defer stop()
		for n := range nums {
			c, ok := next()
			if !ok || !yield(n, c) {
				return
			}
		}
	}
}

func pullDemo() {
	nums := slices.Values([]int{17, 18, 19, 20})
	counters := slices.Values([]string{"1 號櫃台", "2 號櫃台", "3 號櫃台"})
	for n, c := range assign(nums, counters) {
		fmt.Println(n, "→", c)
	}
}

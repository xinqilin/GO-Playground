package main

import (
	"cmp"
	"fmt"
)

// Clamp 把 v 限制在 lo 到 hi 之間。T 可以是任何能用 < 比較的型別。
func Clamp[T cmp.Ordered](v, lo, hi T) T {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// CountOf 數 s 裡有幾個 v。只需要 ==，所以限制是 comparable。
func CountOf[T comparable](s []T, v T) int {
	n := 0
	for _, x := range s {
		if x == v {
			n++
		}
	}
	return n
}

func paramsDemo() {
	fmt.Println(Clamp(20, 1, 12))
	fmt.Println(Clamp(0.35, 0.0, 0.3))
	fmt.Println(Clamp("Z", "A", "M"))
	fmt.Println(Clamp[float64](5, 1, 12.5))
	fmt.Println(CountOf([]string{"排骨飯", "雞腿飯", "排骨飯"}, "排骨飯"))
}

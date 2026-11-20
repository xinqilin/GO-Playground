package main

import (
	"fmt"
	"strings"
)

// sum 接受任意個 int。在函式裡，nums 就是一個 []int。
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// joinWith 的可變參數一定要放在最後，前面可以有一般參數。
func joinWith(sep string, parts ...string) string {
	return strings.Join(parts, sep)
}

// variadic 示範三種呼叫方式：不給、給好幾個、把現成的 slice 用 ... 展開。
func variadic() {
	fmt.Println(sum(), sum(85), sum(85, 90, 15))
	prices := []int{85, 90, 15}
	fmt.Println(sum(prices...))
	fmt.Println(joinWith("、", "排骨", "雞腿", "滷蛋"))
}

package main

import (
	"fmt"
	"slices"
)

// copyBasics 示範 copy 複製的個數是兩邊長度的較小值，以及用 slices.Clone 複製一份獨立的 slice。
func copyBasics() {
	src := []int{85, 90, 15}
	dst := make([]int, 2)
	n := copy(dst, src)
	fmt.Println(n, dst)

	clone := slices.Clone(src)
	clone[0] = 0
	fmt.Println(src, clone)
}

// insertDelete 示範用 slices.Delete 和 slices.Insert 在中間刪除、插入。
func insertDelete() {
	queue := []string{"A01", "A02", "A03", "A04"}
	queue = slices.Delete(queue, 1, 2) // 刪掉 queue[1:2]，也就是 A02
	fmt.Println(queue)
	queue = slices.Insert(queue, 1, "V01")
	fmt.Println(queue)
}

// filterInPlace 示範原地過濾：s[:0] 是一個長度 0、但沿用同一個陣列的 slice，留下來的元素依序寫回前面。
func filterInPlace() {
	prices := []int{85, 120, 90, 150, 70}
	cheap := prices[:0]
	for _, p := range prices {
		if p < 100 {
			cheap = append(cheap, p)
		}
	}
	fmt.Println(cheap)
	fmt.Println(prices) // 前三格被覆寫了，後面是舊資料

	others := []int{85, 120, 90, 150, 70}
	others = slices.DeleteFunc(others, func(p int) bool { return p >= 100 })
	fmt.Println(others)
}

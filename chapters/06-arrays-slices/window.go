package main

import "fmt"

// sliceLiteral 示範建立 slice 的兩種方法：字面值和 make。
func sliceLiteral() {
	menu := []string{"排骨", "雞腿", "素食"} // 方括號裡沒有長度：這是 slice
	fmt.Println(menu, len(menu), cap(menu))

	orders := make([]int, 3)    // len 3、cap 3，內容是零值
	queue := make([]int, 0, 10) // len 0、cap 10：先保留空間
	fmt.Println(orders, len(orders), cap(orders))
	fmt.Println(queue, len(queue), cap(queue))
}

// sliceWindow 示範 slice 是看向陣列的窗：透過 slice 修改，改的是底層陣列。
func sliceWindow() {
	arr := [5]int{10, 20, 30, 40, 50}
	s := arr[1:4] // 索引 1、2、3，不含 4
	fmt.Println(s, len(s), cap(s))
	s[0] = 99
	fmt.Println(arr)
}

// sliceBounds 示範 [low:high] 可以省略兩端，以及 slice 再切一次，只要不超過 cap 就能往後延伸。
func sliceBounds() {
	s := []int{0, 1, 2, 3, 4, 5}
	fmt.Println(s[:2], s[4:], s[:])

	part := s[1:3]
	fmt.Println(part, len(part), cap(part))
	more := part[:4] // part 的長度只有 2，但容量有 5
	fmt.Println(more)
}

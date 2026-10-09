// 第 06 章：陣列與切片。
//
// 陣列是一排固定數量的格子；slice 是看向陣列的一扇窗，記著從哪裡開始看、看幾格、最多能看到哪。
// 執行：go run ./chapters/06-arrays-slices
package main

import (
	"fmt"
	"slices"
)

func main() {
	fmt.Println("== 陣列是值 ==")
	arrays()

	fmt.Println("== slice 是一扇窗 ==")
	sliceWindow()

	fmt.Println("== append 與容量 ==")
	appendGrowth()

	fmt.Println("== 共用底層陣列的陷阱 ==")
	sharedBacking()

	fmt.Println("== copy 與 slices 套件 ==")
	copyAndHelpers()

	fmt.Println("== 二維 slice ==")
	twoD()
}

// arrays 示範陣列：長度是型別的一部分，指定（=）會複製整個陣列。
func arrays() {
	a := [3]int{1, 2, 3}
	b := a // 複製一份
	b[0] = 100
	fmt.Println(a, b, len(a))
	c := [...]int{1, 2, 3} // ... 讓編譯器自己數長度
	fmt.Println(a == c)    // 同型別的陣列可以直接比較
}

// sliceWindow 示範 slice 透過窗戶改到的是底層陣列。
func sliceWindow() {
	arr := [5]int{10, 20, 30, 40, 50}
	s := arr[1:4] // 看見 20 30 40；cap 從起點算到陣列結尾
	fmt.Println(s, len(s), cap(s))
	s[0] = 99
	fmt.Println(arr)
}

// appendGrowth 示範 append：容量夠就原地放，不夠就換一個更大的底層陣列。
// 新容量是多少是 runtime 的實作細節，不同 Go 版本可能不同，程式不要依賴它。
func appendGrowth() {
	var s []int // nil slice：len 0、cap 0，可以直接 append
	fmt.Println(s == nil, len(s), cap(s))

	full := make([]int, 3)   // len 3、cap 3：已經滿了
	grown := append(full, 4) // 放不下，append 換了一個更大的新陣列
	grown[0] = 99
	fmt.Println(full, grown, cap(grown) > cap(full)) // full 沒被改到：兩者不再共用
}

// sharedBacking 示範兩個 slice 共用底層陣列時，append 可能蓋掉別人的資料。
func sharedBacking() {
	s := []int{1, 2, 3, 4}
	head := s[:2]           // len 2、cap 4：後面還有空間
	head = append(head, 99) // 容量夠，直接寫進底層陣列的第 3 格
	fmt.Println(s, head)    // s 的 3 被蓋掉了

	t := []int{1, 2, 3, 4}
	safe := t[:2:2]         // 完整切片運算式：把 cap 也限制成 2
	safe = append(safe, 99) // 容量不夠，換新陣列，不影響 t
	fmt.Println(t, safe)
}

// copyAndHelpers 示範 copy 與 Go 1.21 加入的 slices 套件。
func copyAndHelpers() {
	src := []int{3, 1, 2}
	dst := make([]int, len(src))
	n := copy(dst, src) // 回傳實際複製的個數
	slices.Sort(dst)
	fmt.Println(n, src, dst)
	fmt.Println(slices.Contains(dst, 2), slices.Index(dst, 3), slices.Max(dst))
}

// twoD 示範二維 slice：每一列可以不一樣長。
func twoD() {
	grid := make([][]int, 3)
	for i := range grid {
		grid[i] = make([]int, i+1)
		for j := range grid[i] {
			grid[i][j] = i + j
		}
	}
	fmt.Println(grid)
}

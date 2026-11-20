package main

import "fmt"

// sharedBacking 示範兩個 slice 共用底層陣列時，append 可能蓋掉別人的資料。
func sharedBacking() {
	s := []int{1, 2, 3, 4}
	part := s[:2]           // len 2、cap 4：後面還有空間
	part = append(part, 99) // 容量夠，直接寫進底層陣列的索引 2
	fmt.Println(s, part)
}

// fullSliceExpr 示範 s[low:high:max] 同時限制容量，之後的 append 一定會換到新陣列。
func fullSliceExpr() {
	s := []int{1, 2, 3, 4}
	part := s[:2:2] // len 2、cap 2
	part = append(part, 99)
	fmt.Println(s, part)
}

// raise 修改元素：呼叫端看得到，因為兩邊看向同一個陣列。
func raise(prices []int) {
	for i := range prices {
		prices[i] += 5
	}
}

// addItem 想用 append 加一項，但只改到自己那份 slice，呼叫端的長度不會變。
func addItem(cart []string, item string) {
	cart = append(cart, item)
}

// withItem 是慣用的寫法：回傳新的 slice，由呼叫端接住。
func withItem(cart []string, item string) []string {
	return append(cart, item)
}

// passSlice 示範把 slice 傳給函式：改元素看得到，append 看不到。
func passSlice() {
	prices := []int{85, 90}
	raise(prices)
	fmt.Println(prices)

	cart := []string{"排骨"}
	addItem(cart, "滷蛋")
	fmt.Println(cart, len(cart))
	cart = withItem(cart, "滷蛋")
	fmt.Println(cart, len(cart))
}

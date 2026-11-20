package main

import "fmt"

// appendBasics 示範 append 回傳新的 slice，一定要接住回傳值。
func appendBasics() {
	var cart []string // nil slice，可以直接 append
	cart = append(cart, "排骨")
	cart = append(cart, "雞腿", "滷蛋")
	drinks := []string{"紅茶", "味噌湯"}
	cart = append(cart, drinks...)
	fmt.Println(cart, len(cart))
}

// appendGrowth 示範容量不夠時，append 會換一個更大的底層陣列，新舊 slice 從此分開；
// 容量夠時則直接寫進原本的陣列，新舊 slice 共用同一塊。
func appendGrowth() {
	full := make([]int, 3) // len 3、cap 3：已經滿了
	grown := append(full, 4)
	grown[0] = 99
	fmt.Println(full, grown, cap(grown) > cap(full))

	roomy := make([]int, 3, 10) // len 3、cap 10：後面還有空間
	same := append(roomy, 4)
	same[0] = 99
	fmt.Println(roomy, same, &roomy[0] == &same[0])
}

// countMoves 依序 append n 個數字，回傳底層陣列換了幾次。
func countMoves(s []int, n int) int {
	moves := 0
	for i := range n {
		before := cap(s)
		s = append(s, i)
		if cap(s) != before {
			moves++
		}
	}
	return moves
}

// appendPrealloc 示範事先知道數量時，用 make 保留容量，append 就不必一直換陣列。
func appendPrealloc() {
	fmt.Println("沒有保留容量，換過陣列：", countMoves(nil, 1000) > 0)
	fmt.Println("保留 1000 格，換陣列的次數：", countMoves(make([]int, 0, 1000), 1000))
}

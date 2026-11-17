package main

import "fmt"

// precedence 示範運算子優先順序：乘除和位移比加減先算，比較又比 && 先算。
func precedence() {
	fmt.Println(2+3*4, (2+3)*4)
	fmt.Println(1 + 2<<1) // << 和 * 同一級，比 + 先算
	n := 6
	fmt.Println(n%2 == 0 && n > 5)
	fmt.Println(n&1 == 0) // & 也比 == 先算
}

// expensive 回傳平均單價是否超過 100 元。
// && 左邊是 false 時，右邊根本不會執行，所以 qty 是 0 也不會除以 0。
func expensive(total, qty int) bool {
	return qty > 0 && total/qty > 100
}

// shortCircuit 示範 && 的短路求值。
func shortCircuit() {
	fmt.Println(expensive(450, 3), expensive(0, 0))
}

// incDec 示範 ++、-- 和複合指定。
func incDec() {
	queue := 5
	queue++    // 來了一位
	queue--    // 叫走一位
	queue += 3 // 一次來三位
	price := 85
	price *= 2
	fmt.Println(queue, price)
}

// bitFlags 示範用位元運算組合、檢查、移除、切換旗標。%04b 以 4 位二進位印出。
func bitFlags() {
	order := Spicy | ExtraRice
	fmt.Printf("%04b\n", order)
	fmt.Println(order&Spicy != 0, order&NoOnion != 0)
	order |= Takeout
	order &^= Spicy
	fmt.Printf("%04b\n", order)
	order ^= NoOnion
	fmt.Printf("%04b\n", order)
}

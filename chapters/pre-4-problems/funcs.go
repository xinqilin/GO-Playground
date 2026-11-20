package main

import "fmt"

// greet 有一個參數 name，沒有回傳值：它只負責印出一行字。
func greet(name string) {
	fmt.Println("早安，" + name)
}

func greetDemo() {
	greet("小明")
	greet("阿美")
}

// subtotal 有兩個參數，回傳一個 int：單價乘以數量。
func subtotal(price, qty int) int {
	return price * qty
}

func subtotalDemo() {
	a := subtotal(85, 3)
	b := subtotal(90, 2)
	fmt.Println(a, b, a+b)
}

// memberPrice 遇到 return 就立刻結束，把結果交回去；後面的程式不會執行。
func memberPrice(price int, member bool) int {
	if !member {
		return price
	}
	return price * 9 / 10
}

func memberPriceDemo() {
	fmt.Println(memberPrice(100, false))
	fmt.Println(memberPrice(100, true))
	fmt.Println(memberPrice(85, true))
}

// hoursMinutes 一次回傳兩個值。
func hoursMinutes(minutes int) (int, int) {
	return minutes / 60, minutes % 60
}

func hoursMinutesDemo() {
	h, m := hoursMinutes(135)
	fmt.Println(h, "小時", m, "分")
}

// addOne 改的是自己的參數 n，呼叫的人手上的變數不受影響。
func addOne(n int) {
	n = n + 1
	fmt.Println("函式裡：", n)
}

func copyDemo() {
	count := 5
	addOne(count)
	fmt.Println("函式外：", count)
}

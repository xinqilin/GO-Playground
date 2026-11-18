package main

import "fmt"

// compare 比較兩個值，結果是 true 或 false。
func compare() {
	price := 85
	fmt.Println(price > 80)
	fmt.Println(price == 90)
	fmt.Println(price != 90)
	fmt.Println(price <= 85)
	fmt.Println("排骨" == "排骨")
}

// logic 用 ||（或）、&&（而且）、!（不是）組合條件。
func logic() {
	age := 70
	member := false
	total := 120
	fmt.Println(age < 12 || age >= 65)
	fmt.Println(member && total >= 100)
	fmt.Println(!member)
}

// traceMe 是追蹤表的示範：先自己追蹤，再執行對答案。
func traceMe() {
	x := 3
	y := x * 2
	x = y + 1
	y = x - y
	fmt.Println(x, y)
}

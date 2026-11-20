package main

import "fmt"

// arithmetic 示範五種算術。兩個整數相除，結果還是整數，小數部分直接丟掉。
func arithmetic() {
	fmt.Println(7 + 2)
	fmt.Println(7 - 2)
	fmt.Println(7 * 2)
	fmt.Println(7 / 2)
	fmt.Println(7 % 2)
	fmt.Println(7.0 / 2)
}

// splitBill 把 100 元平分給 3 個人：/ 算每人幾元，% 算剩幾元。
func splitBill() {
	total := 100
	people := 3
	fmt.Println("每人", total/people, "元")
	fmt.Println("剩下", total%people, "元")
}

// duration 把 135 分鐘換成幾小時幾分。
func duration() {
	minutes := 135
	fmt.Println(minutes/60, "小時", minutes%60, "分")
}

// discount 打九折：整數和小數不能直接相乘，要先把整數轉成小數。
func discount() {
	price := 85
	rate := 0.9
	fmt.Println(float64(price) * rate)
	fmt.Println(int(float64(price) * rate))
}

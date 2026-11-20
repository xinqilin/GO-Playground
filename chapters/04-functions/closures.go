package main

import "fmt"

// counter 回傳一個 closure。closure 會記住它用到的外層變數 n，每次呼叫都在同一個 n 上累加。
func counter() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}

// closureCounter 示範每個 closure 有自己的一份 n。
func closureCounter() {
	next := counter()
	fmt.Println(next(), next(), next())
	other := counter()
	fmt.Println(other(), next())
}

// closureShared 示範 closure 抓住的是變數本身，不是當時的值：之後外面改了，closure 也看得到。
func closureShared() {
	price := 85
	show := func() { fmt.Println("目前價格：", price) }
	show()
	price = 90
	show()
}

// greeter 是一個「產生函式的函式」：每次呼叫都做出一個帶著不同 prefix 的新函式。
func greeter(prefix string) func(string) string {
	return func(name string) string {
		return prefix + name
	}
}

// closureFactory 示範用 closure 預先設定好一部分參數。
func closureFactory() {
	welcome := greeter("歡迎光臨，")
	call := greeter("請取餐：")
	fmt.Println(welcome("小明"))
	fmt.Println(call("A07"))
}

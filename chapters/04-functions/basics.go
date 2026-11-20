package main

import "fmt"

// bentoTotal 計算便當的總價。參數寫成「名稱 型別」，相同型別的相鄰參數可以合併成 qty, price int。
func bentoTotal(qty, price int) int {
	return qty * price
}

// callBasics 示範呼叫函式。
func callBasics() {
	fmt.Println(bentoTotal(3, 85))
	fmt.Println(bentoTotal(2, 90) + bentoTotal(1, 25))
}

// addOne 改的是自己的參數。參數是呼叫時複製進來的一份，改了也不會影響呼叫端。
func addOne(n int) {
	n++
	fmt.Println("函式裡：", n)
}

// passByValue 示範 Go 的參數一律傳值。
func passByValue() {
	stock := 10
	addOne(stock)
	fmt.Println("呼叫端：", stock)
}

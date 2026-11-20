// 先修 2：值、變數與運算。
//
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
//
// 執行：go run ./chapters/pre-2-values
package main

import "fmt"

func main() {
	section("1. 值與型別")
	values()
	valueTypes()
	textVsNumber()

	section("2. 變數")
	variables()

	section("3. 改變變數的值")
	reassign()
	shortcuts()

	section("4. 交換兩個變數")
	wrongSwap()
	tempSwap()
	goSwap()

	section("5. 算術與整數除法")
	arithmetic()
	splitBill()
	duration()
	discount()

	section("6. 比較與邏輯")
	compare()
	logic()

	section("7. 追蹤表")
	traceMe()
}

func section(title string) {
	fmt.Println()
	fmt.Println("==", title, "==")
}

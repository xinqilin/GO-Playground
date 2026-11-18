// 先修 4：拆解問題：函式與清單。
//
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
//
// 執行：go run ./chapters/pre-4-problems
package main

import "fmt"

func main() {
	section("1. 函式")
	greetDemo()
	subtotalDemo()
	memberPriceDemo()
	hoursMinutesDemo()
	copyDemo()

	section("2. 清單")
	listBasics()
	stringList()
	emptyList()

	section("3. 走過清單")
	loopIndex()
	loopRange()

	section("4. 四個基本演算法")
	sumDemo()
	countDemo()
	maxDemo()
	wrongMaxDemo()
	searchDemo()

	section("5. 拆解問題")
	weeklyReport()

	section("6. 用 Println 除錯")
	debugDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

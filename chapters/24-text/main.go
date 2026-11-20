// 第 24 章：文字處理。
//
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
//
// 執行：go run ./chapters/24-text
package main

import "fmt"

func main() {
	section("1. 寫一個小解析器")
	parseDemo()
	numErrorDemo()

	section("2. regexp：比對與擷取")
	matchDemo()
	submatchDemo()

	section("3. regexp：取代與限制")
	replaceDemo()
	replaceFuncDemo()
	re2Demo()

	section("4. text/template")
	receiptDemo()
	execErrorDemo()

	section("5. 函式與組合範本")
	noticeDemo()

	section("6. html/template")
	escapeDemo()
	trustedDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

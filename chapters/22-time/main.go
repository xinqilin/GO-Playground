// 第 22 章：時間。
//
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
// 示範裡的時間都是寫死的，不用 time.Now()，輸出才會每次都一樣。
//
// 執行：go run ./chapters/22-time
package main

import "fmt"

func main() {
	section("1. Time 與 Duration")
	etaDemo()
	durationDemo()

	section("2. 格式化與解析")
	layoutDemo()
	parseDemo()

	section("3. 時區")
	zonesDemo()
	parseZoneDemo()

	section("4. 日曆的計算")
	calendarDemo()
	truncateDemo()
	dstDemo()

	section("5. 比較、單調時鐘與 Unix 時間")
	equalDemo()
	monotonicDemo()
	unixDemo()

	section("6. 計時器")
	afterFuncDemo()
	tickerDemo()

	section("7. 讓時間可以測試")
	checkoutDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

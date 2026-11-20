// 第 11 章：錯誤處理。
//
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
//
// 執行：go run ./chapters/11-errors
package main

import "fmt"

func main() {
	section("1. 錯誤是一個值")
	basicsDemo()

	section("2. sentinel error")
	sentinelDemo()

	section("3. 包裝：%w 與錯誤鏈")
	wrapDemo()
	wrapVsFormat()

	section("4. 自訂錯誤型別")
	asDemo()

	section("5. 一次回報多個錯誤")
	joinDemo()

	section("6. Close 也會失敗")
	closeDemo()

	section("7. 什麼時候用 panic")
	mustDemo()

	section("8. 分層處理錯誤")
	layersDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

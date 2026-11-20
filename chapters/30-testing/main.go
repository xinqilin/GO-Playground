// 第 30 章：測試。
//
// 這一章的主角是 parking 套件和它的測試檔（parking/*_test.go）。
// main 只是把 parking 跑一遍，讓你先看看被測試的程式在做什麼；
// main_test.go 用 Example 鎖住這些輸出，module_test.go 實際跑 testdata 裡會失敗的測試，
// 確認頁面上的失敗訊息是真的。
//
// 執行：go run ./chapters/30-testing
// 測試：go test ./chapters/30-testing/...
package main

import "fmt"

func main() {
	section("1. 停車費")
	feeDemo()

	section("2. 柵欄機與日報表")
	gateDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

// 第 31 章：效能與 profiling。
//
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，並放著頁面上用到的 benchmark。
// 時間會隨機器不同，所以 Example 只印結果和配置次數，不印時間。
//
// 執行：go run ./chapters/31-performance
// benchmark：go test -run='^$' -bench=. -benchmem ./chapters/31-performance
package main

import "fmt"

func main() {
	section("1. 先確認結果正確")
	revenueDemo()

	section("3. 數配置的次數")
	allocDemo()

	section("4. 逃逸分析")
	escapeDemo()

	section("5. 預先配置")
	preallocDemo()

	section("6. sync.Pool")
	poolDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

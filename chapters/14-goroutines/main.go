// 第 14 章：goroutine。
//
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
// 輸出會隨時間或機器而不同的示範放在子目錄：nowait、procs、crash，它們沒有測試。
//
// 執行：go run ./chapters/14-goroutines
package main

import "fmt"

func main() {
	section("2. 等它們做完：sync.WaitGroup")
	kitchenDemo()
	kitchenAddDone()
	argsNowDemo()

	section("3. 把結果交回來")
	sumDemo()
	storesDemo()

	section("4. goroutine 很便宜")
	spawnDemo()
	growDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

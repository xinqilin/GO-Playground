// 第 19 章：並行模式。
//
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
//
// 執行：go run ./chapters/19-concurrency-patterns
package main

import "fmt"

func main() {
	section("1. worker pool")
	poolDemo()

	section("2. semaphore")
	semaphoreDemo()

	section("3. 可以取消的管線")
	pipelineDemo()

	section("4. 結果的順序")
	orderDemo()

	section("5. 第一個錯誤就全部取消")
	groupDemo()

	section("6. 限制速率")
	tickerDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

// 第 16 章：select。
//
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
//
// 執行：go run ./chapters/16-select
package main

import "fmt"

func main() {
	section("1. 同時等好幾個 channel")
	readyDemo()
	fairDemo()

	section("2. 逾時")
	timeoutDemo()
	deadlineDemo()

	section("3. default：不等")
	tryReceiveDemo()
	dropDemo()

	section("4. for-select 與 nil channel")
	loopDemo()

	section("5. 叫送的一方停下來")
	findDemo()

	section("6. 優先順序")
	priorityDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

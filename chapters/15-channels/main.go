// 第 15 章：channel。
//
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
// 一定會當掉的死結示範放在 deadlock 子目錄，沒有測試。
//
// 執行：go run ./chapters/15-channels
package main

import "fmt"

func main() {
	section("1. 送與收：無緩衝的 channel")
	handoffDemo()
	doneDemo()

	section("2. 有緩衝的 channel")
	bufferDemo()

	section("3. close 與 range")
	rangeDemo()
	closedDemo()
	closePanicDemo()

	section("4. 方向")
	directionDemo()

	section("5. 管線")
	pipelineDemo()

	section("6. fan-out 與 fan-in")
	fanDemo()

	section("7. 死結與洩漏")
	findDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

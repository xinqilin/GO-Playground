// 第 13 章：疊代器。
//
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
//
// 執行：go run ./chapters/13-iterators
package main

import "fmt"

func main() {
	section("1. range 一個函式")
	rangeFuncDemo()
	byHand()

	section("2. yield 的約定：提早結束")
	earlyStopDemo()
	ignoreYieldDemo()

	section("3. 替自己的型別寫疊代器")
	seq2Demo()
	treeDemo()

	section("4. 標準庫的疊代器")
	stdlibDemo()

	section("5. 串成管線")
	pipelineDemo()
	lazyDemo()

	section("6. iter.Pull")
	pullBasics()
	pullDemo()

	section("7. 疊代器、slice 與 channel")
	reuseDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

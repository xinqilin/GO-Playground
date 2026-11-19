// 第 18 章：context。
//
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
//
// 執行：go run ./chapters/18-context
package main

import "fmt"

func main() {
	section("1. 為什麼需要 context")
	basicsDemo()

	section("2. 取消樹")
	treeDemo()

	section("3. 期限")
	errDemo()
	deadlineDemo()

	section("4. 取消的原因")
	causeDemo()

	section("5. 寫一個聽 ctx 的函式")
	listenDemo()
	tallyDemo()

	section("6. 帶值")
	valueDemo()

	section("7. WithoutCancel 與 AfterFunc")
	detachDemo()
	afterDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

// 第 04 章：函式。
//
// 函式是一台機器：放進參數、拿出結果。在 Go 裡，函式本身也是一種值，可以存在變數裡、傳來傳去。
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
//
// 執行：go run ./chapters/04-functions
package main

import "fmt"

func main() {
	section("1. 宣告與呼叫")
	callBasics()
	passByValue()

	section("2. 多個回傳值")
	multiReturn()
	valueAndError()
	namedResults()

	section("3. 可變參數")
	variadic()

	section("4. 函式是值")
	funcValues()
	funcTypes()
	sortWithFunc()

	section("5. closure")
	closureCounter()
	closureShared()
	closureFactory()

	section("6. 遞迴")
	recursion()

	section("7. defer")
	deferOrder()
	deferArgs()
	fmt.Println(deferResult())
	deferInLoop()
	deferInLoopFixed()

	section("8. panic 與 recover")
	panicRecover()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

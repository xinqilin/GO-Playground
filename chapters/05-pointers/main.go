// 第 05 章：指標。
//
// 指標是一張寫著「門牌」的紙條：拿著它就能找到、也能改動原本的變數。
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
// 指標的實際位址每次執行都不同，所以示範裡只比較位址，不印出來。
//
// 執行：go run ./chapters/05-pointers
package main

import "fmt"

func main() {
	section("1. 取位址與解參考")
	addressAndDeref()
	modifyThroughPointer()
	aliasing()

	section("2. 傳值與傳指標")
	passValueVsPointer()
	structCopy()
	shallowCopy()

	section("3. 指向 struct 的指標")
	structPointer()
	newVsLiteral()
	newWithValue()

	section("4. nil 指標")
	nilPointer()
	optionalField()

	section("5. 回傳區域變數的指標")
	returnLocalPointer()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

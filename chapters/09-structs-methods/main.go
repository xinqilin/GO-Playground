// 第 09 章：自訂型別、struct 與方法。
//
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
//
// 執行：go run ./chapters/09-structs-methods
package main

import "fmt"

func main() {
	section("1. 自訂型別")
	definedTypes()
	aliasDemo()

	section("2. struct")
	structLiterals()
	structPointer()
	anonymousStruct()

	section("3. 複製與比較")
	copyStruct()
	shallowCopy()
	compareStructs()

	section("4. 值 receiver 與指標 receiver")
	receiverDemo()
	autoAddress()
	methodSetDemo()

	section("5. 建構函式與零值")
	zeroValueDemo()
	constructorDemo()
	tallyDemo()

	section("6. 嵌入")
	embedDemo()
	promotedKeys()
	shadowDemo()
	notInheritance()

	section("7. struct tag")
	tagsDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

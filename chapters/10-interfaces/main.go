// 第 10 章：interface。
//
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
//
// 執行：go run ./chapters/10-interfaces
package main

import "fmt"

func main() {
	section("1. interface 是一份方法清單")
	implicitDemo()

	section("2. interface 值")
	interfaceValue()
	pointerInInterface()
	compareInterfaces()

	section("3. type assertion")
	assertDemo()
	assertPanic()
	optionalDemo()

	section("4. type switch 與 any")
	typeSwitchDemo()

	section("5. nil interface 陷阱")
	nilTrapDemo()

	section("6. 標準庫裡的 interface")
	writerDemo()
	sortDemo()

	section("7. 設計 interface")
	embedInterfaceDemo()
	remindDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

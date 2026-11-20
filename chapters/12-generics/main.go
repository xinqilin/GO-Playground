// 第 12 章：泛型。
//
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
//
// 執行：go run ./chapters/12-generics
package main

import "fmt"

func main() {
	section("1. 為什麼需要泛型")
	whyDemo()

	section("2. 型別參數與限制")
	paramsDemo()

	section("3. 型別集合與 ~")
	typeSetDemo()

	section("4. 型別推論")
	inferDemo()

	section("5. 泛型型別")
	setDemo()

	section("6. generic method")
	methodDemo()

	section("7. 什麼時候用泛型")
	whenDemo()
	kindDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

// 第 25 章：日誌。
//
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
//
// 執行：go run ./chapters/25-logging
package main

import "fmt"

func main() {
	section("1. log 套件")
	logDemo()

	section("2. slog：結構化的紀錄")
	textJSONDemo()

	section("3. 屬性與 With")
	attrsDemo()
	withDemo()

	section("4. 層級")
	levelsDemo()
	levelVarDemo()

	section("5. LogValuer")
	valuerDemo()
	lazyDemo()

	section("6. Handler")
	multiDemo()
	contextDemo()

	section("7. 測試日誌")
	checkoutDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

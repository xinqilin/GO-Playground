// 第 29 章：package 與 module。
//
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
// module_test.go 對 testdata 底下的幾個小 module 執行 go 指令，確認頁面上的指令輸出是真的。
//
// 執行：go run ./chapters/29-packages-modules
package main

import "fmt"

func main() {
	section("1. package 與 import")
	exportDemo()
	aliasDemo()

	section("2. 初始化的順序")
	initOrderDemo()

	section("5. 版本")
	goVersionDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

// 第 08 章：字串、byte 與 rune。
//
// 字串是一排唯讀的 byte，內容用 UTF-8 編碼：英文字佔 1 個 byte，中文字佔 3 個。
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
//
// 執行：go run ./chapters/08-strings-runes
package main

import "fmt"

func main() {
	section("1. 字串的基本操作")
	stringBasics()
	literals()

	section("2. UTF-8：byte 與 rune")
	bytesVsRunes()
	runeLiteral()
	cutInMiddle()

	section("3. range 字串")
	rangeString()
	rangeInvalid()

	section("4. 轉換")
	conversions()
	editRunes()
	numberToString()

	section("5. strings 套件")
	searchFuncs()
	splitJoin()
	trimCut()
	replaceFuncs()

	section("6. strings.Builder")
	builderDemo()

	section("7. strconv 與 unicode")
	parseNumbers()
	formatNumbers()
	classify()
	halfwidth()

	section("8. 中文的顯示寬度")
	alignedMenu()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

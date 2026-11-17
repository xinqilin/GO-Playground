// 第 02 章：輸入與輸出。
//
// 程式像一台果汁機：從 io.Reader 倒進原料，從 io.Writer 倒出結果。
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
// 讀標準輸入的完整程式在 sum 子目錄。
//
// 執行：go run ./chapters/02-input-output
package main

import "fmt"

func main() {
	section("1. Print 家族")
	printSpacing()
	sprintDemo()
	fprintDemo()

	section("2. 格式化動詞")
	verbsGeneral()
	verbsInt()
	verbsFloat()
	verbsString()

	section("3. 寬度與精度")
	widthBasics()
	fareTable()
	widthCJK()

	section("4. Stringer")
	stringerDemo()

	section("5. 標準輸出與標準錯誤")
	streams()

	section("6. bufio.Scanner")
	scanLines()
	scanWords()
	scanTooLong()
	scanBigBuffer()
	sscanDemo()

	section("7. io.Reader 與 io.Writer")
	numberLinesDemo()
	bufferedWriter()
	copyDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

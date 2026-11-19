// 第 20 章：檔案與 I/O。
//
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
// 示範只在自己開的暫存目錄裡讀寫，結束時刪掉，不會碰到其他檔案。
//
// 執行：go run ./chapters/20-files
package main

import "fmt"

func main() {
	section("1. 一次讀寫整個檔案")
	wholeDemo()
	missingDemo()

	section("2. 開檔與關檔")
	openDemo()

	section("3. 串流：一塊一塊處理")
	copyFileDemo()
	mergeDemo()
	previewDemo()

	section("4. 路徑與 os.Root")
	pathsDemo()
	traversalDemo()
	rootDemo()

	section("5. 走訪目錄")
	walkDemo()

	section("6. fs.FS 與 embed")
	fsDemo()
	embedDemo()

	section("7. 暫存檔與安全覆寫")
	atomicDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

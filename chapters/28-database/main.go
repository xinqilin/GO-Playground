// 第 28 章：資料庫。
//
// 這一章是獨立的 module（有自己的 go.mod），因為它需要第三方的 SQLite driver；手冊的其他部分只用標準庫。
// 每個示範都在暫存目錄開一個新的資料庫，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
//
// 執行（在 repo 根目錄）：go -C chapters/28-database run .
package main

import "fmt"

func main() {
	section("1. database/sql 與 driver")
	openDemo()

	section("2. 寫入：Exec")
	execDemo()

	section("3. 讀取：Query 與 QueryRow")
	queryDemo()

	section("4. 參數與 SQL injection")
	injectionDemo()

	section("5. NULL")
	nullDemo()

	section("6. 交易")
	txDemo()

	section("7. 連線池與 context")
	poolDemo()

	section("8. migration 與測試")
	migrateDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

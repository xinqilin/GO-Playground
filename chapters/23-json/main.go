// 第 23 章：JSON 與資料編碼。
//
// 本章用 Go 1.27 加入標準庫的 encoding/json/v2，最後一節再對照現有程式碼裡常見的 encoding/json（v1）。
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
//
// 執行：go run ./chapters/23-json
package main

import "fmt"

func main() {
	section("1. Marshal 與 Unmarshal")
	marshalDemo()
	unmarshalDemo()

	section("2. struct tag")
	tagsDemo()
	omitDemo()

	section("3. 解碼的規則與錯誤")
	anyDemo()
	strictDemo()
	decodeErrorsDemo()

	section("4. 結構未知的 JSON")
	lookupDemo()
	eventDemo()

	section("5. 自訂編碼")
	textDemo()
	uuidDemo()
	durationDemo()

	section("6. 串流")
	linesDemo()
	arrayDemo()
	encodeDemo()

	section("7. 現有程式碼裡的 encoding/json")
	v1v2Demo()
	v1APIDemo()

	section("8. encoding/csv")
	csvDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

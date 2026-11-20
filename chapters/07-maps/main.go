// 第 07 章：map。
//
// map 是一排有名牌的置物櫃：用 key 直接找到櫃子，不用從頭找。
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
// map 的走訪順序每次都可能不同，所以需要印出多個 key 的示範都先排序，或交給 fmt（它會排序 key）。
//
// 執行：go run ./chapters/07-maps
package main

import "fmt"

func main() {
	section("1. 建立與基本操作")
	mapBasics()
	makeMap()
	nilMap()

	section("2. 讀取：零值與 comma-ok")
	commaOK()
	countOrders()
	groupOrders()

	section("3. 走訪與刪除")
	rangeMap()
	sortedRange()
	sortByValue()
	deleteWhileRange()

	section("4. key 的型別")
	structKey()
	nanKey()

	section("5. 值是 struct")
	structValue()
	pointerValue()

	section("6. map 是參考")
	mapIsReference()

	section("7. 當成集合")
	mapAsSet()

	section("8. 預先配置")
	indexByName()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

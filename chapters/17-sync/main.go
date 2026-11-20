// 第 17 章：sync 與 data race。
//
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
//
// 執行：go run ./chapters/17-sync
// 故意出錯的版本：go run -race ./chapters/17-sync/race、go run ./chapters/17-sync/mapcrash
package main

import "fmt"

func main() {
	section("1. data race")
	lostUpdateDemo()

	section("2. Mutex")
	stockDemo()
	oversellDemo()

	section("3. RWMutex")
	menuDemo()

	section("4. 只做一次")
	onceDemo()
	onceErrDemo()

	section("5. atomic")
	ticketDemo()
	maxDemo()
	noticeDemo()

	section("6. sync.Map 與 sync.Pool")
	seenDemo()
	poolDemo()

	section("7. happens-before")
	orderDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

// 先修 3：判斷與重複。
//
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
//
// 執行：go run ./chapters/pre-3-decisions
package main

import "fmt"

func main() {
	section("1. if 與 else")
	passFailDemo()
	umbrellaDemo()

	section("2. else if")
	gradeDemo()
	wrongGradeDemo()

	section("3. 組合條件")
	ticketDemo()

	section("4. for")
	countdown()
	laps()

	section("5. 計數器與累加器")
	sumTo100()
	countMultiples()
	savings()

	section("6. 只寫條件的 for、break 與 continue")
	bounce()
	firstMultiple()
	openDays()

	section("7. 巢狀迴圈")
	triangle()
	seats()

	section("8. 無窮迴圈")
	neverTen()
}

func section(title string) {
	fmt.Println()
	fmt.Println("==", title, "==")
}

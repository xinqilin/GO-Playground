// 第 03 章：流程控制。
//
// Go 的流程控制只有 if、for、switch 三種，再加上 break、continue 和很少用到的 goto。
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
//
// 執行：go run ./chapters/03-control-flow
package main

import "fmt"

func main() {
	section("1. if")
	ifElse()
	ifInit()
	guardClause()
	shadowing()

	section("2. for 的四種寫法")
	forClassic()
	forCondition()
	forInfinite()
	forRangeInt()
	forContinue()

	section("3. range")
	rangeCopy()
	rangeArrayCopy()
	rangeString()
	rangeMap()
	rangeChannel()
	rangeFunc()
	rangeOnce()
	loopVarPerIteration()

	section("4. break、continue 與標籤")
	labeledBreak()
	labeledContinue()
	breakInSwitch()

	section("5. switch")
	switchValue()
	switchInit()
	switchTrue()
	switchOrder()
	switchFallthrough()
	typeSwitch()

	section("常見錯誤")
	shadowPitfall()
	fmt.Println(reverse([]string{"A01", "A02", "A03"}))
}

func section(title string) {
	fmt.Println("==", title, "==")
}

// 第 01 章：變數、常數與型別。
//
// 變數是一個貼了名字、只能放某一種東西的盒子；沒放東西時，盒子裡是那個型別的「零值」。
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
//
// 執行：go run ./chapters/01-variables
package main

import "fmt"

func main() {
	section("1. 宣告變數")
	declareVar()
	declareShort()
	declareMulti()
	redeclare()
	typeInference()
	packageLevel()

	section("2. 零值")
	zeroValues()
	zeroUsable()

	section("3. 整數")
	intRanges()
	overflowWrap()
	integerDivision()
	divideByZero()

	section("4. 浮點數")
	floatError()
	floatCompare()
	moneyInTenths()
	floatSpecial()

	section("5. 型別轉換")
	convertNumbers()
	convertStrings()

	section("6. 常數與 iota")
	constBasics()
	untypedConst()
	iotaWeekday()
	iotaPatterns()

	section("7. 作用域")
	scopeLevels()
	shadowPackage()

	section("8. 運算子")
	precedence()
	shortCircuit()
	incDec()
	bitFlags()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

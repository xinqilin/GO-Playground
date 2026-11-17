// 先修 1：程式怎麼跑起來。
//
// 程式是一份給電腦的食譜：電腦從上到下，一行一行照著做。
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
//
// 執行：go run ./chapters/pre-1-running
package main

import "fmt"

func main() {
	hello()
	cookRice()
	calculate()
	wrongAverage()
	rightAverage()
}

// hello 印出一行字。雙引號裡的文字會原樣印出來。
func hello() {
	fmt.Println("你好，Go")
}

// cookRice 依序印出三個步驟。程式從上到下一行一行執行，所以印出來的順序就是寫的順序。
func cookRice() {
	fmt.Println("1. 量米")
	fmt.Println("2. 洗米")
	fmt.Println("3. 按下煮飯鍵")
}

// calculate 讓電腦幫忙算數學。沒有雙引號的 85 * 3 是一個算式，會先算出結果再印出來。
func calculate() {
	fmt.Println(85 * 3)
	fmt.Println("85 * 3")
	fmt.Println("三個排骨便當：", 85*3, "元")
}

// wrongAverage 想算 80 和 90 的平均，但忘了先加再除：乘除會比加減先算。
// 程式能跑，結果卻是錯的。
func wrongAverage() {
	fmt.Println(80 + 90/2)
}

// rightAverage 用括號讓加法先算。
func rightAverage() {
	fmt.Println((80 + 90) / 2)
}

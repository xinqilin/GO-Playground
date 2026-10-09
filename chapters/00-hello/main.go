// 第 00 章：第一支 Go 程式。
//
// 一支可以執行的 Go 程式需要兩樣東西：package main 和 func main()。
// 執行：go run ./chapters/00-hello
package main

// import 把標準庫的套件借進來用。fmt 負責格式化輸入與輸出。
import "fmt"

// main 是程式的入口：它跑完，整支程式就結束。
func main() {
	fmt.Println("Hello, 世界")

	// Println 可以一次印多個值，值與值之間自動補一個空白。
	fmt.Println("整數", 3, "小數", 3.14, "布林", true)

	// 單引號是 rune（一個 Unicode 字元），印出來的是它的編號，第 08 章細講。
	fmt.Println('a', '世')
}

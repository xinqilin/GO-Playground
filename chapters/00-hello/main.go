// 第 00 章：第一支 Go 程式。
//
// 一支可以執行的 Go 程式需要 package main 和 func main()。
// 這一章的程式分成 main.go 和 greet.go 兩個檔案，它們屬於同一個 package。
//
// 執行：go run ./chapters/00-hello
package main

import "fmt"

// main 是程式的入口：它跑完，整支程式就結束。
func main() {
	fmt.Println("哈囉，Go！")

	// Println 可以一次印好幾個值，值和值之間自動補一個空白，最後換行。
	fmt.Println("便當", 3, "個，共", 255, "元，已付款：", true)

	// greeting 寫在 greet.go，同一個 package 裡的函式可以直接呼叫。
	fmt.Println(greeting("小明"))
}

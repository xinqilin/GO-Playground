// crash 示範執行時才發生的錯誤：程式可以編譯，但跑到一半當掉。
// 它故意會當掉，所以沒有測試。
//
// 執行：go run ./chapters/pre-1-running/crash
package main

import "fmt"

func main() {
	total := 300
	people := 0
	fmt.Println("開始分帳")
	fmt.Println(total / people)
	fmt.Println("這一行不會執行")
}

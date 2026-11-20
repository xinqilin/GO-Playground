package main

import "fmt"

// 平台每台機器不同，所以只比對 describe 的輸出。
func Example_describe() {
	fmt.Println(describe())
	// Output: 版本 dev，免費版
}

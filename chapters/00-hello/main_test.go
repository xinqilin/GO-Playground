package main

import "fmt"

// Example 執行 main 並比對輸出；go test 會檢查下方的 Output 註解。
func Example() {
	main()
	// Output:
	// 哈囉，Go！
	// 便當 3 個，共 255 元，已付款： true
	// 哈囉，小明！
}

func Example_greeting() {
	fmt.Println(greeting("阿美"))
	// Output:
	// 哈囉，阿美！
}

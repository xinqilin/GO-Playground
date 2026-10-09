package main

// init 的輸出發生在測試開始之前，不在 Example 的比對範圍內。
func Example() {
	main()
	// Output:
	// 哈囉，Gopher
	// 哈囉，陌生人
	// 嗨，Go
	// 1 2
}

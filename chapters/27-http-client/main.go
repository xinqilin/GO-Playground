// 第 27 章：HTTP 用戶端。
//
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
// 示範的 client 不經過網路，而是把請求直接交給 fakeShop 這個假的店家（見 transport.go）。
//
// 執行：go run ./chapters/27-http-client
package main

import "fmt"

func main() {
	section("1. 發出請求、讀回應")
	getDemo()

	section("3. 網址與 JSON")
	urlDemo()
	jsonDemo()

	section("5. 重試與退避")
	retryDemo()

	section("6. RoundTripper")
	transportDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

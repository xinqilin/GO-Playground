// crash 示範沒有被 recover 的 panic：程式會印出訊息和呼叫路徑，然後以結束碼 2 結束。
// 它故意會當掉，所以沒有測試。
//
// 執行：go run ./chapters/04-functions/crash
package main

import "fmt"

func take(stock, n int) int {
	if n > stock {
		panic(fmt.Sprintf("庫存不足：剩 %d 個，要拿 %d 個", stock, n))
	}
	return stock - n
}

func main() {
	defer fmt.Println("main 的 defer 還是會執行")
	fmt.Println(take(5, 2))
	fmt.Println(take(3, 10))
	fmt.Println("這一行不會執行")
}

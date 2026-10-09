// 故意寫錯的計數器：1000 個 goroutine 同時做 count++，沒有任何保護。
//
// count++ 其實是「讀出來、加一、寫回去」三個步驟，兩個 goroutine 交錯執行就會少算。
// 用 race detector 執行，會看到 WARNING: DATA RACE：
//
//	go run -race ./chapters/17-sync/race
package main

import (
	"fmt"
	"sync"
)

func main() {
	count := 0
	var wg sync.WaitGroup
	for range 1000 {
		wg.Go(func() { count++ }) // data race：多個 goroutine 同時讀寫 count
	}
	wg.Wait()
	fmt.Println("count =", count, "（應該是 1000）")
}

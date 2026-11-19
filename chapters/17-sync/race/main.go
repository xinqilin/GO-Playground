// 故意寫錯的程式：10 個 goroutine 同時對 served 做 served++，沒有任何保護。
//
// served++ 是「讀出來、加一、寫回去」三步，兩個 goroutine 交錯執行就會少算。
// 這個資料夾沒有測試，因為 go test -race 一定會抓到它。
//
//	go run ./chapters/17-sync/race        # 結果常常少於 10000，每次不一樣
//	go run -race ./chapters/17-sync/race  # race detector 印出 WARNING: DATA RACE
package main

import (
	"fmt"
	"sync"
)

func main() {
	served := 0
	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			for range 1000 {
				served++ // data race：好幾個 goroutine 同時讀寫 served
			}
		})
	}
	wg.Wait()
	fmt.Println("賣出", served, "份（應該是 10000）")
}

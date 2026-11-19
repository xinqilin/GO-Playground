// 故意寫錯的程式：好幾個 goroutine 同時寫同一個 map，沒有上鎖。
//
// 執行時期偵測到同時寫入時，會直接讓整個程式結束（不需要 -race，但不保證每次都抓得到）。
// 這個資料夾沒有測試，因為它的目的就是當掉。
//
//	go run ./chapters/17-sync/mapcrash
package main

import (
	"fmt"
	"sync"
)

func main() {
	orders := make(map[int]int)
	var wg sync.WaitGroup
	for shop := range 10 {
		wg.Go(func() {
			for n := range 1000 {
				orders[shop*1000+n] = n // 好幾個 goroutine 同時寫 orders
			}
		})
	}
	wg.Wait()
	fmt.Println("共", len(orders), "筆")
}

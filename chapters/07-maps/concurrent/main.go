// concurrent 示範兩個 goroutine 同時寫同一個 map。runtime 偵測到時會直接結束程式，
// 這不是 panic，recover 也接不住。它故意會當掉，所以沒有測試。
//
// 執行：go run ./chapters/07-maps/concurrent
package main

import (
	"fmt"
	"sync"
)

func main() {
	counts := make(map[int]int)
	var wg sync.WaitGroup
	for g := range 2 {
		wg.Go(func() {
			for i := range 1_000_000 {
				counts[i%100] = g
			}
		})
	}
	wg.Wait()
	fmt.Println("沒有被偵測到，這次運氣好", len(counts))
}

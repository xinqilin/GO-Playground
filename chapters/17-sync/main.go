// 第 17 章：sync 與 data race。
//
// 多個 goroutine 同時改同一個變數，就像好幾個人同時在同一格寫字，結果會亂掉。
// Mutex 是一把鑰匙，一次只讓一個人進去；atomic 是不可分割的小動作。
// 執行：go run ./chapters/17-sync
// 故意出錯的版本：go run -race ./chapters/17-sync/race
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Counter 可以給多個 goroutine 同時使用：mu 保護 counts。
// 含有 Mutex 的 struct 不能複製，所以都用 *Counter 傳遞。
type Counter struct {
	mu     sync.Mutex
	counts map[string]int
}

// NewCounter 建立一個可以直接使用的 Counter。
func NewCounter() *Counter {
	return &Counter{counts: make(map[string]int)}
}

// Inc 把 key 的次數加 1。
func (c *Counter) Inc(key string) {
	c.mu.Lock()
	defer c.mu.Unlock() // 用 defer 解鎖，中途 return 或 panic 也不會忘記
	c.counts[key]++
}

// Get 讀取 key 的次數。讀也要上鎖，否則可能讀到寫到一半的資料。
func (c *Counter) Get(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.counts[key]
}

func main() {
	var wg sync.WaitGroup

	fmt.Println("== Mutex：一次只讓一個人進去 ==")
	c := NewCounter()
	for range 1000 {
		wg.Go(func() { c.Inc("visits") })
	}
	wg.Wait()
	fmt.Println(c.Get("visits"))

	fmt.Println("== atomic：單一數字的計數 ==")
	var hits atomic.Int64
	for range 1000 {
		wg.Go(func() { hits.Add(1) })
	}
	wg.Wait()
	fmt.Println(hits.Load())

	fmt.Println("== OnceValue：只做一次 ==")
	loadConfig := sync.OnceValue(func() string {
		fmt.Println("讀取設定檔…") // 不管幾個 goroutine 呼叫，只會印一次
		return "port=8080"
	})
	for range 3 {
		wg.Go(func() { loadConfig() })
	}
	wg.Wait()
	fmt.Println(loadConfig())
}

package main

import (
	"sync"
	"testing"
)

func Example() {
	main()
	// Output:
	// == Mutex：一次只讓一個人進去 ==
	// 1000
	// == atomic：單一數字的計數 ==
	// 1000
	// == OnceValue：只做一次 ==
	// 讀取設定檔…
	// port=8080
}

// 用 go test -race 執行時，Counter 少了鎖就會被抓到。
func TestCounterConcurrent(t *testing.T) {
	c := NewCounter()
	var wg sync.WaitGroup
	for i := range 100 {
		wg.Go(func() {
			c.Inc("a")
			if i%2 == 0 {
				c.Inc("b")
			}
		})
	}
	wg.Wait()
	if got := c.Get("a"); got != 100 {
		t.Errorf("a = %d, want 100", got)
	}
	if got := c.Get("b"); got != 50 {
		t.Errorf("b = %d, want 50", got)
	}
}

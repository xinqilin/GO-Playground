package main

import (
	"fmt"
	"strings"
	"sync"
)

// prepare 模擬廚房做一份餐。
func prepare(order string) string {
	return order + "好了"
}

// kitchenDemo：每一份餐交給一個 goroutine，wg.Go 開 goroutine 並登記，wg.Wait 等全部做完。
func kitchenDemo() {
	orders := []string{"排骨飯", "雞腿飯", "魚排飯"}
	done := make([]string, len(orders))
	var wg sync.WaitGroup
	for i, o := range orders {
		wg.Go(func() {
			done[i] = prepare(o)
		})
	}
	wg.Wait()
	fmt.Println(strings.Join(done, "、"))
}

// kitchenAddDone 是 Go 1.25 以前的寫法：先 Add(1) 再 go，goroutine 結束時 Done。
func kitchenAddDone() {
	orders := []string{"排骨飯", "雞腿飯", "魚排飯"}
	done := make([]string, len(orders))
	var wg sync.WaitGroup
	for i, o := range orders {
		wg.Add(1)
		go func() {
			defer wg.Done()
			done[i] = prepare(o)
		}()
	}
	wg.Wait()
	fmt.Println(strings.Join(done, "、"))
}

// argsNowDemo：go 後面的函式參數，在 go 這一行就求值了。
func argsNowDemo() {
	var wg sync.WaitGroup
	order := "排骨飯"
	wg.Add(1)
	go func(o string) {
		defer wg.Done()
		fmt.Println("廚房收到", o)
	}(order)
	order = "雞腿飯"
	wg.Wait()
	fmt.Println("櫃台現在的單", order)
}

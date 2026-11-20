package main

import (
	"context"
	"fmt"
	"sync"
)

// cookAll 用 cooks 個廚師（worker）處理 orders。
// 派工的一方把「第幾張單」送進 jobs，每個廚師從同一個 jobs 拿單，做好寫進自己那一格。
// ctx 被取消時不再派新的單，已經拿到單的廚師做完手上這張就結束。
func cookAll(ctx context.Context, orders []string, cooks int, cook func(string) string) []string {
	jobs := make(chan int)
	out := make([]string, len(orders))
	var wg sync.WaitGroup
	for range cooks {
		wg.Go(func() {
			for i := range jobs {
				out[i] = cook(orders[i])
			}
		})
	}

feed:
	for i := range orders {
		select {
		case jobs <- i:
		case <-ctx.Done():
			break feed
		}
	}
	close(jobs)
	wg.Wait()
	return out
}

func poolDemo() {
	orders := []string{"排骨飯", "雞腿飯", "魚排飯", "滷肉飯", "控肉飯"}
	done := cookAll(context.Background(), orders, 2, func(o string) string {
		return o + " 好了"
	})
	for _, d := range done {
		fmt.Println(d)
	}
}

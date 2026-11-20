package main

import (
	"context"
	"fmt"
)

// prepare 在另一個 goroutine 裡依序送出 items，送完就 close。
// 每次送的時候也同時等 ctx.Done()：ctx 被取消就不再送，直接結束。
func prepare(ctx context.Context, items ...string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for _, it := range items {
			select {
			case out <- it:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// listenDemo：拿到雞腿飯就不要了。cancel 通知 prepare 停下來，不用把剩下的收完。
func listenDemo() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for o := range prepare(ctx, "排骨飯", "雞腿飯", "魚排飯", "滷肉飯") {
		fmt.Println(o)
		if o == "雞腿飯" {
			break
		}
	}
}

// tally 加總一大串發票金額。這種一直在算、不會停下來等的工作，
// 沒有 channel 可以 select，就每做一段檢查一次 ctx.Err()。
func tally(ctx context.Context, amounts []int) (int, error) {
	sum := 0
	for i, a := range amounts {
		if i%1000 == 0 {
			if err := ctx.Err(); err != nil {
				return 0, err
			}
		}
		sum += a
	}
	return sum, nil
}

func tallyDemo() {
	amounts := make([]int, 5000)
	for i := range amounts {
		amounts[i] = i + 1
	}
	fmt.Println(tally(context.Background(), amounts))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	fmt.Println(tally(ctx, amounts))
}

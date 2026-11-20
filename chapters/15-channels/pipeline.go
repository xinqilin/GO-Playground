package main

import (
	"fmt"
	"slices"
	"sync"
)

// cook 是管線的一站：從 in 收訂單，做好之後送到 out；in 被 close 之後，它也 close out。
func cook(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for o := range in {
			out <- o + "（熱）"
		}
	}()
	return out
}

// pack 是下一站：把做好的餐裝進外帶袋。
func pack(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for m := range in {
			out <- "外帶袋[" + m + "]"
		}
	}()
	return out
}

// pipelineDemo：三站各在自己的 goroutine 裡，同時處理不同的訂單。
func pipelineDemo() {
	for bag := range pack(cook(orders("排骨飯", "雞腿飯"))) {
		fmt.Println(bag)
	}
}

// kitchen 開 cooks 位廚師從同一個 orders channel 搶單（fan-out），
// 做好的全部送到同一個 out channel（fan-in）。所有廚師都結束之後才 close out。
func kitchen(in <-chan string, cooks int) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup
	for range cooks {
		wg.Go(func() {
			for o := range in {
				out <- o + "好了"
			}
		})
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// fanDemo：哪位廚師做哪一份、哪一份先好，每次都不一樣，所以收完之後排序再印。
func fanDemo() {
	var done []string
	for m := range kitchen(orders("排骨飯", "雞腿飯", "魚排飯", "滷肉飯", "控肉飯"), 3) {
		done = append(done, m)
	}
	slices.Sort(done)
	fmt.Println(len(done), done)
}

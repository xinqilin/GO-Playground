// 第 19 章：並行模式。
//
// 把前幾章的零件組起來：worker pool 像固定人數的工班、fan-in 像好幾條輸送帶匯成一條、
// semaphore 像有限的號碼牌。
// 執行：go run ./chapters/19-concurrency-patterns
package main

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
	"time"
)

// result 是一個工作的結果，ID 是它在輸入 slice 裡的位置。
type result struct{ ID, Value int }

// workerPool 用 n 個 worker 處理 jobs，回傳依 ID 排序的結果。ctx 取消時所有 goroutine 都會收工。
func workerPool(ctx context.Context, jobs []int, n int, work func(int) int) []result {
	jobCh := make(chan int)
	resCh := make(chan result)

	var wg sync.WaitGroup
	for range n {
		wg.Go(func() {
			for id := range jobCh {
				select {
				case resCh <- result{ID: id, Value: work(jobs[id])}:
				case <-ctx.Done():
					return
				}
			}
		})
	}

	go func() { // 派工：送完或被取消就關閉 jobCh，worker 的 range 才會結束
		defer close(jobCh)
		for id := range jobs {
			select {
			case jobCh <- id:
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() { // 所有 worker 都結束後，關閉結果 channel
		wg.Wait()
		close(resCh)
	}()

	var out []result
	for r := range resCh {
		out = append(out, r)
	}
	slices.SortFunc(out, func(a, b result) int { return a.ID - b.ID })
	return out
}

// fanIn 把好幾個 channel 匯成一個；全部來源都關閉後，輸出也關閉。
func fanIn(chans ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	for _, c := range chans {
		wg.Go(func() {
			for v := range c {
				out <- v
			}
		})
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// gen 把幾個數字送進 channel，送完就關閉。
func gen(nums ...int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for _, n := range nums {
			ch <- n
		}
	}()
	return ch
}

// limited 跑 tasks 個工作，但同時最多 limit 個；回傳觀察到的最大同時數。
// sem 是有緩衝的 channel：放得進去才能開始，做完拿出來，就像號碼牌。
func limited(tasks, limit int) int {
	sem := make(chan struct{}, limit)
	var running, peak atomic.Int64
	var wg sync.WaitGroup
	for range tasks {
		wg.Go(func() {
			sem <- struct{}{}
			defer func() { <-sem }()

			now := running.Add(1)
			for { // 記錄最大值：別人可能同時在更新，用 CompareAndSwap 重試
				p := peak.Load()
				if now <= p || peak.CompareAndSwap(p, now) {
					break
				}
			}
			time.Sleep(time.Millisecond) // 模擬工作
			running.Add(-1)
		})
	}
	wg.Wait()
	return int(peak.Load())
}

func main() {
	fmt.Println("== worker pool ==")
	res := workerPool(context.Background(), []int{1, 2, 3, 4, 5}, 3, func(n int) int { return n * n })
	for _, r := range res {
		fmt.Printf("job %d → %d\n", r.ID, r.Value)
	}

	fmt.Println("== fan-in ==")
	var all []int
	for v := range fanIn(gen(1, 2), gen(3), gen(4, 5)) {
		all = append(all, v)
	}
	slices.Sort(all) // 匯流後的順序不固定
	fmt.Println(all)

	fmt.Println("== semaphore：同時最多 3 個 ==")
	fmt.Println(limited(20, 3) <= 3)
}

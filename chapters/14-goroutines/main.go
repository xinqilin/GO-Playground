// 第 14 章：goroutine。
//
// goroutine 是很便宜的分身：在函式呼叫前面加上 go，它就在背景自己跑。
// 執行：go run ./chapters/14-goroutines
package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("== 每個元素一個 goroutine ==")
	fmt.Println(squares([]int{1, 2, 3, 4}))

	fmt.Println("== 切段平行加總 ==")
	nums := make([]int, 1000)
	for i := range nums {
		nums[i] = i + 1
	}
	fmt.Println(sumParallel(nums, 4))

	fmt.Println("== goroutine 很便宜 ==")
	fmt.Println(spawn(10000))
}

// squares 為每個元素開一個 goroutine 算平方。
// 每個 goroutine 只寫自己那一格，彼此不會互相干擾。
func squares(nums []int) []int {
	out := make([]int, len(nums))
	var wg sync.WaitGroup
	for i, n := range nums {
		// Go 1.25 起的 wg.Go，等於 wg.Add(1) 再 go 一個結尾會呼叫 wg.Done() 的函式。
		// Go 1.22 起每一圈的 i、n 都是新的變數，可以放心在 goroutine 裡使用。
		wg.Go(func() {
			out[i] = n * n
		})
	}
	wg.Wait() // 等所有 goroutine 結束；沒有這行，main 可能先跑完，結果還沒算好
	return out
}

// sumParallel 把 nums 切成 parts 段，每段交給一個 goroutine 加總，最後再合起來。
func sumParallel(nums []int, parts int) int {
	parts = max(parts, 1)
	partial := make([]int, parts)
	size := (len(nums) + parts - 1) / parts
	var wg sync.WaitGroup
	for p := range parts {
		lo := min(p*size, len(nums))
		hi := min(lo+size, len(nums))
		wg.Add(1) // 傳統寫法：Add 一定要在 go 之前
		go func() {
			defer wg.Done()
			for _, n := range nums[lo:hi] {
				partial[p] += n
			}
		}()
	}
	wg.Wait()

	total := 0
	for _, s := range partial {
		total += s
	}
	return total
}

// spawn 開 n 個 goroutine，各自在自己的格子打勾，最後數有幾格被打勾。
func spawn(n int) int {
	done := make([]bool, n)
	var wg sync.WaitGroup
	for i := range n {
		wg.Go(func() { done[i] = true })
	}
	wg.Wait()

	count := 0
	for _, d := range done {
		if d {
			count++
		}
	}
	return count
}

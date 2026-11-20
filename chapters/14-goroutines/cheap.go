package main

import (
	"fmt"
	"sync"
)

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

func spawnDemo() {
	fmt.Println(spawn(10000))
}

// depth 遞迴 n 層，每一層都在 stack 上用掉一點空間。
func depth(n int) int {
	if n == 0 {
		return 0
	}
	return 1 + depth(n-1)
}

// growDemo：goroutine 的 stack 一開始很小，不夠用時 runtime 會自動把它變大。
func growDemo() {
	var got int
	var wg sync.WaitGroup
	wg.Go(func() {
		got = depth(100_000)
	})
	wg.Wait()
	fmt.Println(got)
}

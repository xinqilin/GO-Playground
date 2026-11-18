package main

import (
	"errors"
	"fmt"
	"sync"
)

// sumParallel 把 nums 切成 parts 段，每段交給一個 goroutine 加總到自己的格子，最後再合起來。
func sumParallel(nums []int, parts int) int {
	parts = max(parts, 1)
	partial := make([]int, parts)
	size := (len(nums) + parts - 1) / parts
	var wg sync.WaitGroup
	for p := range parts {
		lo := min(p*size, len(nums))
		hi := min(lo+size, len(nums))
		wg.Go(func() {
			for _, n := range nums[lo:hi] {
				partial[p] += n
			}
		})
	}
	wg.Wait()

	total := 0
	for _, s := range partial {
		total += s
	}
	return total
}

func sumDemo() {
	nums := make([]int, 1000)
	for i := range nums {
		nums[i] = i + 1
	}
	fmt.Println(sumParallel(nums, 4))
}

var ErrNoStore = errors.New("沒有這家分店")

// shelves 是各分店的排骨飯庫存。真實的程式會透過網路查詢，這裡用固定的資料模擬。
var shelves = map[string]int{"中山店": 12, "雙連店": 0, "民權店": 7}

func stock(store string) (int, error) {
	n, ok := shelves[store]
	if !ok {
		return 0, fmt.Errorf("%s：%w", store, ErrNoStore)
	}
	return n, nil
}

// checkStores 同時查詢每一家分店。結果和錯誤都寫進自己那一格，等全部查完再整理。
func checkStores(stores []string) ([]int, error) {
	counts := make([]int, len(stores))
	errs := make([]error, len(stores))
	var wg sync.WaitGroup
	for i, s := range stores {
		wg.Go(func() {
			counts[i], errs[i] = stock(s)
		})
	}
	wg.Wait()
	return counts, errors.Join(errs...)
}

func storesDemo() {
	counts, err := checkStores([]string{"中山店", "雙連店", "圓山店", "民權店"})
	fmt.Println(counts)
	fmt.Println(err)
	fmt.Println(errors.Is(err, ErrNoStore))
}

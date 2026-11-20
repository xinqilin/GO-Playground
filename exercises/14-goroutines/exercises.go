//go:build !solution

package crew

import "sync"

// Each 同時呼叫 f(0)、f(1)……f(n-1)，每一個都在自己的 goroutine 裡執行，全部結束才回傳。
// n 小於等於 0 時什麼都不做。
func Each(n int, f func(i int)) {
	// TODO: 實作
}

// MapParallel 對 s 的每個元素同時呼叫 f，每個元素一個 goroutine，
// 回傳的 slice 順序要和 s 一樣。
func MapParallel[T, U any](s []T, f func(T) U) []U {
	// TODO: 實作
	return nil
}

// CheckAll 同時檢查每一家分店，全部檢查完才回傳。
// 所有的錯誤用 errors.Join 合在一起，順序和 stores 一樣；全部沒問題時回傳 nil。
func CheckAll(stores []string, check func(store string) error) error {
	// TODO: 實作
	return nil
}

// CountWords 計算每個詞在所有文件裡一共出現幾次，詞用 strings.Fields 切。
// 每份文件交給一個 goroutine 計算；不可以讓好幾個 goroutine 同時寫同一個 map。
func CountWords(docs []string) map[string]int {
	// TODO: 實作
	return nil
}

// Try 在一個新的 goroutine 裡執行 f，等它結束。
// f 正常結束就回傳 nil；f panic 時不可以讓程式當掉，而是回傳一個錯誤：
// panic 的值是 error 時，用 %w 包起來，讓 errors.Is 找得到；其他的值用 %v。
// 錯誤訊息的開頭是 "panic: "。
func Try(f func()) error {
	// TODO: 實作
	return nil
}

// Total 回傳 b 和它底下所有分店的營業額總和。
// 每一家子分店交給一個 goroutine 遞迴計算；b 是 nil 時回傳 0。
func Total(b *Branch) int {
	// TODO: 實作
	return 0
}

// RunAll 同時執行所有的工作，全部做完才回傳。
//
// 找錯：每個工作都有做到，但它們是一個接一個做的，沒有同時執行。
func RunAll(tasks []func()) {
	var wg sync.WaitGroup
	for _, t := range tasks {
		wg.Go(t)
		wg.Wait()
	}
}

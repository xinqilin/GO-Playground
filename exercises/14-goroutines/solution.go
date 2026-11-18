//go:build solution

package crew

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

func Each(n int, f func(i int)) {
	var wg sync.WaitGroup
	for i := range max(n, 0) {
		wg.Go(func() { f(i) })
	}
	wg.Wait()
}

// MapParallel 讓每個 goroutine 只寫自己那一格，所以不需要鎖。
func MapParallel[T, U any](s []T, f func(T) U) []U {
	out := make([]U, len(s))
	var wg sync.WaitGroup
	for i, v := range s {
		wg.Go(func() { out[i] = f(v) })
	}
	wg.Wait()
	return out
}

// CheckAll 把錯誤寫進各自的格子，等全部結束再 Join；errors.Join 會略過 nil，全部是 nil 時回傳 nil。
func CheckAll(stores []string, check func(store string) error) error {
	errs := make([]error, len(stores))
	var wg sync.WaitGroup
	for i, s := range stores {
		wg.Go(func() { errs[i] = check(s) })
	}
	wg.Wait()
	return errors.Join(errs...)
}

// CountWords 讓每個 goroutine 數到自己的 map 裡，Wait 之後再由呼叫的 goroutine 合併。
func CountWords(docs []string) map[string]int {
	parts := make([]map[string]int, len(docs))
	var wg sync.WaitGroup
	for i, d := range docs {
		wg.Go(func() {
			m := map[string]int{}
			for _, w := range strings.Fields(d) {
				m[w]++
			}
			parts[i] = m
		})
	}
	wg.Wait()

	total := map[string]int{}
	for _, m := range parts {
		for w, n := range m {
			total[w] += n
		}
	}
	return total
}

// Try 的 recover 寫在新的 goroutine 自己的 defer 裡；寫在 Try 本身接不到別的 goroutine 的 panic。
func Try(f func()) (err error) {
	var wg sync.WaitGroup
	wg.Go(func() {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			if e, ok := r.(error); ok {
				err = fmt.Errorf("panic: %w", e)
			} else {
				err = fmt.Errorf("panic: %v", r)
			}
		}()
		f()
	})
	wg.Wait()
	return err
}

// Total 為每一家子分店開一個 goroutine，各自把結果寫進自己的格子。
func Total(b *Branch) int {
	if b == nil {
		return 0
	}
	sums := make([]int, len(b.Children))
	var wg sync.WaitGroup
	for i, c := range b.Children {
		wg.Go(func() { sums[i] = Total(c) })
	}
	wg.Wait()

	total := b.Sales
	for _, s := range sums {
		total += s
	}
	return total
}

// RunAll 先把所有工作都開起來，最後才 Wait 一次。
func RunAll(tasks []func()) {
	var wg sync.WaitGroup
	for _, t := range tasks {
		wg.Go(t)
	}
	wg.Wait()
}

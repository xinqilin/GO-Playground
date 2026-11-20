//go:build !solution

package rush

import (
	"context"
	"sync"
	"time"
)

// Limit 執行所有 tasks，同時最多 n 個在跑，全部做完才回傳。n 至少是 1。
func Limit(n int, tasks ...func()) {
	// TODO: 實作
}

// MapOrdered 用剛好 workers 個 goroutine 處理 in 的每一個值，
// 回傳的結果和 in 的順序相同。workers 至少是 1，比 len(in) 多也沒關係。
func MapOrdered[In, Out any](workers int, in []In, f func(In) Out) []Out {
	// TODO: 實作
	return nil
}

// Batches 從 in 收值，湊滿 size 個就送出一批。
// 這一批的第一個值收到之後過了 every 還沒湊滿，也直接送出；不送出空的批次。
// in 被 close 時，送出還沒送的那一批（如果有），再 close 回傳的 channel；
// ctx 被取消時，不必送出剩下的，直接 close 回傳的 channel。
func Batches[T any](ctx context.Context, in <-chan T, size int, every time.Duration) <-chan []T {
	// TODO: 實作
	out := make(chan []T)
	close(out)
	return out
}

// Team 同時執行好幾個工作，同時最多 limit 個；有一個失敗，就取消其他工作。用 NewTeam 建立。
type Team struct {
	// TODO: 加上需要的欄位
}

// NewTeam 回傳一個 Team，以及要交給工作使用的 ctx：它從 ctx 衍生，
// 在第一個工作失敗、或 Wait 回傳時被取消。limit 小於 1 表示不限制同時的數量。
func NewTeam(ctx context.Context, limit int) (*Team, context.Context) {
	// TODO: 實作
	return &Team{}, ctx
}

// Go 在新的 goroutine 執行 f。已經有 limit 個工作在跑時，Go 會等到有空位才開始並回傳。
func (t *Team) Go(f func() error) {
	// TODO: 實作
}

// Wait 等所有工作結束，回傳第一個失敗的錯誤；全部成功回傳 nil。
func (t *Team) Wait() error {
	// TODO: 實作
	return nil
}

// Prices 用 workers 個 worker 查每一道菜的價格，回傳菜名對價格的 map。
//
// 找錯：只要有菜要查，它就永遠不會回傳。
func Prices(workers int, dishes []string, lookup func(string) int) map[string]int {
	type priced struct {
		dish  string
		price int
	}
	jobs := make(chan string)
	results := make(chan priced)

	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			for d := range jobs {
				results <- priced{d, lookup(d)}
			}
		})
	}
	go func() {
		defer close(jobs)
		for _, d := range dishes {
			jobs <- d
		}
	}()

	wg.Wait()
	close(results)
	out := make(map[string]int)
	for r := range results {
		out[r.dish] = r.price
	}
	return out
}

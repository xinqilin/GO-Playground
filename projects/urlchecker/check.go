package main

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// Config 是檢查的設定。
type Config struct {
	Workers   int           // 同時最多幾個請求，至少 1
	Timeout   time.Duration // 每一次請求最多等多久
	PerSecond int           // 每秒最多送出幾個請求；0 是不限制
	Retries   int           // 失敗後最多再試幾次
	Backoff   time.Duration // 第一次重試前等多久，之後每次加倍
}

// Result 是一個網址的檢查結果。
type Result struct {
	URL      string
	Status   int           // HTTP 狀態碼；沒有收到回應時是 0
	Err      error         // 沒有收到回應的原因
	Attempts int           // 實際送出幾次請求
	Took     time.Duration // 包含等待限速和重試的時間
}

// OK 回傳這個網址是否正常：收到回應，而且狀態碼小於 400。
func (r Result) OK() bool { return r.Err == nil && r.Status < 400 }

// Check 用 cfg.Workers 個 worker 檢查 urls，回傳的結果順序和 urls 相同。
func Check(ctx context.Context, urls []string, cfg Config, f Fetcher) []Result {
	var tick <-chan time.Time // nil 代表不限速
	if cfg.PerSecond > 0 {
		t := time.NewTicker(time.Second / time.Duration(cfg.PerSecond))
		defer t.Stop()
		tick = t.C
	}

	results := make([]Result, len(urls))
	jobs := make(chan int)
	var wg sync.WaitGroup
	for range max(cfg.Workers, 1) {
		wg.Go(func() {
			for i := range jobs {
				results[i] = checkOne(ctx, urls[i], cfg, f, tick)
			}
		})
	}
	for i := range urls {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	return results
}

// checkOne 檢查一個網址：等限速、送出請求，遇到值得重試的失敗就等一下再試。
func checkOne(ctx context.Context, rawURL string, cfg Config, f Fetcher, tick <-chan time.Time) Result {
	r := Result{URL: rawURL}
	if err := validateURL(rawURL); err != nil {
		r.Err = err
		return r
	}
	start := time.Now()
	backoff := cfg.Backoff
	for {
		if err := waitTurn(ctx, tick); err != nil {
			r.Status, r.Err = 0, err
			break
		}
		r.Attempts++
		r.Status, r.Err = fetchOnce(ctx, rawURL, cfg.Timeout, f)
		if r.Attempts > cfg.Retries || !retryable(ctx, r) {
			break
		}
		if err := sleep(ctx, backoff); err != nil {
			r.Status, r.Err = 0, err
			break
		}
		backoff *= 2
	}
	r.Took = time.Since(start)
	return r
}

// fetchOnce 送出一次請求，最多等 timeout。
func fetchOnce(ctx context.Context, rawURL string, timeout time.Duration, f Fetcher) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return f.Fetch(ctx, rawURL)
}

// retryable 回傳 r 是否值得再試一次：伺服器暫時有問題（5xx、429），
// 或這一次沒有收到回應（例如逾時），而且整個檢查還沒被取消。
func retryable(ctx context.Context, r Result) bool {
	if ctx.Err() != nil {
		return false
	}
	return r.Err != nil || r.Status == http.StatusTooManyRequests || r.Status >= 500
}

// waitTurn 等到節拍器的下一拍；tick 是 nil 時不用等。ctx 被取消時回傳 ctx.Err()。
func waitTurn(ctx context.Context, tick <-chan time.Time) error {
	if tick == nil {
		return ctx.Err()
	}
	select {
	case <-tick:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// sleep 等 d，或 ctx 被取消。
func sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

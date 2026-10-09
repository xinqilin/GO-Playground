// 專題 3：並行網址檢查器。
//
// 同時檢查一批網址，最多同時發出 workers 個請求，每個請求最多等 timeout。
// 用到：goroutine 與 WaitGroup（第 14 章）、semaphore（第 19 章）、context 逾時（第 18 章）、
// interface 讓測試可以換掉真正的網路請求（第 10 章）。
//
// 執行：go run ./projects/urlchecker https://go.dev https://github.com
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// Result 是一個網址的檢查結果。
type Result struct {
	URL    string
	Status int
	Err    error
	Took   time.Duration
}

// String 把結果排成一行，例如「OK   200 https://go.dev」。
func (r Result) String() string {
	switch {
	case r.Err != nil:
		return fmt.Sprintf("FAIL --- %s（%v）", r.URL, r.Err)
	case r.Status >= 400:
		return fmt.Sprintf("FAIL %d %s", r.Status, r.URL)
	default:
		return fmt.Sprintf("OK   %d %s", r.Status, r.URL)
	}
}

// Fetcher 送出請求、回傳 HTTP 狀態碼。真正上線用 httpFetcher，測試時換成假的。
type Fetcher interface {
	Fetch(ctx context.Context, url string) (int, error)
}

// httpFetcher 用 HEAD 請求檢查網址，不下載內容。
type httpFetcher struct{ client *http.Client }

func (f httpFetcher) Fetch(ctx context.Context, url string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return 0, err
	}
	resp, err := f.client.Do(req)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()
	return resp.StatusCode, nil
}

// Check 同時檢查 urls，回傳的結果順序和 urls 相同。
func Check(ctx context.Context, urls []string, workers int, timeout time.Duration, f Fetcher) []Result {
	results := make([]Result, len(urls))
	sem := make(chan struct{}, max(workers, 1))
	var wg sync.WaitGroup
	for i, u := range urls {
		wg.Go(func() {
			select {
			case sem <- struct{}{}: // 拿到號碼牌才開始
			case <-ctx.Done():
				results[i] = Result{URL: u, Err: ctx.Err()}
				return
			}
			defer func() { <-sem }()

			reqCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			start := time.Now()
			status, err := f.Fetch(reqCtx, u)
			results[i] = Result{URL: u, Status: status, Err: err, Took: time.Since(start)}
		})
	}
	wg.Wait()
	return results
}

func main() {
	workers := flag.Int("workers", 4, "同時最多幾個請求")
	timeout := flag.Duration("timeout", 5*time.Second, "每個請求最多等多久")
	flag.Parse()
	urls := flag.Args()
	if len(urls) == 0 {
		fmt.Fprintln(os.Stderr, "用法：urlchecker [-workers 4] [-timeout 5s] URL...")
		os.Exit(2)
	}

	f := httpFetcher{client: &http.Client{}} // 逾時由 context 控制
	failed := 0
	for _, r := range Check(context.Background(), urls, *workers, *timeout, f) {
		fmt.Printf("%s  %v\n", r, r.Took.Round(time.Millisecond))
		if r.Err != nil || r.Status >= 400 {
			failed++
		}
	}
	if failed > 0 {
		os.Exit(1) // 有任何失敗就用非 0 結束，方便接在 CI 或排程裡
	}
}

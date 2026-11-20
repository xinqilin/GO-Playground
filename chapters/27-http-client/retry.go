package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"time"
)

// retryable 回傳這次的結果值不值得再試一次。
func retryable(resp *http.Response, err error) bool {
	if err != nil {
		// ctx 被取消或到期，是呼叫的一方不要了，再試也沒有用；其他的網路錯誤可能是暫時的。
		return !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded)
	}
	switch resp.StatusCode {
	case http.StatusTooManyRequests, http.StatusBadGateway,
		http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	}
	return false // 2xx 成功了；其他 4xx 是請求本身的問題，再送一次結果也一樣
}

// getWithRetry 送出 GET，遇到暫時性的失敗就等一下再試，最多 attempts 次。
// 等待的時間從 base 開始每次加倍，再加上最多一半的隨機抖動。
// 只用在冪等的請求：GET 重送幾次結果都一樣，POST 重送可能會下兩次單。
func getWithRetry(ctx context.Context, client *http.Client, u string, attempts int, base time.Duration) (*http.Response, error) {
	var resp *http.Response
	var err error
	for i := range attempts {
		if i > 0 {
			wait := base << (i - 1)
			wait += rand.N(wait/2 + 1)
			select {
			case <-time.After(wait):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
		req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
		if reqErr != nil {
			return nil, reqErr
		}
		resp, err = client.Do(req)
		if !retryable(resp, err) || i == attempts-1 {
			break
		}
		if resp != nil {
			// 要再試一次，這個回應用不到了：讀完並關掉，連線才能重用。
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}
	return resp, err
}

// flaky 是一個前兩次都說「太忙」的店家。
func flaky() http.Handler {
	calls := 0
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		fmt.Printf("店家收到第 %d 次請求\n", calls)
		if calls <= 2 {
			http.Error(w, "太忙了", http.StatusServiceUnavailable)
			return
		}
		fmt.Fprint(w, "今天的便當：排骨飯")
	})
}

func retryDemo() {
	client := &http.Client{Timeout: 5 * time.Second, Transport: handlerTransport{flaky()}}
	resp, err := getWithRetry(context.Background(), client, "http://shop.test/today", 4, 10*time.Millisecond)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode, string(body))
}

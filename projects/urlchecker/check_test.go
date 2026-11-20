package main

import (
	"context"
	"errors"
	"slices"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

// fakeFetcher 依網址決定怎麼回應，用計時器模擬網路延遲。
// 在 synctest.Test 裡計時器用的是假的時鐘：測試不用真的等，時間也完全精確。
type fakeFetcher struct {
	mu      sync.Mutex
	calls   map[string]int // 每個網址被請求了幾次
	running int            // 現在有幾個請求在進行
	peak    int            // running 曾經到過的最大值
	starts  []time.Time    // 每個請求開始的時間
}

func newFake() *fakeFetcher { return &fakeFetcher{calls: map[string]int{}} }

func (f *fakeFetcher) Fetch(ctx context.Context, rawURL string) (int, error) {
	f.mu.Lock()
	f.calls[rawURL]++
	n := f.calls[rawURL]
	f.running++
	f.peak = max(f.peak, f.running)
	f.starts = append(f.starts, time.Now())
	f.mu.Unlock()
	defer func() {
		f.mu.Lock()
		f.running--
		f.mu.Unlock()
	}()

	switch rawURL {
	case "https://ok.test":
		return respond(ctx, 120*time.Millisecond, 200)
	case "https://missing.test":
		return respond(ctx, 80*time.Millisecond, 404)
	case "https://flaky.test": // 前兩次 503，第三次才成功
		if n < 3 {
			return respond(ctx, 50*time.Millisecond, 503)
		}
		return respond(ctx, 50*time.Millisecond, 200)
	}
	<-ctx.Done() // 其他網址永遠不回應，直到逾時或被取消
	return 0, ctx.Err()
}

// respond 等 d 之後回傳 status；等的途中 ctx 結束就回傳 ctx.Err()。
func respond(ctx context.Context, d time.Duration, status int) (int, error) {
	select {
	case <-time.After(d):
		return status, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

// TestCheckWorkers：20 個網址、3 個 worker，同時最多 3 個請求。
// 每個請求 120ms，要 7 輪（3+3+3+3+3+3+2），總共剛好 840ms。
func TestCheckWorkers(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		urls := slices.Repeat([]string{"https://ok.test"}, 20)
		f := newFake()
		start := time.Now()
		results := Check(t.Context(), urls, Config{Workers: 3, Timeout: time.Second}, f)
		if took := time.Since(start); took != 840*time.Millisecond {
			t.Errorf("花了 %v，想要 840ms", took)
		}
		if f.peak != 3 {
			t.Errorf("同時最多 %d 個請求，想要 3 個", f.peak)
		}
		for i, r := range results {
			if !r.OK() || r.URL != urls[i] {
				t.Errorf("results[%d] = %+v", i, r)
			}
		}
	})
}

// TestCheckRate：每秒 10 個請求，就是每 100ms 一個。
// 5 個 worker 都閒著，請求還是照節拍在 100、200、300、400、500ms 送出。
func TestCheckRate(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		urls := slices.Repeat([]string{"https://ok.test"}, 5)
		f := newFake()
		start := time.Now()
		Check(t.Context(), urls, Config{Workers: 5, Timeout: time.Second, PerSecond: 10}, f)

		var got []time.Duration
		for _, s := range f.starts {
			got = append(got, s.Sub(start))
		}
		want := []time.Duration{100 * time.Millisecond, 200 * time.Millisecond,
			300 * time.Millisecond, 400 * time.Millisecond, 500 * time.Millisecond}
		if !slices.Equal(got, want) {
			t.Errorf("請求開始的時間 = %v，想要 %v", got, want)
		}
	})
}

func TestCheckRetry(t *testing.T) {
	tests := []struct {
		url          string
		retries      int
		wantStatus   int
		wantErr      error
		wantAttempts int
		wantTook     time.Duration
	}{
		{"https://flaky.test", 0, 503, nil, 1, 50 * time.Millisecond},
		{"https://flaky.test", 1, 503, nil, 2, 300 * time.Millisecond},                    // 50 + 等 200 + 50
		{"https://flaky.test", 2, 200, nil, 3, 750 * time.Millisecond},                    // 再等 400 + 50
		{"https://flaky.test", 5, 200, nil, 3, 750 * time.Millisecond},                    // 成功了就不再試
		{"https://missing.test", 2, 404, nil, 1, 80 * time.Millisecond},                   // 404 不會因為重試就變好
		{"https://slow.test", 1, 0, context.DeadlineExceeded, 2, 1200 * time.Millisecond}, // 逾時 500 + 等 200 + 逾時 500
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			cfg := Config{Workers: 1, Timeout: 500 * time.Millisecond, Retries: tt.retries, Backoff: 200 * time.Millisecond}
			r := Check(t.Context(), []string{tt.url}, cfg, newFake())[0]
			if r.Status != tt.wantStatus || !errors.Is(r.Err, tt.wantErr) || r.Attempts != tt.wantAttempts || r.Took != tt.wantTook {
				t.Errorf("%s 重試 %d 次：狀態 %d、錯誤 %v、送了 %d 次、花了 %v；想要 %d、%v、%d 次、%v",
					tt.url, tt.retries, r.Status, r.Err, r.Attempts, r.Took, tt.wantStatus, tt.wantErr, tt.wantAttempts, tt.wantTook)
			}
		})
	}
}

// TestCheckCanceled：300ms 時取消整個檢查。正在等的請求馬上結束、不再重試；
// 還沒輪到的網址一個請求都不送，直接標成已取消。
func TestCheckCanceled(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		time.AfterFunc(300*time.Millisecond, cancel)

		urls := []string{"https://slow.test", "https://ok.test", "https://ok.test"}
		cfg := Config{Workers: 1, Timeout: time.Second, Retries: 2}
		start := time.Now()
		results := Check(ctx, urls, cfg, newFake())
		if took := time.Since(start); took != 300*time.Millisecond {
			t.Errorf("花了 %v，想要 300ms", took)
		}
		for i, r := range results {
			wantAttempts := 0
			if i == 0 {
				wantAttempts = 1
			}
			if !errors.Is(r.Err, context.Canceled) || r.Attempts != wantAttempts {
				t.Errorf("results[%d]：錯誤 %v、送了 %d 次；想要 Canceled、%d 次", i, r.Err, r.Attempts, wantAttempts)
			}
		}
	})
}

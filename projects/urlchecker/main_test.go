package main

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

// fakeFetcher 依網址決定回應，不碰網路；slow.test 會一直等到逾時。
type fakeFetcher struct {
	running, peak atomic.Int64
}

func (f *fakeFetcher) Fetch(ctx context.Context, url string) (int, error) {
	now := f.running.Add(1)
	defer f.running.Add(-1)
	for {
		p := f.peak.Load()
		if now <= p || f.peak.CompareAndSwap(p, now) {
			break
		}
	}

	switch url {
	case "https://ok.test":
		return 200, nil
	case "https://missing.test":
		return 404, nil
	default:
		<-ctx.Done()
		return 0, ctx.Err()
	}
}

func Example() {
	urls := []string{"https://ok.test", "https://missing.test", "https://slow.test"}
	for _, r := range Check(context.Background(), urls, 2, 50*time.Millisecond, &fakeFetcher{}) {
		fmt.Println(r)
	}
	// Output:
	// OK   200 https://ok.test
	// FAIL 404 https://missing.test
	// FAIL --- https://slow.test（context deadline exceeded）
}

func TestCheckKeepsOrderAndLimit(t *testing.T) {
	urls := make([]string, 20)
	for i := range urls {
		urls[i] = "https://ok.test"
	}
	urls[7] = "https://slow.test"

	f := &fakeFetcher{}
	res := Check(context.Background(), urls, 3, 20*time.Millisecond, f)
	for i, r := range res {
		wantErr := i == 7
		if (r.Err != nil) != wantErr {
			t.Errorf("res[%d] = %v", i, r)
		}
	}
	if !errors.Is(res[7].Err, context.DeadlineExceeded) {
		t.Errorf("slow.test 應該逾時，got %v", res[7].Err)
	}
	if p := f.peak.Load(); p > 3 {
		t.Errorf("同時最多 3 個請求，實際 %d 個", p)
	}
}

func TestCheckCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for _, r := range Check(ctx, []string{"https://slow.test"}, 1, time.Hour, &fakeFetcher{}) {
		if !errors.Is(r.Err, context.Canceled) {
			t.Errorf("已取消的 context 應該回傳 Canceled，got %v", r.Err)
		}
	}
}

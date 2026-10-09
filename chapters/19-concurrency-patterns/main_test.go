package main

import (
	"context"
	"testing"
)

func Example() {
	main()
	// Output:
	// == worker pool ==
	// job 0 → 1
	// job 1 → 4
	// job 2 → 9
	// job 3 → 16
	// job 4 → 25
	// == fan-in ==
	// [1 2 3 4 5]
	// == semaphore：同時最多 3 個 ==
	// true
}

func TestWorkerPool(t *testing.T) {
	jobs := make([]int, 50)
	for i := range jobs {
		jobs[i] = i
	}
	res := workerPool(context.Background(), jobs, 4, func(n int) int { return n + 1 })
	if len(res) != len(jobs) {
		t.Fatalf("got %d results, want %d", len(res), len(jobs))
	}
	for i, r := range res {
		if r.ID != i || r.Value != i+1 {
			t.Errorf("res[%d] = %+v, want {ID:%d Value:%d}", i, r, i, i+1)
		}
	}
}

// 已經取消的 context：worker pool 要能結束，而且不會多做完所有工作。
func TestWorkerPoolCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	res := workerPool(ctx, make([]int, 1000), 4, func(n int) int { return n })
	if len(res) == 1000 {
		t.Errorf("取消後不該把 1000 個工作全部做完")
	}
}

func TestLimited(t *testing.T) {
	for _, limit := range []int{1, 2, 5} {
		if peak := limited(30, limit); peak > limit {
			t.Errorf("limit %d: peak %d", limit, peak)
		}
	}
}

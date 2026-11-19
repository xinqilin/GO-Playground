//go:build solution

package rush

import (
	"context"
	"sync"
	"time"
)

// Limit 用容量 n 的 channel 當號碼牌。在開 goroutine 之前就拿號碼牌，
// 同時存在的 goroutine 也不會超過 n 個。
func Limit(n int, tasks ...func()) {
	sem := make(chan struct{}, n)
	var wg sync.WaitGroup
	for _, task := range tasks {
		sem <- struct{}{}
		wg.Go(func() {
			defer func() { <-sem }()
			task()
		})
	}
	wg.Wait()
}

// MapOrdered 派工時送的是索引，worker 把結果寫進自己那一格，順序自然和 in 相同。
func MapOrdered[In, Out any](workers int, in []In, f func(In) Out) []Out {
	jobs := make(chan int)
	out := make([]Out, len(in))
	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			for i := range jobs {
				out[i] = f(in[i])
			}
		})
	}
	for i := range in {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	return out
}

// Batches 用一個計時 channel flush：這一批是空的時候它是 nil，select 不會選到；
// 收到這一批的第一個值時才開始計時。
func Batches[T any](ctx context.Context, in <-chan T, size int, every time.Duration) <-chan []T {
	out := make(chan []T)
	go func() {
		defer close(out)
		var batch []T
		var flush <-chan time.Time

		send := func() bool {
			select {
			case out <- batch:
				batch, flush = nil, nil
				return true
			case <-ctx.Done():
				return false
			}
		}

		for {
			select {
			case v, ok := <-in:
				if !ok {
					if len(batch) > 0 {
						send()
					}
					return
				}
				if len(batch) == 0 {
					flush = time.After(every)
				}
				batch = append(batch, v)
				if len(batch) == size && !send() {
					return
				}
			case <-flush:
				if !send() {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// Team 是第 5 節的 Group 加上號碼牌：Go 先拿號碼牌才開 goroutine，做完還回去。
type Team struct {
	wg     sync.WaitGroup
	sem    chan struct{}
	cancel context.CancelCauseFunc
	once   sync.Once
	err    error
}

func NewTeam(ctx context.Context, limit int) (*Team, context.Context) {
	ctx, cancel := context.WithCancelCause(ctx)
	t := &Team{cancel: cancel}
	if limit > 0 {
		t.sem = make(chan struct{}, limit)
	}
	return t, ctx
}

func (t *Team) Go(f func() error) {
	if t.sem != nil {
		t.sem <- struct{}{}
	}
	t.wg.Go(func() {
		if t.sem != nil {
			defer func() { <-t.sem }()
		}
		if err := f(); err != nil {
			t.once.Do(func() {
				t.err = err
				t.cancel(err)
			})
		}
	})
}

func (t *Team) Wait() error {
	t.wg.Wait()
	t.cancel(t.err)
	return t.err
}

// Prices 原本在讀結果之前就 wg.Wait()：worker 送結果時沒有人收，全部卡在送，
// Wait 也就永遠等不到。改成另開一個 goroutine 等 worker 結束再 close(results)，
// 這邊同時一路收結果。
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
	go func() {
		wg.Wait()
		close(results)
	}()

	out := make(map[string]int)
	for r := range results {
		out[r.dish] = r.price
	}
	return out
}

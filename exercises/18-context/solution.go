//go:build solution

package courier

import (
	"context"
	"errors"
	"slices"
	"sync"
	"time"
)

// Sleep 用 time.NewTimer 而不是 time.After，是為了提早回來時能 Stop 掉計時器；
// Go 1.23 起不 Stop 也會被回收，兩種寫法都對。
func Sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func ShouldStart(ctx context.Context, need time.Duration) bool {
	if ctx.Err() != nil {
		return false
	}
	deadline, ok := ctx.Deadline()
	return !ok || time.Until(deadline) >= need
}

type routeKey struct{}

// WithStop 一定要做一份新的 slice：直接 append 到上一層的 slice，
// 底層陣列還有空間時，從同一個 ctx 分出去的兩條路線會寫到同一格。
func WithStop(ctx context.Context, stop string) context.Context {
	prev, _ := ctx.Value(routeKey{}).([]string)
	next := append(slices.Clip(prev), stop)
	return context.WithValue(ctx, routeKey{}, next)
}

// Route 回傳一份複本，呼叫的人怎麼改都不會影響 ctx 裡存的值。
func Route(ctx context.Context) []string {
	route, _ := ctx.Value(routeKey{}).([]string)
	return slices.Clone(route)
}

// Pump 的收和送都放在 select 裡，另一個 case 是 ctx.Done()。
func Pump[T any](ctx context.Context, in <-chan T, out chan<- T) error {
	for {
		var v T
		select {
		case got, ok := <-in:
			if !ok {
				return nil
			}
			v = got
		case <-ctx.Done():
			return ctx.Err()
		}
		select {
		case out <- v:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Fastest 從 ctx 衍生出一個可以取消的 ctx 交給每一間廚房。
// results 的容量和廚房一樣多，誰送結果都不會卡住；有人成功就 cancel，再 Wait 大家結束。
func Fastest(ctx context.Context, kitchens ...func(context.Context) (string, error)) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	type result struct {
		dish string
		err  error
	}
	results := make(chan result, len(kitchens))
	var wg sync.WaitGroup
	for _, k := range kitchens {
		wg.Go(func() {
			dish, err := k(ctx)
			results <- result{dish, err}
		})
	}

	var errs []error
	for range kitchens {
		r := <-results
		if r.err == nil {
			cancel()
			wg.Wait()
			return r.dish, nil
		}
		errs = append(errs, r.err)
	}
	return "", errors.Join(errs...)
}

// Quote 原本用 context.Background() 當上層，和呼叫的人的 ctx 斷了關係，
// 取消傳不下來。從傳進來的 ctx 衍生，就同時有「最多 1 秒」和「呼叫的人取消就停」。
func Quote(ctx context.Context, ask func(context.Context) (time.Duration, error)) (time.Duration, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return ask(ctx)
}

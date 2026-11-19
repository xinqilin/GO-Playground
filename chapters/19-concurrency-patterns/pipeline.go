package main

import (
	"context"
	"fmt"
	"strconv"
)

// emit 在另一個 goroutine 裡依序送出 items；送完或 ctx 被取消就 close。
func emit[T any](ctx context.Context, items ...T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for _, it := range items {
			select {
			case out <- it:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// stage 是管線的一站：把 in 的每個值經過 f 送出去。
// 收的時候不用 select ctx：上游也聽 ctx，被取消時會 close，range 自然結束。
// 送的時候一定要 select ctx：下游不收了，這一站才不會永遠卡在送。
func stage[In, Out any](ctx context.Context, in <-chan In, f func(In) Out) <-chan Out {
	out := make(chan Out)
	go func() {
		defer close(out)
		for v := range in {
			select {
			case out <- f(v):
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// amount 是一張發票的金額；err 不是 nil 表示這張發票有問題，沿著管線往下傳。
type amount struct {
	v   int
	err error
}

func parse(s string) amount {
	v, err := strconv.Atoi(s)
	if err != nil {
		return amount{err: fmt.Errorf("發票金額「%s」不是數字", s)}
	}
	return amount{v: v}
}

// addTax 加上 5% 營業稅，無條件捨去到整數。
func addTax(a amount) amount {
	if a.err != nil {
		return a
	}
	return amount{v: a.v * 105 / 100}
}

// pipelineDemo：遇到第一張有問題的發票就停。break 之後 defer 的 cancel 讓每一站都結束。
func pipelineDemo() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	raw := emit(ctx, "120", "85", "九十", "300")
	taxed := stage(ctx, stage(ctx, raw, parse), addTax)

	total := 0
	for a := range taxed {
		if a.err != nil {
			fmt.Println("停下來：", a.err)
			break
		}
		total += a.v
	}
	fmt.Println("含稅合計", total)
}

package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Group 同時執行好幾個工作。第一個失敗的工作會取消其他工作，Wait 回傳那個錯誤。
type Group struct {
	wg     sync.WaitGroup
	cancel context.CancelCauseFunc
	once   sync.Once
	err    error
}

// NewGroup 回傳一個 Group，以及交給工作用的 ctx：有工作失敗、或 Wait 回傳時，它會被取消。
func NewGroup(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancelCause(ctx)
	return &Group{cancel: cancel}, ctx
}

// Go 在新的 goroutine 執行 f。f 回傳錯誤時，只有第一個會被記下來並取消 ctx。
func (g *Group) Go(f func() error) {
	g.wg.Go(func() {
		if err := f(); err != nil {
			g.once.Do(func() {
				g.err = err
				g.cancel(err)
			})
		}
	})
}

// Wait 等所有工作結束，回傳第一個錯誤。WaitGroup 和 Once 保證這裡讀得到 g.err 最後的值。
func (g *Group) Wait() error {
	g.wg.Wait()
	g.cancel(g.err)
	return g.err
}

// groupDemo：三家分店同時備料。西區一開始就失敗，東區和北區還在做，收到取消就停下來。
func groupDemo() {
	g, ctx := NewGroup(context.Background())
	branches := []string{"東區", "西區", "北區"}
	status := make([]string, len(branches))
	for i, b := range branches {
		g.Go(func() error {
			if b == "西區" {
				status[i] = "缺貨"
				return errors.New("西區 排骨缺貨")
			}
			<-ctx.Done() // 模擬一件很久的工作，只會因為取消而停下來
			status[i] = "停下來：" + ctx.Err().Error()
			return ctx.Err()
		})
	}

	fmt.Println("錯誤：", g.Wait())
	for i, b := range branches {
		fmt.Println(b, status[i])
	}
}

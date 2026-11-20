package main

import (
	"context"
	"fmt"
	"time"
)

// notifyAll 依序通知每一位客人取餐，但兩次通知之間至少隔 every（簡訊商限制發送的速率）。
// 每一次通知前都等下一個 tick；ctx 被取消就停下來。
func notifyAll(ctx context.Context, phones []string, every time.Duration, send func(string)) error {
	t := time.NewTicker(every)
	defer t.Stop()
	for _, p := range phones {
		select {
		case <-t.C:
			send(p)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func tickerDemo() {
	phones := []string{"0912-000-111", "0922-000-222", "0933-000-333"}
	err := notifyAll(context.Background(), phones, time.Millisecond, func(p string) {
		fmt.Println("已通知", p)
	})
	fmt.Println("錯誤：", err)
}

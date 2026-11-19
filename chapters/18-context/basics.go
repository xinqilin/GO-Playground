package main

import (
	"context"
	"fmt"
	"time"
)

// cook 模擬一道要煮 d 的菜。ctx 被取消（客人取消、超過期限）就不煮了，回傳原因。
func cook(ctx context.Context, d time.Duration) error {
	select {
	case <-time.After(d):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func basicsDemo() {
	if err := cook(context.Background(), time.Millisecond); err == nil {
		fmt.Println("滷蛋 好了")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if err := cook(ctx, time.Hour); err != nil {
		fmt.Println("佛跳牆 放棄：", err)
	}
}

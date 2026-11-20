package main

import (
	"context"
	"fmt"
	"time"
)

// orderIDKey 是存訂單編號用的 key。用不匯出的型別，別的套件不可能做出相同的 key。
type orderIDKey struct{}

// withOrderID 回傳一個帶著訂單編號的 ctx。
func withOrderID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, orderIDKey{}, id)
}

// orderID 取出訂單編號；ctx 裡沒有的話，ok 是 false。
func orderID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(orderIDKey{}).(string)
	return id, ok
}

// logLine 印出一行紀錄，開頭帶上這個請求的訂單編號。
func logLine(ctx context.Context, msg string) {
	id, ok := orderID(ctx)
	if !ok {
		id = "-"
	}
	fmt.Printf("[%s] %s\n", id, msg)
}

// valueDemo：值會跟著衍生出來的 ctx 一路往下傳。
func valueDemo() {
	ctx := withOrderID(context.Background(), "A-07")
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	logLine(ctx, "開始備餐")
	logLine(context.Background(), "沒有訂單的紀錄")
}

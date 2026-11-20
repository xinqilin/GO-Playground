package main

import (
	"context"
	"fmt"
)

// detachDemo：請求結束（被取消）之後還要做的事，例如寫對帳紀錄，
// 用 WithoutCancel 拿到一個「保留值、但不會被取消」的 ctx。
func detachDemo() {
	req, cancel := context.WithCancel(withOrderID(context.Background(), "A-07"))
	audit := context.WithoutCancel(req)
	cancel()

	id, _ := orderID(audit)
	fmt.Println("請求：", req.Err())
	fmt.Println("對帳紀錄：", audit.Err(), id)
}

// afterDemo：AfterFunc 讓 ctx 被取消時，在另一個 goroutine 執行一個函式。
// stop 可以在那之前取消這個安排。
func afterDemo() {
	ctx, cancel := context.WithCancel(context.Background())
	notified := make(chan struct{})
	context.AfterFunc(ctx, func() {
		fmt.Println("通知外送員：訂單取消")
		close(notified)
	})
	cancel()
	<-notified

	ctx2, cancel2 := context.WithCancel(context.Background())
	stop := context.AfterFunc(ctx2, func() {
		fmt.Println("這一行不會印出來")
	})
	fmt.Println("撤回通知：", stop())
	cancel2()
}

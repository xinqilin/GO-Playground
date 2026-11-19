package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	errSoldOut = errors.New("排骨賣完了")
	errSlow    = errors.New("廚房出餐太慢")
)

// causeDemo：Err 只會是 Canceled 或 DeadlineExceeded；真正的原因用 context.Cause 取出。
func causeDemo() {
	ctx, cancel := context.WithCancelCause(context.Background())
	child, cancelChild := context.WithCancel(ctx)
	defer cancelChild()
	cancel(errSoldOut)
	fmt.Println(ctx.Err(), "/", context.Cause(ctx))
	fmt.Println("衍生的 ctx 也拿得到原因：", context.Cause(child))

	slow, cancelSlow := context.WithTimeoutCause(context.Background(), 10*time.Millisecond, errSlow)
	defer cancelSlow()
	<-slow.Done()
	fmt.Println(slow.Err(), "/", context.Cause(slow))

	plain, cancelPlain := context.WithCancel(context.Background())
	cancelPlain()
	fmt.Println("沒有指定原因：", context.Cause(plain))
}

package main

import (
	"context"
	"fmt"
	"time"
)

// holdSeat 幫客人保留座位 d 這麼久；時間到了還沒付款，就呼叫 release 把座位放出來。
// 回傳的 Timer 在客人付款時呼叫 Stop，release 就不會被呼叫。
func holdSeat(seat string, d time.Duration, release func(seat string)) *time.Timer {
	return time.AfterFunc(d, func() { release(seat) })
}

// afterFuncDemo 一個座位等到逾時，另一個在逾時前付款。
func afterFuncDemo() {
	released := make(chan string, 1)
	release := func(seat string) { released <- seat }

	holdSeat("A1", 10*time.Millisecond, release)
	fmt.Println("逾時放出：", <-released)

	t := holdSeat("B2", time.Hour, release)
	fmt.Println("B2 付款了，計時停下來：", t.Stop())
	fmt.Println("再 Stop 一次：", t.Stop())
}

// callNumbers 每隔 every 叫一個號碼，從 first 叫到 last；ctx 取消就提早停下來。
func callNumbers(ctx context.Context, first, last int, every time.Duration, call func(n int)) {
	t := time.NewTicker(every)
	defer t.Stop()
	for n := first; n <= last; n++ {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			call(n)
		}
	}
}

// tickerDemo 叫號機每 10 毫秒叫一號。
func tickerDemo() {
	callNumbers(context.Background(), 101, 103, 10*time.Millisecond, func(n int) {
		fmt.Printf("請 %d 號到櫃台\n", n)
	})
}

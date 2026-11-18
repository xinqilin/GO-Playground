package main

import "fmt"

// handoffDemo：無緩衝的 channel 是出餐口，送的一方要等到有人來拿，才算交出去。
func handoffDemo() {
	window := make(chan string)
	go func() {
		window <- "排骨飯"
		window <- "雞腿飯"
	}()
	fmt.Println("拿到", <-window)
	fmt.Println("拿到", <-window)
}

// doneDemo：用一個 channel 等另一個 goroutine 做完。close 之前寫的東西，收到通知之後一定看得到。
func doneDemo() {
	done := make(chan struct{})
	var receipt string
	go func() {
		receipt = "排骨飯 x2，NT$170"
		close(done)
	}()
	<-done
	fmt.Println(receipt)
}

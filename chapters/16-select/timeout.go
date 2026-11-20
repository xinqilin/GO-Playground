package main

import (
	"fmt"
	"time"
)

// ready 回傳一個已經放好結果的 channel。
func ready(v string) <-chan string {
	ch := make(chan string, 1)
	ch <- v
	return ch
}

// never 回傳一個永遠不會有結果的 channel，模擬一直沒有回應的廚房。
func never() <-chan string {
	return make(chan string)
}

// waitFor 等 ch 的結果，最多等 limit。
func waitFor(ch <-chan string, limit time.Duration) string {
	select {
	case v := <-ch:
		return "拿到 " + v
	case <-time.After(limit):
		return "逾時"
	}
}

func timeoutDemo() {
	fmt.Println(waitFor(ready("排骨飯"), time.Second))
	fmt.Println(waitFor(never(), 10*time.Millisecond))
}

// collectUntil 一直收 ch，直到 ch 被 close，或從開始算起過了 limit。
// 計時器在迴圈外面只建立一次，所以 limit 是整段的期限，不會每收到一個值就重新計時。
func collectUntil(ch <-chan string, limit time.Duration) []string {
	var got []string
	deadline := time.After(limit)
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return got
			}
			got = append(got, v)
		case <-deadline:
			return got
		}
	}
}

func deadlineDemo() {
	ch := make(chan string, 3)
	ch <- "排骨飯"
	ch <- "雞腿飯"
	fmt.Println(collectUntil(ch, 20*time.Millisecond))
}

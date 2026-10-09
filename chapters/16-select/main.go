// 第 16 章：select。
//
// select 是同時盯著好幾個 channel 的櫃檯：哪一個先準備好，就先處理哪一個。
// 執行：go run ./chapters/16-select
package main

import (
	"fmt"
	"slices"
	"time"
)

func main() {
	fmt.Println("== 逾時 ==")
	fmt.Println(withTimeout(fetch(0), time.Second))
	fmt.Println(withTimeout(fetch(-1), 10*time.Millisecond))

	fmt.Println("== default：不等 ==")
	ch := make(chan int, 1)
	fmt.Println(tryReceive(ch))
	ch <- 5
	fmt.Println(tryReceive(ch))

	fmt.Println("== 合併兩個 channel ==")
	fmt.Println(merge(from(1, 3, 5), from(2, 4)))
}

// fetch 模擬一個要等 d 才有結果的工作；d < 0 代表永遠不會有結果。
func fetch(d time.Duration) <-chan string {
	ch := make(chan string, 1) // 有緩衝：就算沒人來收，送的 goroutine 也不會卡住而洩漏
	if d < 0 {
		return ch
	}
	go func() {
		time.Sleep(d)
		ch <- "資料"
	}()
	return ch
}

// withTimeout 等 ch 的結果，最多等 limit。
func withTimeout(ch <-chan string, limit time.Duration) string {
	select {
	case v := <-ch:
		return "拿到" + v
	case <-time.After(limit):
		return "逾時"
	}
}

// tryReceive 用 default 做非阻塞的接收：沒有東西就馬上離開。
func tryReceive(ch <-chan int) string {
	select {
	case v := <-ch:
		return fmt.Sprint("收到 ", v)
	default:
		return "現在沒有"
	}
}

// from 把幾個數字送進 channel，送完就關閉。
func from(nums ...int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for _, n := range nums {
			ch <- n
		}
	}()
	return ch
}

// merge 用 for-select 同時收兩個 channel。
// 某一邊關閉後把它設成 nil：從 nil channel 收會永遠卡住，select 就不會再選它。
func merge(a, b <-chan int) []int {
	var out []int
	for a != nil || b != nil {
		select {
		case v, ok := <-a:
			if !ok {
				a = nil
				continue
			}
			out = append(out, v)
		case v, ok := <-b:
			if !ok {
				b = nil
				continue
			}
			out = append(out, v)
		}
	}
	slices.Sort(out) // 兩邊交錯的順序每次都可能不同，排序後才好比對
	return out
}

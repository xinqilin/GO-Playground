package main

import "fmt"

// orders 在另一個 goroutine 裡依序送出每一份訂單，送完就 close。
// 回傳型別 <-chan string 表示拿到的人只能收、不能送，也不能 close。
func orders(items ...string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for _, it := range items {
			out <- it
		}
	}()
	return out
}

// rangeDemo：range 一個 channel 會一直收，直到 channel 被 close 而且裡面的值都收完。
func rangeDemo() {
	for o := range orders("排骨飯", "雞腿飯", "魚排飯") {
		fmt.Println("出餐", o)
	}
}

// closedDemo：close 之後，還沒收的值照樣收得到；收完之後，每次都立刻拿到零值和 false。
func closedDemo() {
	c := make(chan int, 2)
	c <- 17
	c <- 18
	close(c)
	for range 4 {
		v, ok := <-c
		fmt.Println(v, ok)
	}
}

// try 執行 f，把 panic 的內容印出來，讓示範可以繼續。
func try(f func()) {
	defer func() {
		fmt.Println("panic:", recover())
	}()
	f()
}

// closePanicDemo：對已經 close 的 channel 送值、再 close 一次，或 close 一個 nil channel，都會 panic。
func closePanicDemo() {
	c := make(chan int, 1)
	close(c)
	try(func() { c <- 1 })
	try(func() { close(c) })
	var never chan int
	try(func() { close(never) })
}

package main

import "fmt"

// bufferDemo：有緩衝的 channel 是一排號碼牌架，還有空位就能直接放，先放的先拿。
func bufferDemo() {
	tickets := make(chan int, 3)
	tickets <- 17
	tickets <- 18
	tickets <- 19
	fmt.Println(len(tickets), cap(tickets))

	fmt.Println("叫號", <-tickets)
	tickets <- 20
	fmt.Println(len(tickets), cap(tickets))

	fmt.Println("叫號", <-tickets, <-tickets, <-tickets)
	fmt.Println(len(tickets), cap(tickets))
}

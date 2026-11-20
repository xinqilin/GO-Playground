package main

import "fmt"

// next 先看急件，沒有急件才同時等兩邊。
// 只寫一個 select 的話，兩邊都有單時是隨機挑，急件可能一直被普通單插隊。
func next(urgent, normal <-chan string) string {
	select {
	case o := <-urgent:
		return o
	default:
	}
	select {
	case o := <-urgent:
		return o
	case o := <-normal:
		return o
	}
}

func priorityDemo() {
	urgent := make(chan string, 3)
	normal := make(chan string, 3)
	normal <- "普通 1"
	normal <- "普通 2"
	urgent <- "急件 1"
	urgent <- "急件 2"

	for range 4 {
		fmt.Println(next(urgent, normal))
	}
}

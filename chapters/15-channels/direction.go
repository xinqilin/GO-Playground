package main

import "fmt"

// announce 只負責送：參數型別是 chan<- int，在這個函式裡只能送、不能收。
// 送完由它 close，因為只有送的一方知道什麼時候沒有東西了。
func announce(nums []int, out chan<- int) {
	for _, n := range nums {
		out <- n
	}
	close(out)
}

// display 只負責收：參數型別是 <-chan int。
func display(in <-chan int) {
	for n := range in {
		fmt.Println("請", n, "號")
	}
}

func directionDemo() {
	ch := make(chan int)
	go announce([]int{17, 18, 19}, ch)
	display(ch)
}

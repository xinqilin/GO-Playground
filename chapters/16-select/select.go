package main

import "fmt"

// readyDemo：select 同時等外帶和外送兩個 channel，哪一個有單就處理哪一個。
func readyDemo() {
	pickup := make(chan string, 1)
	delivery := make(chan string, 1)
	delivery <- "外送單 A-07"

	select {
	case o := <-pickup:
		fmt.Println("外帶", o)
	case o := <-delivery:
		fmt.Println("處理", o)
	}
}

// fairDemo：好幾個 case 同時都能進行時，select 隨機挑一個，不是由上往下挑第一個。
// 已經 close 的 channel 永遠收得到東西，所以這兩個 case 每一次都同時準備好。
func fairDemo() {
	a := make(chan struct{})
	b := make(chan struct{})
	close(a)
	close(b)

	na, nb := 0, 0
	for range 1000 {
		select {
		case <-a:
			na++
		case <-b:
			nb++
		}
	}
	fmt.Println(na > 0, nb > 0, na+nb)
}

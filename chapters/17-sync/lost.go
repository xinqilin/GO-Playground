package main

import "fmt"

// lostUpdateDemo 把兩個 goroutine 同時執行 served++ 的過程，一步一步排出來：
// served++ 其實是「讀出來、加一、寫回去」三步，兩個店員都先讀，再各自寫回，就少算一份。
func lostUpdateDemo() {
	served := 5
	a := served    // 店員 A 讀到 5
	b := served    // 店員 B 也讀到 5
	served = a + 1 // A 寫回 6
	served = b + 1 // B 也寫回 6
	fmt.Println("兩個店員各記一份，結果是", served)
}

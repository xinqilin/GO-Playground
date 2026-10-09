// 第 15 章：channel。
//
// channel 是 goroutine 之間的輸送帶：一邊放、一邊拿。不要靠共用變數溝通，要靠溝通來共用資料。
// 執行：go run ./chapters/15-channels
package main

import "fmt"

func main() {
	fmt.Println("== 無緩衝：一手交錢一手交貨 ==")
	ch := make(chan string)
	go func() { ch <- "ping" }() // 送出會卡住，直到有人來收
	fmt.Println(<-ch)

	fmt.Println("== 有緩衝：先放進信箱 ==")
	box := make(chan int, 2)
	box <- 1
	box <- 2 // 容量 2，放兩個還不會卡住；放第三個才會
	fmt.Println(len(box), cap(box))
	fmt.Println(<-box, <-box)

	fmt.Println("== close 與 range ==")
	for n := range produce(3) {
		fmt.Println("收到", n)
	}

	fmt.Println("== 關閉之後再收 ==")
	c := make(chan int, 1)
	c <- 7
	close(c)
	v, ok := <-c // 關閉前放進去的還收得到
	fmt.Println(v, ok)
	v, ok = <-c // 收完之後：零值和 false
	fmt.Println(v, ok)

	fmt.Println("== 管線：產生 → 平方 → 加總 ==")
	fmt.Println(sum(square(produce(4))))
}

// produce 依序送出 1 到 n，送完就關閉 channel。
// 回傳型別 <-chan int 代表呼叫端只能收、不能送。
func produce(n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out) // 送的一方負責關閉，range 才知道什麼時候結束
		for i := 1; i <= n; i++ {
			out <- i
		}
	}()
	return out
}

// square 從 in 收數字、算平方、送到 out。
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// sum 把 channel 裡的數字全部加起來，直到 channel 關閉。
func sum(in <-chan int) int {
	total := 0
	for n := range in {
		total += n
	}
	return total
}

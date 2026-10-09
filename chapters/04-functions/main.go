// 第 04 章：函式。
//
// 函式是一台機器：放進參數、拿出結果。在 Go 裡，函式本身也是一種值，可以傳來傳去。
// 執行：go run ./chapters/04-functions
package main

import "fmt"

func main() {
	fmt.Println("== 多個回傳值 ==")
	q, r := divmod(17, 5)
	fmt.Println(q, r)

	fmt.Println("== 可變參數 ==")
	fmt.Println(sum(), sum(1, 2), sum([]int{1, 2, 3}...))

	fmt.Println("== 函式是值 ==")
	double := func(n int) int { return n * 2 }
	fmt.Println(apply([]int{1, 2, 3}, double))

	fmt.Println("== closure ==")
	next := counter()
	fmt.Println(next(), next(), next())
	other := counter() // 另一個 closure，有自己的 n
	fmt.Println(other())

	fmt.Println("== defer ==")
	deferOrder()
}

// divmod 同時回傳商和餘數。
func divmod(a, b int) (int, int) {
	return a / b, a % b
}

// sum 接受任意個 int；在函式裡，nums 就是一個 []int。
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// apply 把 f 套用到每個元素上：函式可以當參數傳。
func apply(nums []int, f func(int) int) []int {
	out := make([]int, 0, len(nums))
	for _, n := range nums {
		out = append(out, f(n))
	}
	return out
}

// counter 回傳一個 closure。closure 記得它用到的變數 n，每次呼叫都會累加。
func counter() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}

// deferOrder 示範 defer：函式結束前才執行，後註冊的先執行；參數在註冊當下就決定。
func deferOrder() {
	for i := range 3 {
		defer fmt.Println("defer", i)
	}
	x := 10
	defer fmt.Println("註冊 defer 時 x 是", x)
	x = 20
	fmt.Println("函式本體結束，x =", x)
}

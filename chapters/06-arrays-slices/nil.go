package main

import "fmt"

// nilVsEmpty 示範 nil slice 和空 slice：長度都是 0，用起來一樣，只有和 nil 比較、用 %#v 印出時不同。
func nilVsEmpty() {
	var a []int
	b := []int{}
	c := make([]int, 0)
	fmt.Println(a == nil, b == nil, c == nil)
	fmt.Println(len(a), len(b), len(c))
	fmt.Printf("%v %v\n", a, b)
	fmt.Printf("%#v %#v\n", a, b)

	for range a {
		fmt.Println("不會執行")
	}
	a = append(a, 1)
	fmt.Println(a, len(a))
}

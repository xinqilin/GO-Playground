package main

import "fmt"

// addressAndDeref 示範 & 取得變數的位址，* 照著位址取值。
func addressAndDeref() {
	stock := 5
	p := &stock // p 的型別是 *int：「指向 int 的指標」
	fmt.Println("stock =", stock)
	fmt.Println("*p =", *p)
	fmt.Println("p == &stock:", p == &stock)
	fmt.Printf("%T\n", p)
}

// modifyThroughPointer 示範透過指標改動原本的變數。
func modifyThroughPointer() {
	stock := 5
	p := &stock
	*p = 10 // 改的是 p 指向的那一格，也就是 stock
	*p += 3
	fmt.Println("stock =", stock)
}

// aliasing 示範兩個指標指向同一個變數：從哪一個改，另一個都看得到。
func aliasing() {
	stock := 5
	a := &stock
	b := &stock
	*a = 8
	fmt.Println(*b, a == b)

	other := 8
	c := &other
	fmt.Println(*a == *c, a == c) // 值相等，但指向不同的變數
}

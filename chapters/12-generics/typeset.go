package main

import "fmt"

// NTD 是新台幣，底層型別是 int。
type NTD int

func (n NTD) String() string {
	return fmt.Sprintf("NT$%d", int(n))
}

// Number 是型別限制：底層型別是 int、int64 或 float64 的型別都可以。
type Number interface {
	~int | ~int64 | ~float64
}

// Total 和 Sum 一樣是加總，但限制用了 ~，所以 NTD 這種自訂型別也能用。
func Total[N Number](s []N) N {
	var total N
	for _, v := range s {
		total += v
	}
	return total
}

// typeSetDemo：回傳的型別就是傳進去的 NTD，所以 String 方法還在。
func typeSetDemo() {
	prices := []NTD{85, 95, 30}
	fmt.Println(Total(prices))
	fmt.Printf("%T\n", Total(prices))
	fmt.Println(Total([]int64{1, 2}))
}

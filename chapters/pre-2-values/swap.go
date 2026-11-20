package main

import "fmt"

// wrongSwap 想把兩杯飲料對調，但第一步就把紅茶倒掉了。
func wrongSwap() {
	a := "紅茶"
	b := "綠茶"
	a = b
	b = a
	fmt.Println(a, b)
}

// tempSwap 先把紅茶倒進第三個杯子暫放。
func tempSwap() {
	a := "紅茶"
	b := "綠茶"
	tmp := a
	a = b
	b = tmp
	fmt.Println(a, b)
}

// goSwap 一次指定兩個變數：右邊的值全部先拿好，才放進左邊。
func goSwap() {
	a := "紅茶"
	b := "綠茶"
	a, b = b, a
	fmt.Println(a, b)
}

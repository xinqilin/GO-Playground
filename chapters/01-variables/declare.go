package main

import (
	"fmt"
	"strconv"
)

// shopName 宣告在套件層級：同一個 package 的所有函式都看得到它。
// 套件層級只能用 var，不能用 :=。
var shopName = "阿姨便當"

// declareVar 示範 var 的三種寫法。
func declareVar() {
	var price int = 85 // 名稱、型別、初始值都寫出來
	var qty = 3        // 省略型別，由右邊的值推斷成 int
	var paid bool      // 省略初始值，得到零值 false
	fmt.Println(price, qty, paid)
}

// declareShort 示範簡短宣告 :=，這是函式裡最常用的寫法。
func declareShort() {
	item := "排骨便當"
	price := 85
	price = 90 // 之後要改值用 =，不是 :=
	fmt.Println(item, price)
}

// declareMulti 示範一次宣告好幾個變數，以及交換兩個變數。
func declareMulti() {
	var (
		opens  = 10
		closes = 20
	)
	name, qty := "雞腿便當", 2
	fmt.Println(opens, closes, name, qty)

	first, second := "A01", "A02"
	first, second = second, first // 右邊先全部算好，再一起指定，不需要暫存變數
	fmt.Println(first, second)
}

// redeclare 示範 := 左邊只要有一個新變數就可以，已經存在的變數會被重新指定。
func redeclare() {
	price, err := strconv.Atoi("85")
	fmt.Println(price, err)
	qty, err := strconv.Atoi("三") // qty 是新的，err 沿用上面那一個
	fmt.Println(qty, err)
}

// typeInference 印出 := 推斷出來的型別。%T 會印出值的型別。
func typeInference() {
	a := 85
	b := 85.0
	c := 'A'
	d := "A"
	e := true
	fmt.Printf("%T %T %T %T %T\n", a, b, c, d, e)
}

// packageLevel 使用套件層級的變數。
func packageLevel() {
	fmt.Println("店名：" + shopName)
}

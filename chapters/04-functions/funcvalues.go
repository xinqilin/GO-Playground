package main

import (
	"cmp"
	"fmt"
	"slices"
)

func double(n int) int { return n * 2 }

// funcValues 示範函式可以存在變數裡，也可以寫成沒有名字的函式字面值。
func funcValues() {
	f := double
	fmt.Println(f(21))

	square := func(n int) int { return n * n }
	fmt.Println(square(9))

	fmt.Printf("%T\n", square)

	var g func(int) int
	fmt.Println(g == nil)
}

// PriceRule 是一個函式型別：任何「收一個 int、回傳一個 int」的函式都屬於它。
type PriceRule func(int) int

// applyRules 依序把每條規則套用到價格上。
func applyRules(price int, rules []PriceRule) int {
	for _, rule := range rules {
		price = rule(price)
	}
	return price
}

// funcTypes 示範把函式放進 slice，再當成參數傳給另一個函式。
func funcTypes() {
	memberDiscount := func(p int) int { return p * 9 / 10 }
	coupon := func(p int) int { return p - 50 }
	rules := []PriceRule{memberDiscount, coupon}
	fmt.Println(applyRules(500, rules))
	fmt.Println(applyRules(500, []PriceRule{coupon, memberDiscount}))
}

type bento struct {
	name  string
	price int
}

// sortWithFunc 示範標準庫用函式當參數的例子：排序的比較方式由你決定。
func sortWithFunc() {
	menu := []bento{{"雞腿", 90}, {"素食", 70}, {"排骨", 85}}
	slices.SortFunc(menu, func(a, b bento) int {
		return cmp.Compare(a.price, b.price)
	})
	fmt.Println(menu)
}

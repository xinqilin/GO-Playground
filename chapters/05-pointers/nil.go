package main

import "fmt"

// nilPointer 示範指標的零值是 nil，對 nil 解參考會 panic。
// 這裡用 recover 把 panic 接住才能印出訊息（第 04 章）。
func nilPointer() {
	var p *int
	fmt.Println(p == nil)
	defer func() {
		fmt.Println("panic:", recover())
	}()
	fmt.Println(*p)
}

// priceAfter 回傳打折後的價格。discount 是 nil 代表「沒有設定折扣」，和「折扣 0%」不一樣。
func priceAfter(price int, discount *int) string {
	if discount == nil {
		return fmt.Sprintf("%d 元（沒有設定折扣）", price)
	}
	return fmt.Sprintf("%d 元（折扣 %d%%）", price*(100-*discount)/100, *discount)
}

// optionalField 示範用指標表示「可以沒有值」。
func optionalField() {
	zero := 0
	fmt.Println(priceAfter(100, nil))
	fmt.Println(priceAfter(100, &zero))
	fmt.Println(priceAfter(100, new(20)))
}

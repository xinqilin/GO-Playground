package main

import (
	"fmt"
	"strconv"
	"strings"
)

// divmod 同時回傳商和餘數。
func divmod(a, b int) (int, int) {
	return a / b, a % b
}

// multiReturn 示範接收多個回傳值，以及用 _ 丟掉不要的值。
func multiReturn() {
	q, r := divmod(17, 5)
	fmt.Println(q, r)
	_, rest := divmod(1000, 3)
	fmt.Println(rest)
}

// parsePrice 把 "85 元" 這種字串轉成整數。最後一個回傳值是 error，成功時是 nil。
func parsePrice(s string) (int, error) {
	num, found := strings.CutSuffix(s, " 元")
	if !found {
		return 0, fmt.Errorf("%q 沒有「 元」結尾", s)
	}
	n, err := strconv.Atoi(num)
	if err != nil {
		return 0, fmt.Errorf("%q 的金額不是整數", s)
	}
	return n, nil
}

// valueAndError 示範「結果加錯誤」的慣用寫法：呼叫完馬上檢查 err。
func valueAndError() {
	for _, s := range []string{"85 元", "八十五 元", "85"} {
		price, err := parsePrice(s)
		if err != nil {
			fmt.Println("錯誤：", err)
			continue
		}
		fmt.Println("價格：", price)
	}
}

// splitBill 的回傳值有名字：each 和 rest 在函式一開始就是零值，
// 不帶值的 return 會把它們當下的值回傳出去。
func splitBill(total, people int) (each, rest int) {
	if people < 1 {
		return // 回傳 0, 0
	}
	each = total / people
	rest = total % people
	return
}

// namedResults 示範具名回傳值。
func namedResults() {
	fmt.Println(splitBill(1000, 3))
	fmt.Println(splitBill(1000, 0))
}

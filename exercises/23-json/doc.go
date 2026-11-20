// Package takeout 是第 23 章「JSON 與資料編碼」的練習題。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/23-json/
//
// 只想跑某一題時加上 -run：
//
//	go test ./exercises/23-json/ -run TestMoney
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/23-json/
//
// 除了練習 5 用 encoding/csv，其他題目都用 encoding/json/v2。
package takeout

import "fmt"

// Order 是一筆外帶訂單（練習 1、2）。
type Order struct {
	ID       string `json:"id"`
	Customer string `json:"customer"`
	Items    []Line `json:"items"`
	Note     string `json:"note,omitempty"`
}

// Line 是訂單裡的一個品項。
type Line struct {
	Dish string `json:"dish"`
	Qty  int    `json:"qty"`
}

// Money 是新台幣金額，單位是元（練習 4）。
type Money int

// Receipt 是收據（練習 4）。
type Receipt struct {
	Order string `json:"order"`
	Total Money  `json:"total"`
}

// MenuRow 是從 CSV 讀進來的一道菜（練習 5）。
type MenuRow struct {
	Name     string
	Category string
	Price    int
}

// CupSize 是飲料的杯型（練習 7）。零值 UnknownSize 表示沒有指定。
type CupSize int

const (
	UnknownSize CupSize = iota
	MediumCup           // 中杯
	LargeCup            // 大杯
)

// MarshalText 把杯型寫成「中杯」「大杯」。
func (c CupSize) MarshalText() ([]byte, error) {
	switch c {
	case MediumCup:
		return []byte("中杯"), nil
	case LargeCup:
		return []byte("大杯"), nil
	}
	return nil, fmt.Errorf("杯型 %d 不正確", int(c))
}

// Drink 是一杯飲料（練習 7）。
type Drink struct {
	Name string  `json:"name"`
	Size CupSize `json:"size"`
}

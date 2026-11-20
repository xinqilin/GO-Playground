// Package shop 是第 09 章「自訂型別、struct 與方法」的練習題。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的方法，然後執行：
//
//	go test ./exercises/09-structs-methods/
//
// 只想跑某一題時加上 -run：
//
//	go test ./exercises/09-structs-methods/ -run TestDoubleBooked
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/09-structs-methods/
package shop

import "errors"

// Seat 是一個座位，例如 A 排 5 號是 Seat{Row: "A", No: 5}。
type Seat struct {
	Row string
	No  int
}

// Line 是購物車裡的一行：品名、單價、數量。
type Line struct {
	Name  string
	Price int
	Qty   int
}

// Cart 是購物車。零值就是一台空的購物車。
type Cart struct {
	Lines []Line
}

// LoggedCart 是會記錄每一次加入動作的購物車。
type LoggedCart struct {
	Cart
	Log []string
}

// Card 是一張儲值的交通卡。
type Card struct {
	Balance int
}

// ErrInsufficient 表示餘額或庫存不足。
var ErrInsufficient = errors.New("不足")

// ErrBadAmount 表示金額不符合規定。
var ErrBadAmount = errors.New("金額不符合規定")

// Inventory 是庫存表。零值就能直接使用。
type Inventory struct {
	stock map[string]int
}

// Counter 是計數器。
type Counter struct {
	n int
}

// Value 回傳目前的計數。
func (c *Counter) Value() int {
	return c.n
}

// Employee 是一位員工。
type Employee struct {
	Name   string
	Salary int
}

// Package orders 是第 11 章「錯誤處理」的練習題。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/11-errors/
//
// 只想跑某一題時加上 -run：
//
//	go test ./exercises/11-errors/ -run TestParsePickup
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/11-errors/
package orders

import "errors"

// 練習用到的 sentinel error。
var (
	ErrTimeFormat = errors.New("格式要是 HH:MM")
	ErrClosed     = errors.New("不在營業時間")
	ErrSoldOut    = errors.New("賣完了")
	ErrPermanent  = errors.New("不能重試")
)

// LimitError 表示某樣東西超過了每人限購的數量。
type LimitError struct {
	Item string
	Max  int
}

// Package inout 是第 02 章「輸入與輸出」的練習題。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/02-input-output/
//
// 一開始全部是紅燈，寫完一題就會綠一題。只想跑某一題時加上 -run：
//
//	go test ./exercises/02-input-output/ -run TestPercent
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/02-input-output/
package inout

// Money 是以「元」為單位的金額。它的 String 方法在 exercises.go。
type Money int

// Item 是明細上的一個品項。
type Item struct {
	Name  string
	Qty   int
	Price int
}

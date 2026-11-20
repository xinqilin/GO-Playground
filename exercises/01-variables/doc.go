// Package variables 是第 01 章「變數、常數與型別」的練習題。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/01-variables/
//
// 一開始全部是紅燈，寫完一題就會綠一題。只想跑某一題時加上 -run：
//
//	go test ./exercises/01-variables/ -run TestSplitBill
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/01-variables/
package variables

// Perm 是便當店收銀系統的員工權限。每個權限佔一個位元，一位員工可以同時有好幾個權限。
type Perm uint8

const (
	Checkout  Perm = 1 << iota // 收銀
	Refund                     // 退款
	EditPrice                  // 改價
	CloseDay                   // 關帳
)

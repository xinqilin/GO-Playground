// Package memo 是第 24 章「文字處理」的練習題。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/24-text/
//
// 只想跑某一題時加上 -run：
//
//	go test ./exercises/24-text/ -run TestMask
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/24-text/
package memo

// Review 是顧客的評論（練習 2）。Name 和 Comment 是使用者輸入的。
type Review struct {
	Name    string
	Stars   int
	Comment string
}

// Invoice 是從文字裡找出來的一張發票（練習 4）。
type Invoice struct {
	Number string // 例如 AB-12345678
	Date   string // 一律是 2006-01-02 的格式
	Amount int    // 新台幣，單位是元
}

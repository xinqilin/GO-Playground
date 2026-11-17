// Package text 是第 08 章「字串與 rune」的練習題。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/08-strings-runes/
//
// 一開始全部是紅燈，寫完一題就會綠一題。只想跑某一題時加上 -run：
//
//	go test ./exercises/08-strings-runes/ -run TestMaskName
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/08-strings-runes/
package text

import "unicode"

// MenuItem 是菜單上的一個品項。
type MenuItem struct {
	Name  string
	Price int
}

// RuneWidth 回傳一個字在等寬字型裡大約佔幾格：中日韓文字和全形符號佔 2 格，其他佔 1 格。
// 和第 08 章 width.go 的 runeWidth 相同，練習 6 可以直接使用。
func RuneWidth(r rune) int {
	switch {
	case unicode.Is(unicode.Han, r),
		unicode.In(r, unicode.Hiragana, unicode.Katakana, unicode.Hangul),
		r >= 0x3000 && r <= 0x303F,
		r >= 0xFF01 && r <= 0xFF60,
		r >= 0xFFE0 && r <= 0xFFE6:
		return 2
	}
	return 1
}

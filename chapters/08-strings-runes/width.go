package main

import (
	"fmt"
	"strings"
	"unicode"
)

// runeWidth 回傳一個字在等寬字型裡大約佔幾格：中日韓文字和全形符號佔 2 格，其他佔 1 格。
// 這是簡化的版本，涵蓋常見的中文情況；完整的規則（Unicode 的 East Asian Width）
// 要用標準庫以外的套件，例如 golang.org/x/text/width。
func runeWidth(r rune) int {
	switch {
	case unicode.Is(unicode.Han, r),
		unicode.In(r, unicode.Hiragana, unicode.Katakana, unicode.Hangul),
		r >= 0x3000 && r <= 0x303F, // 中日韓的標點，例如「、」「。」「「」
		r >= 0xFF01 && r <= 0xFF60, // 全形的英數字和標點，例如「，」「！」
		r >= 0xFFE0 && r <= 0xFFE6: // 全形的貨幣符號等
		return 2
	}
	return 1
}

// displayWidth 回傳字串在等寬字型裡大約佔幾格。
func displayWidth(s string) int {
	w := 0
	for _, r := range s {
		w += runeWidth(r)
	}
	return w
}

// padRight 在 s 後面補空白，讓顯示寬度變成 width；已經超過時原樣回傳。
func padRight(s string, width int) string {
	if gap := width - displayWidth(s); gap > 0 {
		return s + strings.Repeat(" ", gap)
	}
	return s
}

// alignedMenu 示範 %-8s 依字數補空白，中文會對不齊；padRight 依顯示寬度補空白才對得齊。
func alignedMenu() {
	items := []string{"Coffee", "排骨飯", "雞腿便當"}
	for _, it := range items {
		fmt.Printf("|%-8s|\n", it)
	}
	for _, it := range items {
		fmt.Printf("|%s|\n", padRight(it, 8))
	}
}

package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// parseNumbers 示範把字串轉成數字。格式不對或超出範圍時回傳 error。
func parseNumbers() {
	n, err := strconv.Atoi("85")
	fmt.Println(n+5, err)
	_, err = strconv.Atoi("85元")
	fmt.Println(err)
	f, _ := strconv.ParseFloat("3.5", 64)
	fmt.Println(f * 2)
	big, err := strconv.ParseInt("99999999999", 10, 32)
	fmt.Println(big, err)
	h, _ := strconv.ParseInt("ff", 16, 64)
	ok, _ := strconv.ParseBool("true")
	fmt.Println(h, ok)
}

// formatNumbers 示範把數字轉成字串。
func formatNumbers() {
	fmt.Println(strconv.Itoa(85) + "元")
	fmt.Println(strconv.FormatInt(255, 2), strconv.FormatInt(255, 16))
	fmt.Println(strconv.FormatFloat(3.14159, 'f', 2, 64))
	fmt.Println(strconv.Quote("排骨\n飯"), strconv.QuoteToASCII("排骨"))
}

// classify 示範 unicode 套件判斷字元的種類。
func classify() {
	for _, r := range "A9語 ，" {
		fmt.Printf("%q 字母=%v 數字=%v 漢字=%v 空白=%v 標點=%v\n",
			r, unicode.IsLetter(r), unicode.IsDigit(r), unicode.Is(unicode.Han, r), unicode.IsSpace(r), unicode.IsPunct(r))
	}
}

// toHalfwidth 把全形的英數字和符號轉成半形，例如有人在表單裡輸入「ＡＢ１２３」。
// 全形的 ！（U+FF01）到 ～（U+FF5E）和半形的 ! 到 ~ 一一對應，碼位剛好差 0xFEE0。
func toHalfwidth(s string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r >= '！' && r <= '～':
			return r - 0xFEE0
		case r == '　': // 全形空白
			return ' '
		}
		return r
	}, s)
}

// halfwidth 示範使用 toHalfwidth。
func halfwidth() {
	fmt.Println(toHalfwidth("ＡＢ－１２３４５６７８"))
}

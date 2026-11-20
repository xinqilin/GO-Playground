package main

import (
	"fmt"
	"unicode/utf8"
)

// bytesVsRunes 示範 len 算的是 byte 數；要算字數用 utf8.RuneCountInString。
func bytesVsRunes() {
	s := "Go語言"
	fmt.Println(len(s), utf8.RuneCountInString(s))
	fmt.Println(s[0], s[2]) // 71 是 'G'；232 只是「語」的第一個 byte
	fmt.Printf("% x\n", s)  // 每個 byte 用十六進位印出來
}

// runeLiteral 示範 rune 是 int32 的別名，存的是一個字的 Unicode 碼位。
func runeLiteral() {
	r := '語'
	fmt.Println(r)
	fmt.Printf("%c %U %T\n", r, r, r)
	fmt.Println(utf8.RuneLen(r), utf8.RuneLen('G'))
	fmt.Println(string(r) + "言")
}

// cutInMiddle 示範用 byte 位置切字串，可能把一個中文字切成兩半，得到不合法的 UTF-8。
func cutInMiddle() {
	s := "Go語言"
	fmt.Printf("%q %v\n", s[:3], utf8.ValidString(s[:3]))
	fmt.Printf("%q %v\n", s[:5], utf8.ValidString(s[:5]))
}

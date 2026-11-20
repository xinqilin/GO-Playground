package main

import "fmt"

// rangeString 示範 range 字串：i 是這個字開始的 byte 位置，r 是一整個字。
func rangeString() {
	for i, r := range "Go語言" {
		fmt.Println(i, string(r), r)
	}
}

// rangeInvalid 示範遇到不合法的 byte 時，range 拿到 U+FFFD（替代字元），而且只前進 1 個 byte。
func rangeInvalid() {
	for i, r := range "A\xffB" {
		fmt.Printf("%d %U\n", i, r)
	}
}

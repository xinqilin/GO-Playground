package main

import (
	"fmt"
	"strconv"
)

// conversions 示範字串和 []byte、[]rune 之間的轉換。每次轉換都會複製一份資料。
func conversions() {
	s := "Go語言"
	b := []byte(s)
	r := []rune(s)
	fmt.Println(len(b), len(r))
	fmt.Println(b[:2], r[2:])
	fmt.Println(string(r[2]), string(r[2:]))
	b[0] = 'g' // 改的是複本
	fmt.Println(s, string(b))
}

// replaceAt 把第 i 個字換成 c。字串不能改，所以轉成 []rune，改完再轉回字串。
func replaceAt(s string, i int, c rune) string {
	r := []rune(s)
	r[i] = c
	return string(r)
}

// editRunes 示範使用 replaceAt。
func editRunes() {
	fmt.Println(replaceAt("大杯紅茶", 0, '中'))
}

// numberToString 示範數字轉成十進位字串要用 strconv.Itoa；string(rune(n)) 是把 n 當成碼位。
func numberToString() {
	n := 65
	fmt.Println(strconv.Itoa(n), string(rune(n)))
	fmt.Println(fmt.Sprint(n) + "元")
}

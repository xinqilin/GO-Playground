// 第 08 章：字串、byte 與 rune。
//
// 字串是一排唯讀的 byte，內容用 UTF-8 編碼：英文字佔 1 個 byte，中文字佔 3 個。
// 執行：go run ./chapters/08-strings-runes
package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func main() {
	fmt.Println("== 字串是 byte 序列 ==")
	bytesVsRunes()

	fmt.Println("== range 一次拿一個 rune ==")
	rangeString("Go語言")

	fmt.Println("== strings 套件 ==")
	stringsPkg()

	fmt.Println("== 組字串用 Builder ==")
	fmt.Println(repeatJoin("Go", 3))

	fmt.Println("== 數字與字串互轉 ==")
	conversions()
}

// bytesVsRunes 示範 len 算的是 byte 數，用索引拿到的也是 byte。
func bytesVsRunes() {
	s := "Go語言"
	fmt.Println(len(s))                    // 8：中文字在 UTF-8 佔 3 個 byte
	fmt.Println(utf8.RuneCountInString(s)) // 4：字數
	fmt.Println(s[0], s[2])                // 71 是 'G'；232 只是「語」的第一個 byte
	r := []rune(s)
	fmt.Println(string(r[2])) // 轉成 []rune 之後才能按字取
}

// rangeString 示範 range 字串：i 是 byte 位置，r 是一整個字（rune）。
func rangeString(s string) {
	for i, r := range s {
		fmt.Println(i, string(r), r)
	}
}

// stringsPkg 示範 strings 套件的常用函式。
func stringsPkg() {
	t := strings.TrimSpace("  Hello, Go  ")
	fmt.Println(strings.ToUpper(t), strings.Contains(t, "Go"), strings.HasPrefix(t, "He"))
	fmt.Println(strings.Split("a,b,c", ","), strings.Fields(" a  b "))
	key, value, found := strings.Cut("lang=go", "=")
	fmt.Println(key, value, found)
}

// repeatJoin 用 strings.Builder 組字串。字串不能修改，在迴圈裡用 + 串接，每次都會產生新字串。
func repeatJoin(word string, n int) string {
	var b strings.Builder
	for i := range n {
		if i > 0 {
			b.WriteString("-")
		}
		b.WriteString(word)
	}
	return b.String()
}

// conversions 示範 strconv：數字和字串之間的轉換。
func conversions() {
	n, err := strconv.Atoi("123")
	fmt.Println(n+1, err)
	_, err = strconv.Atoi("12a")
	fmt.Println(err)
	f, err := strconv.ParseFloat("3.5", 64)
	fmt.Println(f*2, err)
	fmt.Println(strconv.Itoa(42)+"!", strconv.Quote("你好\n"))
	raw := `C:\path\no\escape` // 反引號是原始字串，不處理跳脫字元
	fmt.Println(raw)
}

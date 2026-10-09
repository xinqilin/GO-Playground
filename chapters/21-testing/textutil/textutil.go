// Package textutil 是第 21 章用來寫測試的小套件。
package textutil

import "strings"

// Reverse 把字串反轉。先轉成 []rune 再反轉，中文字才不會被切壞。
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// IsPalindrome 判斷 s 去掉空白、不分大小寫之後，正著讀和反著讀是否一樣。
func IsPalindrome(s string) bool {
	s = strings.ToLower(strings.Join(strings.Fields(s), ""))
	return s == Reverse(s)
}

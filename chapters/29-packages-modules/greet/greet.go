// Package greet 示範一個小套件：大寫開頭的名稱可以給其他套件用，小寫開頭的只在套件內部可見。
package greet

import "strings"

// Prefix 是匯出的變數。
var Prefix = "哈囉"

// Hello 是匯出的函式，其他套件用 greet.Hello 呼叫。
func Hello(name string) string {
	return Prefix + "，" + normalize(name)
}

// normalize 沒有匯出，只有 greet 套件自己能呼叫。
func normalize(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "陌生人"
	}
	return name
}

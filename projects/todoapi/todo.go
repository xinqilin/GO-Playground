package main

import (
	"errors"
	"maps"
	"slices"
	"strings"
	"time"
)

// Todo 是一筆待辦事項，也是 API 回應和存檔的格式。
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	Due       string    `json:"due,omitzero"` // 期限，格式是 2006-01-02；空字串代表沒有期限
	CreatedAt time.Time `json:"created_at"`
}

// ErrNotFound 表示找不到指定的待辦事項。
var ErrNotFound = errors.New("找不到這筆待辦")

// ValidationError 記錄每個欄位的問題。一次把全部的問題告訴客戶端，
// 而不是改好一個欄位、送出去，才發現下一個。
type ValidationError map[string]string

func (e ValidationError) Error() string {
	var b strings.Builder
	for i, field := range slices.Sorted(maps.Keys(e)) {
		if i > 0 {
			b.WriteString("；")
		}
		b.WriteString(field + "：" + e[field])
	}
	return b.String()
}

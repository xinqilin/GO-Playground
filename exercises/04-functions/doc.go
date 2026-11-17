// Package functions 是第 04 章「函式」的練習題。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/04-functions/
//
// 一開始全部是紅燈，寫完一題就會綠一題。只想跑某一題時加上 -run：
//
//	go test ./exercises/04-functions/ -run TestMinMax
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/04-functions/
package functions

// Rule 是一條價格規則：收到原價，回傳套用規則之後的價格。
type Rule func(price int) int

// Tracker 記錄資源的開啟、使用和關閉，用來檢查順序對不對。
type Tracker struct {
	Events []string
}

// Open 記錄開啟 name。
func (t *Tracker) Open(name string) { t.Events = append(t.Events, "open "+name) }

// Use 記錄使用 name。
func (t *Tracker) Use(name string) { t.Events = append(t.Events, "use "+name) }

// Close 記錄關閉 name。
func (t *Tracker) Close(name string) { t.Events = append(t.Events, "close "+name) }

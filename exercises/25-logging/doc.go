// Package journal 是第 25 章「日誌」的練習題。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/25-logging/
//
// 只想跑某一題時加上 -run：
//
//	go test ./exercises/25-logging/ -run TestRedact
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/25-logging/
package journal

// Summary 是一份 log 的統計（練習 3）。
type Summary struct {
	Counts map[string]int // 每個 level 字串（例如 "INFO"、"WARN+2"）出現幾筆
	Errors []string       // ERROR 以上（含 ERROR+1 這類）每一筆的 msg，依照出現的順序
}

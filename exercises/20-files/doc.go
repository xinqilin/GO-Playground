// Package cabinet 是第 20 章「檔案與 I/O」的練習題。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/20-files/
//
// 只想跑某一題時加上 -run：
//
//	go test ./exercises/20-files/ -run TestMerge
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/20-files/
//
// 測試只在 t.TempDir() 開的暫存目錄裡讀寫，測試結束時會自動刪掉，不會碰到其他檔案。
package cabinet

import "errors"

// ErrTooLarge 表示上傳的內容超過上限（練習 5）。
var ErrTooLarge = errors.New("上傳的檔案太大")

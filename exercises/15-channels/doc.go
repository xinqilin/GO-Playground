// Package relay 是第 15 章「channel」的練習題。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/15-channels/
//
// 只想跑某一題時加上 -run：
//
//	go test ./exercises/15-channels/ -run TestTickets
//
// 寫完之後也用 -race 跑一次，確認沒有 data race：
//
//	go test -race ./exercises/15-channels/
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/15-channels/
//
// 回傳 channel 的題目，題目檔先回傳一個已經 close 的 channel，讓測試不會卡住。
// 測試收 channel 時有逾時保護：忘了 close，測試會在一秒後失敗，而不是永遠等下去。
//
// TestFirstOK 用 testing/synctest 檢查 goroutine 洩漏。有 goroutine 卡著沒結束時，
// 整個測試程式會 panic，訊息是 "blocked goroutines remain"，後面列出卡住的 goroutine 停在哪一行。
package relay

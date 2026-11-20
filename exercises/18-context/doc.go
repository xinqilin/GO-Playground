// Package courier 是第 18 章「context」的練習題。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/18-context/
//
// 只想跑某一題時加上 -run：
//
//	go test ./exercises/18-context/ -run TestSleep
//
// 寫完之後也用 -race 跑一次，確認沒有 data race：
//
//	go test -race ./exercises/18-context/
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/18-context/
//
// 和時間有關的題目，測試用 testing/synctest 的假時間執行：time.Sleep、time.After 和
// context.WithTimeout 都不會真的等，而且結果每次都一樣。泡泡結束時還有 goroutine 卡著，
// 整個測試程式會 panic，訊息是 "blocked goroutines remain"，後面列出卡住的 goroutine 停在哪一行。
package courier

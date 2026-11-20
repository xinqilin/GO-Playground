// Package rush 是第 19 章「並行模式」的練習題。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test -race ./exercises/19-concurrency-patterns/
//
// 只想跑某一題時加上 -run：
//
//	go test -race ./exercises/19-concurrency-patterns/ -run TestLimit
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -race -tags solution ./exercises/19-concurrency-patterns/
//
// 和時間有關的題目，測試用 testing/synctest 的假時間執行，可以精確檢查花了多久、
// 同時有幾個工作在跑。泡泡結束時還有 goroutine 卡著，整個測試程式會 panic，
// 訊息是 "blocked goroutines remain"，後面列出卡住的 goroutine 停在哪一行。
package rush

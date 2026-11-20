// Package lineup 是第 13 章「疊代器」的練習題。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/13-iterators/
//
// 只想跑某一題時加上 -run：
//
//	go test ./exercises/13-iterators/ -run TestEvery
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/13-iterators/
//
// 每一題的疊代器都要遵守 yield 的約定：yield 回傳 false 之後，立刻停止，不可以再呼叫 yield。
package lineup

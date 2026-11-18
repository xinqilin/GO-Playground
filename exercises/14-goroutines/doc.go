// Package crew 是第 14 章「goroutine」的練習題。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/14-goroutines/
//
// 只想跑某一題時加上 -run：
//
//	go test ./exercises/14-goroutines/ -run TestEach
//
// 寫完之後也用 -race 跑一次，確認沒有 data race：
//
//	go test -race ./exercises/14-goroutines/
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/14-goroutines/
//
// 有幾題的測試會檢查工作是不是真的「同時」在跑：一個接一個做的寫法會等到逾時才失敗，所以會慢個一兩秒。
package crew

// Branch 是一家分店，底下可能還有它負責的其他分店。
type Branch struct {
	Name     string
	Sales    int
	Children []*Branch
}

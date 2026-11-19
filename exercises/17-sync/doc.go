// Package stall 是第 17 章「sync 與 data race」的練習題。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test -race ./exercises/17-sync/
//
// 這一章的題目一定要加 -race：少了鎖的寫法，常常不加 -race 也能通過測試。
// 只想跑某一題時加上 -run：
//
//	go test -race ./exercises/17-sync/ -run TestJar
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -race -tags solution ./exercises/17-sync/
//
// 好幾個 goroutine 同時寫同一個 map 時，執行時期可能直接讓整個測試程式結束，
// 訊息是 "fatal error: concurrent map writes"，其他題目的結果也就看不到了。
package stall

import (
	"errors"
	"sync"
)

// ErrTaken 表示這張桌子已經有人訂了。
var ErrTaken = errors.New("這張桌子已經有人訂了")

// Settings 是店裡目前的設定。
type Settings struct {
	Open     bool
	Discount int // 打幾折，0 表示不打折
	Served   int // 今天賣出幾份
}

// Queue 是候位名單，可以給多個 goroutine 同時使用，零值就能直接用。
type Queue struct {
	mu    sync.Mutex
	names []string
}

// Join 把 name 排到最後面。
func (q *Queue) Join(name string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.names = append(q.names, name)
}

// Len 回傳目前有幾組在等。
func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.names)
}

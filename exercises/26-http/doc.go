// Package ticket 是第 26 章「HTTP 服務」的練習題：排隊叫號機的 API。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/26-http/
//
// 只想跑某一題時加上 -run：
//
//	go test ./exercises/26-http/ -run TestRequireKey
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/26-http/
package ticket

import (
	"slices"
	"sync"
)

// Queue 是叫號機：取號、查前面還有幾個人、叫下一號。零值就能用，可以被很多 goroutine 同時使用。
// 練習 5 直接用它，不用修改。
type Queue struct {
	mu      sync.Mutex
	last    int   // 最後發出去的號碼
	waiting []int // 還在等的號碼，依照順序
}

// Take 發一張新的號碼牌，回傳號碼和前面還有幾個人。號碼從 1 開始。
func (q *Queue) Take() (number, ahead int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.last++
	q.waiting = append(q.waiting, q.last)
	return q.last, len(q.waiting) - 1
}

// Position 回傳 number 前面還有幾個人；這個號碼不在等候中（還沒發、已經叫過）時 ok 是 false。
func (q *Queue) Position(number int) (ahead int, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	i := slices.Index(q.waiting, number)
	return i, i >= 0
}

// Call 叫下一號；沒有人在等時 ok 是 false。
func (q *Queue) Call() (number int, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.waiting) == 0 {
		return 0, false
	}
	number = q.waiting[0]
	q.waiting = q.waiting[1:]
	return number, true
}

// Settings 是叫號機的設定（練習 6）。可以被很多 goroutine 同時使用。
type Settings struct {
	mu   sync.Mutex
	name string
}

// NewSettings 建立店名是 name 的設定。
func NewSettings(name string) *Settings {
	return &Settings{name: name}
}

// Name 回傳目前的店名。
func (s *Settings) Name() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.name
}

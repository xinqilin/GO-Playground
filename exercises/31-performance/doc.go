// Package rush 是第 31 章「效能與 profiling」的練習題：午餐尖峰時段，便當店的程式要又快又省。
//
// 除了結果要對，有幾題的測試還會用 testing.AllocsPerRun 檢查配置了幾次記憶體。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/31-performance/
//
// 只想跑某一題時加上 -run：
//
//	go test ./exercises/31-performance/ -run TestReceipt
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/31-performance/
package rush

import (
	"errors"
	"sync"
)

// ErrBadLine 表示訂單紀錄的格式不對。
var ErrBadLine = errors.New("訂單格式不對")

// Order 是訂單紀錄的一行：「12:05,排骨飯,2,90」。
type Order struct {
	Time  string
	Dish  string
	Qty   int
	Price int
}

// Item 是收據上的一項。
type Item struct {
	Dish  string
	Qty   int
	Price int // 單價
}

// BufferPool 重複使用 bytes.Buffer。零值可以直接用。
type BufferPool struct {
	MaxCap int // Cap() 超過 MaxCap 的緩衝區不放回去；0 表示不限制

	pool sync.Pool
}

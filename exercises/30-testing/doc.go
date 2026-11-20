// Package testkit 是第 30 章「測試」的練習題：自己做幾個寫測試時會用到的小工具。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/30-testing/
//
// 只想跑某一題時加上 -run：
//
//	go test ./exercises/30-testing/ -run TestShrink
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/30-testing/
package testkit

import (
	"errors"
	"sync"
)

// ErrUnreachable 是 Recorder 對 FailOn 裡的電話號碼回傳的錯誤。
var ErrUnreachable = errors.New("聯絡不到")

// Recorder 是測試用的假簡訊服務：記下每一則送出的簡訊。
// FailOn 裡的電話號碼會收到 ErrUnreachable，而且不會被記下。
// 建立之後 FailOn 不會再被修改；Notify 會被好幾個 goroutine 同時呼叫。
type Recorder struct {
	FailOn map[string]bool

	mu   sync.Mutex
	sent []string
}

// Cleanups 收集清理函式，用法和 t.Cleanup 一樣。
type Cleanups struct {
	fns []func()
}

// Add 加入一個清理函式。
func (c *Cleanups) Add(f func()) {
	c.fns = append(c.fns, f)
}

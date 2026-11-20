// Package menu 是第 12 章「泛型」的練習題。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/12-generics/
//
// 只想跑某一題時加上 -run：
//
//	go test ./exercises/12-generics/ -run TestFilter
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/12-generics/
package menu

import "cmp"

// Number 是數字型別的限制：底層型別是 int、int64 或 float64 的型別都可以。
type Number interface {
	~int | ~int64 | ~float64
}

// Tally 計算每個 key 出現了幾次。零值就能直接用。
type Tally[K cmp.Ordered] struct {
	counts map[K]int
}

// Result 是一個可能失敗的計算結果：成功時 Err 是 nil，值在 Val；失敗時 Err 不是 nil。
type Result[T any] struct {
	Val T
	Err error
}

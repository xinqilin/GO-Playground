// Package modver 是第 29 章「package 與 module」的練習題：自己算一次 go 指令替你做的版本決定。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/29-packages-modules/
//
// 只想跑某一題時加上 -run：
//
//	go test ./exercises/29-packages-modules/ -run TestSelect
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/29-packages-modules/
package modver

import "fmt"

// Version 是語意化版本 vMAJOR.MINOR.PATCH，可以帶一個預發行標記：
// v1.2.0-rc.1 的 Pre 是 "rc.1"，正式版的 Pre 是空字串。
type Version struct {
	Major, Minor, Patch int
	Pre                 string
}

// String 把版本寫回 v1.2.3 或 v1.2.3-rc.1 的格式。
func (v Version) String() string {
	s := fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Pre != "" {
		s += "-" + v.Pre
	}
	return s
}

// Require 是 go.mod 裡的一筆 require。
type Require struct {
	Path     string
	Version  string
	Indirect bool // 行尾有 // indirect
}

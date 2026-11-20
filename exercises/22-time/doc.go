// Package timetable 是第 22 章「時間」的練習題。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test -race ./exercises/22-time/
//
// 只想跑某一題時加上 -run：
//
//	go test -race ./exercises/22-time/ -run TestBillingDates
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -race -tags solution ./exercises/22-time/
//
// 和計時器有關的題目，測試用 testing/synctest 的假時間執行，15 分鐘的逾時一瞬間就跑完。
package timetable

import (
	"time"
	_ "time/tzdata" // 時區資料庫編進來，每個系統都能載入 Asia/Taipei
)

// Taipei 是台北的時區。題目裡說的「幾點」都是台北時間。
var Taipei = func() *time.Location {
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		panic(err)
	}
	return loc
}()

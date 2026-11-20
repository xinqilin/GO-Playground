package main

import (
	"fmt"
	"time"
)

// monthEnd 回傳 t 那個月的最後一天，時刻和時區不變。
// time.Date 會自動進位：下個月的第 0 天，就是這個月的最後一天。
func monthEnd(t time.Time) time.Time {
	y, m, _ := t.Date()
	h, mi, s := t.Clock()
	return time.Date(y, m+1, 0, h, mi, s, t.Nanosecond(), t.Location())
}

// calendarDemo 會員從 1 月 31 日開始，一個月後到期：AddDate 的結果不是 2 月底。
func calendarDemo() {
	joined := time.Date(2026, time.January, 31, 12, 0, 0, 0, taipei)

	fmt.Println("加一個月：", joined.AddDate(0, 1, 0).Format(time.DateOnly))
	fmt.Println("2 月 31 日：", time.Date(2026, time.February, 31, 0, 0, 0, 0, taipei).Format(time.DateOnly))
	fmt.Println("2 月的最後一天：", monthEnd(time.Date(2026, time.February, 1, 0, 0, 0, 0, taipei)).Format(time.DateOnly))

	today := time.Date(2026, time.November, 19, 0, 0, 0, 0, taipei)
	xmas := time.Date(2026, time.December, 25, 0, 0, 0, 0, taipei)
	fmt.Println("離聖誕節還有", int(xmas.Sub(today).Hours()/24), "天")
}

// dstDemo 紐約在 2026-11-01 凌晨結束夏令時間，時鐘撥回一小時，那一天有 25 個小時。
// AddDate 加的是「日曆上的一天」，Add 加的是固定的 24 小時。
func dstDemo() {
	ny := mustLoad("America/New_York")
	t := time.Date(2026, time.October, 31, 12, 0, 0, 0, ny)
	const layout = time.DateTime + " MST"

	fmt.Println("AddDate：", t.AddDate(0, 0, 1).Format(layout))
	fmt.Println("Add 24h：", t.Add(24*time.Hour).Format(layout))
	fmt.Println("那一天其實有", t.AddDate(0, 0, 1).Sub(t))
}

// truncateDemo Truncate(24*time.Hour) 是用 UTC 的整天去切，不是你所在時區的午夜。
func truncateDemo() {
	t := time.Date(2026, time.November, 19, 7, 45, 0, 0, taipei) // 台北早上 7:45
	fmt.Println("Truncate：", t.Truncate(24*time.Hour).Format(time.DateTime))

	y, m, d := t.Date()
	fmt.Println("當地午夜：", time.Date(y, m, d, 0, 0, 0, 0, t.Location()).Format(time.DateTime))
}

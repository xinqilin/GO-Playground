package main

import (
	"fmt"
	"time"
)

// weekdays 是中文的星期，索引就是 time.Weekday（星期日是 0）。
var weekdays = [...]string{"日", "一", "二", "三", "四", "五", "六"}

// receiptStamp 把時間印成發票上的樣子，例如 2026/11/19（四）18:30。
// Weekday 的 String 是英文，中文要自己對照。
func receiptStamp(t time.Time) string {
	return t.Format("2006/01/02") + "（" + weekdays[t.Weekday()] + "）" + t.Format("15:04")
}

// layoutDemo 用同一個時間點試幾種版面。版面裡的數字是一個固定的參考時間：
// 2006 年 1 月 2 日 15 點 04 分 05 秒，時區 -0700。
func layoutDemo() {
	t := time.Date(2026, time.November, 19, 18, 30, 5, 0, time.UTC)
	fmt.Println(t.Format(time.DateTime))
	fmt.Println(t.Format(time.RFC3339))
	fmt.Println(t.Format("01/02 15:04"), t.Format("1/2 3:04PM"), t.Format("Mon Jan 2"))
	fmt.Println(receiptStamp(t))
	fmt.Println(t.Format("YYYY-MM-DD hh:mm")) // 其他語言的寫法：沒有一個是 Go 認得的數字
}

// parseDemo 把字串解析成時間：版面要和字串的樣子完全一樣。
func parseDemo() {
	const layout = "2006/01/02 15:04"
	t, err := time.Parse(layout, "2026/11/19 18:30")
	fmt.Println(t, err)

	_, err = time.Parse(layout, "2026-11-19 18:30") // 分隔符號不一樣
	fmt.Println(err)

	_, err = time.Parse(time.DateOnly, "2026-02-30")
	fmt.Println(err)
}

//go:build !solution

package timetable

import "time"

// Remaining 回傳離打烊還有多久，給店門口的看板用：
//   - now 不在 closing 之前："已經打烊"
//   - 剩不到 1 分鐘："不到 1 分鐘"
//   - 不到 1 小時："還有 M 分"
//   - 其他："還有 H 小時 M 分"
//
// 分鐘數無條件捨去。now 和 closing 可能在不同的時區。
func Remaining(now, closing time.Time) string {
	// TODO: 實作
	return ""
}

// IsOpen 回傳 t 這個時刻店有沒有開。營業時間（台北時間）：
// 午餐 10:30 到 14:00，晚餐 17:00 到 20:30，開始的那一分鐘算營業、結束的那一分鐘不算；
// 星期一公休。t 可能在任何時區。
func IsOpen(t time.Time) bool {
	// TODO: 實作
	return false
}

// BillingDates 回傳從 first 開始、每個月一次、共 n 次的扣款時間。
// 每次都在 first 的那個「日」扣款；那個月沒有這一天時，改在那個月的最後一天。
// 時刻和時區都和 first 一樣。n 小於 1 時回傳空的 slice。
func BillingDates(first time.Time, n int) []time.Time {
	// TODO: 實作
	return nil
}

// Holds 管理保留中的座位：保留到期還沒確認，就呼叫 onExpire 把座位放出來。
// 所有方法都可能被不同的 goroutine 同時呼叫。
type Holds struct {
	// TODO: 加上需要的欄位
}

// NewHolds 建立一個 Holds，座位逾時時呼叫 onExpire。
func NewHolds(onExpire func(seat string)) *Holds {
	// TODO: 實作
	return &Holds{}
}

// Hold 保留 seat，d 之後逾時。seat 已經在保留中時，從現在重新計時。
func (h *Holds) Hold(seat string, d time.Duration) {
	// TODO: 實作
}

// Confirm 確認 seat：還在保留中就取消計時、不再保留，回傳 true；
// 已經逾時、或根本沒有保留時回傳 false。
func (h *Holds) Confirm(seat string) bool {
	// TODO: 實作
	return false
}

// Len 回傳還在保留中的座位數。
func (h *Holds) Len() int {
	// TODO: 實作
	return 0
}

// CountByDay 依照台北的日期（格式 2006-01-02）統計每天有幾筆訂單。
//
// 找錯：台北早上 8 點以前的訂單，都被算到前一天去了。
func CountByDay(orders []time.Time) map[string]int {
	counts := map[string]int{}
	for _, t := range orders {
		day := t.In(Taipei).Truncate(24 * time.Hour)
		counts[day.Format(time.DateOnly)]++
	}
	return counts
}

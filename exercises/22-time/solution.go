//go:build solution

package timetable

import (
	"fmt"
	"sync"
	"time"
)

// Remaining 先處理打烊和不到一分鐘的情況，再把剩下的時間截到分鐘。
func Remaining(now, closing time.Time) string {
	left := closing.Sub(now)
	switch {
	case left <= 0:
		return "已經打烊"
	case left < time.Minute:
		return "不到 1 分鐘"
	}
	h := int(left.Hours())
	m := int(left.Minutes()) % 60
	if h == 0 {
		return fmt.Sprintf("還有 %d 分", m)
	}
	return fmt.Sprintf("還有 %d 小時 %d 分", h, m)
}

// IsOpen 先換成台北時間，再把時刻換成「從午夜起算的分鐘數」比較。
func IsOpen(t time.Time) bool {
	t = t.In(Taipei)
	if t.Weekday() == time.Monday {
		return false
	}
	m := t.Hour()*60 + t.Minute()
	return (m >= 10*60+30 && m < 14*60) || (m >= 17*60 && m < 20*60+30)
}

// BillingDates 每一次都從 first 算月份，不從上一次的結果往後加：
// 1/31 → 2/28 之後，下一次要回到 3/31，而不是 3/28。
func BillingDates(first time.Time, n int) []time.Time {
	y, m, d := first.Date()
	h, mi, s := first.Clock()
	dates := make([]time.Time, 0, max(n, 0))
	for i := range n {
		last := time.Date(y, m+time.Month(i)+1, 0, 0, 0, 0, 0, first.Location()).Day()
		dates = append(dates, time.Date(y, m+time.Month(i), min(d, last), h, mi, s, first.Nanosecond(), first.Location()))
	}
	return dates
}

// hold 是一筆保留。計時器的 callback 用它的指標確認「我還是目前的那一筆」。
type hold struct {
	timer *time.Timer
}

// Holds 用 Mutex 保護 map；逾時的 callback 和 Confirm 誰先拿到鎖，誰就決定結果。
type Holds struct {
	mu       sync.Mutex
	holds    map[string]*hold
	onExpire func(string)
}

func NewHolds(onExpire func(seat string)) *Holds {
	return &Holds{holds: map[string]*hold{}, onExpire: onExpire}
}

func (h *Holds) Hold(seat string, d time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if old, ok := h.holds[seat]; ok {
		old.timer.Stop()
	}
	hd := &hold{}
	hd.timer = time.AfterFunc(d, func() { h.expire(seat, hd) })
	h.holds[seat] = hd
}

// expire 只在這一筆還是目前的保留時才放出座位；onExpire 在鎖外面呼叫，它要做什麼都不會卡住 Holds。
func (h *Holds) expire(seat string, hd *hold) {
	h.mu.Lock()
	current := h.holds[seat] == hd
	if current {
		delete(h.holds, seat)
	}
	h.mu.Unlock()
	if current {
		h.onExpire(seat)
	}
}

func (h *Holds) Confirm(seat string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	hd, ok := h.holds[seat]
	if !ok {
		return false
	}
	hd.timer.Stop()
	delete(h.holds, seat)
	return true
}

func (h *Holds) Len() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.holds)
}

// CountByDay：Truncate(24h) 是用 UTC 的整天切的。直接用台北時間格式化日期就好。
func CountByDay(orders []time.Time) map[string]int {
	counts := map[string]int{}
	for _, t := range orders {
		counts[t.In(Taipei).Format(time.DateOnly)]++
	}
	return counts
}

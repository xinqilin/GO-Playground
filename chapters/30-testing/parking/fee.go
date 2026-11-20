// Package parking 是第 30 章用來練習寫測試的停車場：計費、車牌、柵欄機和日報表。
package parking

import (
	"errors"
	"time"
)

// 計費規則。
const (
	Grace     = 10 * time.Minute // 整段停車不到 10 分鐘免費
	Unit      = 30 * time.Minute // 每 30 分鐘一個單位，不足一個單位也算一個
	UnitPrice = 20               // 每個單位 20 元
	DailyCap  = 240              // 每個日曆日最多收 240 元
)

// ErrOutBeforeIn 表示出場時間比進場時間早。
var ErrOutBeforeIn = errors.New("出場時間比進場時間早")

// Fee 計算從 in 停到 out 的停車費。
// 整段停車不到 Grace 免費；否則在 in 所在時區的每個午夜切開，
// 每一天各自以 Unit 為單位進位計費，每天最多收 DailyCap。
func Fee(in, out time.Time) (int, error) {
	if out.Before(in) {
		return 0, ErrOutBeforeIn
	}
	if out.Sub(in) < Grace {
		return 0, nil
	}
	total := 0
	for start := in; start.Before(out); {
		end := nextMidnight(start)
		if end.After(out) {
			end = out
		}
		total += dayFee(end.Sub(start))
		start = end
	}
	return total, nil
}

// nextMidnight 回傳 t 之後的第一個午夜。
func nextMidnight(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d+1, 0, 0, 0, 0, t.Location())
}

// dayFee 是同一天裡停了 d 的費用。
func dayFee(d time.Duration) int {
	units := int((d + Unit - 1) / Unit)
	return min(units*UnitPrice, DailyCap)
}

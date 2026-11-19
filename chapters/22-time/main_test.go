package main

import (
	"context"
	"testing"
	"testing/synctest"
	"time"
)

func Example_etaDemo() {
	etaDemo()
	// Output:
	// 總共要花 40m30s
	// 預計送達 12:32:30
	// 比答應的晚了 2m30s
}

func Example_durationDemo() {
	durationDemo()
	// Output:
	// 1h30m0s 1.5 5400000000000
	// 1h15m30s 75.5 <nil>
	// 排隊要等 21m0s
	// 2m0s 1m0s
}

func Example_layoutDemo() {
	layoutDemo()
	// Output:
	// 2026-11-19 18:30:05
	// 2026-11-19T18:30:05Z
	// 11/19 18:30 11/19 6:30PM Thu Nov 19
	// 2026/11/19（四）18:30
	// YYYY-MM-DD hh:mm
}

func Example_parseDemo() {
	parseDemo()
	// Output:
	// 2026-11-19 18:30:00 +0000 UTC <nil>
	// parsing time "2026-11-19 18:30" as "2006/01/02 15:04": cannot parse "-11-19 18:30" as "/"
	// parsing time "2026-02-30": day out of range
}

func Example_zonesDemo() {
	zonesDemo()
	// Output:
	// 2026-11-19 10:30:00 UTC
	// 2026-11-19 18:30:00 CST
	// 2026-11-19 19:30:00 JST
	// 是同一個瞬間？ true
}

func Example_parseZoneDemo() {
	parseZoneDemo()
	// Output:
	// UTC 2026-11-19T18:30:00Z
	// Asia/Taipei 2026-11-19T18:30:00+08:00
	// 兩個差了 8h0m0s
}

func Example_calendarDemo() {
	calendarDemo()
	// Output:
	// 加一個月： 2026-03-03
	// 2 月 31 日： 2026-03-03
	// 2 月的最後一天： 2026-02-28
	// 離聖誕節還有 36 天
}

func Example_truncateDemo() {
	truncateDemo()
	// Output:
	// Truncate： 2026-11-18 08:00:00
	// 當地午夜： 2026-11-19 00:00:00
}

func Example_dstDemo() {
	dstDemo()
	// Output:
	// AddDate： 2026-11-01 12:00:00 EST
	// Add 24h： 2026-11-01 11:00:00 EST
	// 那一天其實有 25h0m0s
}

func Example_equalDemo() {
	equalDemo()
	// Output:
	// false true 0
	// true true -1
}

func Example_monotonicDemo() {
	monotonicDemo()
	// Output:
	// Now 帶著單調時鐘： true
	// Round(0) 之後： false
	// 經過的時間不會是負的： true
}

func Example_unixDemo() {
	unixDemo()
	// Output:
	// 1795084200 1795084200000
	// 2026-11-19T10:30:00Z true
}

func Example_afterFuncDemo() {
	afterFuncDemo()
	// Output:
	// 逾時放出： A1
	// B2 付款了，計時停下來： true
	// 再 Stop 一次： false
}

func Example_tickerDemo() {
	tickerDemo()
	// Output:
	// 請 101 號到櫃台
	// 請 102 號到櫃台
	// 請 103 號到櫃台
}

func Example_checkoutDemo() {
	checkoutDemo()
	// Output:
	// 11/11 23:59 雙11 可以用？ true
	// 11/12 00:00 雙11 可以用？ false
}

// TestMonthEnd：大月、小月、平年和閏年的 2 月、跨年的 12 月。
func TestMonthEnd(t *testing.T) {
	tests := []struct {
		y    int
		m    time.Month
		want int
	}{
		{2026, time.January, 31}, {2026, time.February, 28}, {2028, time.February, 29},
		{2026, time.April, 30}, {2026, time.December, 31},
	}
	for _, tt := range tests {
		got := monthEnd(time.Date(tt.y, tt.m, 15, 9, 30, 0, 0, taipei))
		if got.Year() != tt.y || got.Month() != tt.m || got.Day() != tt.want || got.Hour() != 9 || got.Location() != taipei {
			t.Errorf("monthEnd(%d-%02d-15 09:30) = %v；want 同一個月的 %d 日 09:30（台北）", tt.y, tt.m, got, tt.want)
		}
	}
}

// TestCheckoutAccept 把櫃台的時鐘換成固定的時間，不必等到真的過期。
func TestCheckoutAccept(t *testing.T) {
	cp := Coupon{Code: "週年慶", Expires: time.Date(2026, time.December, 1, 0, 0, 0, 0, taipei)}
	tests := []struct {
		now  time.Time
		want bool
	}{
		{time.Date(2026, time.November, 30, 23, 59, 59, 0, taipei), true},
		{time.Date(2026, time.December, 1, 0, 0, 0, 0, taipei), false},
		{time.Date(2026, time.November, 30, 16, 0, 0, 0, time.UTC), false}, // 台北已經是 12/1 00:00
	}
	for _, tt := range tests {
		c := Checkout{now: func() time.Time { return tt.now }}
		if got := c.Accept(cp); got != tt.want {
			t.Errorf("在 %v 結帳：Accept = %v，want %v", tt.now, got, tt.want)
		}
	}
}

// TestHoldSeat 在 synctest 的泡泡裡用假時間：保留 15 分鐘，不必真的等 15 分鐘。
// release 在 AfterFunc 自己的 goroutine 裡執行，所以用 channel 把結果交回來。
func TestHoldSeat(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		start := time.Now() // 泡泡裡的時間從 2000-01-01 00:00:00 UTC 開始
		released := make(chan string, 2)
		release := func(seat string) { released <- seat }
		holdSeat("A1", 15*time.Minute, release)
		paid := holdSeat("B2", 15*time.Minute, release)

		time.Sleep(14*time.Minute + 59*time.Second)
		synctest.Wait()
		if len(released) != 0 {
			t.Fatalf("還沒到 15 分鐘就放出了 %q", <-released)
		}
		paid.Stop() // B2 在最後一秒付款

		time.Sleep(time.Second)
		synctest.Wait()
		if n := len(released); n != 1 {
			t.Fatalf("15 分鐘時放出了 %d 個座位，want 1", n)
		}
		if seat := <-released; seat != "A1" || time.Since(start) != 15*time.Minute {
			t.Errorf("放出 %q，經過 %v；want A1、15m0s", seat, time.Since(start))
		}
	})
}

// TestCallNumbers 每 30 秒叫一號，叫完第 2 號就關門：第 3 號不會被叫到，也不用等第 90 秒。
func TestCallNumbers(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		start := time.Now()
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		var at []time.Duration
		callNumbers(ctx, 1, 3, 30*time.Second, func(n int) {
			at = append(at, time.Since(start))
			if n == 2 {
				cancel()
			}
		})
		if len(at) != 2 || at[0] != 30*time.Second || at[1] != time.Minute || time.Since(start) != time.Minute {
			t.Errorf("叫號的時間 %v，結束在 %v；want [30s 1m0s]，結束在 1m0s", at, time.Since(start))
		}
	})
}

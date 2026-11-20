package timetable

import (
	"maps"
	"testing"
	"testing/synctest"
	"time"
)

// tp 是台北時間的簡寫。
func tp(y int, m time.Month, d, h, mi int) time.Time {
	return time.Date(y, m, d, h, mi, 0, 0, Taipei)
}

func TestRemaining(t *testing.T) {
	closing := tp(2026, time.November, 19, 20, 30)
	tests := []struct {
		now  time.Time
		want string
	}{
		{tp(2026, time.November, 19, 18, 5), "還有 2 小時 25 分"},
		{tp(2026, time.November, 19, 19, 30), "還有 1 小時 0 分"},
		{closing.Add(-(44*time.Minute + 59*time.Second)), "還有 44 分"},
		{closing.Add(-30 * time.Second), "不到 1 分鐘"},
		{closing, "已經打烊"},
		{closing.Add(time.Hour), "已經打烊"},
		{time.Date(2026, time.November, 19, 11, 0, 0, 0, time.UTC), "還有 1 小時 30 分"}, // 台北 19:00
	}
	for _, tt := range tests {
		if got := Remaining(tt.now, closing); got != tt.want {
			t.Errorf("Remaining(%v) = %q，want %q", tt.now, got, tt.want)
		}
	}
}

func TestIsOpen(t *testing.T) {
	la, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		t    time.Time
		want bool
	}{
		{tp(2026, time.November, 19, 10, 29), false}, // 星期四
		{tp(2026, time.November, 19, 10, 30), true},
		{tp(2026, time.November, 19, 13, 59), true},
		{tp(2026, time.November, 19, 14, 0), false},
		{tp(2026, time.November, 19, 17, 0), true},
		{tp(2026, time.November, 19, 20, 29), true},
		{tp(2026, time.November, 19, 20, 30), false},
		{tp(2026, time.November, 23, 12, 0), false},                      // 星期一公休
		{time.Date(2026, time.November, 19, 4, 0, 0, 0, time.UTC), true}, // 台北 12:00
		{time.Date(2026, time.November, 22, 20, 0, 0, 0, la), false},     // 洛杉磯星期日 20:00，台北已經是星期一 12:00
		{time.Date(2026, time.November, 17, 21, 0, 0, 0, la), true},      // 洛杉磯星期二 21:00，台北是星期三 13:00
	}
	for _, tt := range tests {
		if got := IsOpen(tt.t); got != tt.want {
			t.Errorf("IsOpen(%v) = %v，want %v", tt.t.In(Taipei), got, tt.want)
		}
	}
}

func TestBillingDates(t *testing.T) {
	format := func(ts []time.Time) []string {
		var out []string
		for _, x := range ts {
			out = append(out, x.Format("2006-01-02 15:04 MST"))
		}
		return out
	}
	tests := []struct {
		first time.Time
		n     int
		want  []string
	}{
		{tp(2026, time.January, 31, 9, 0), 5, []string{
			"2026-01-31 09:00 CST", "2026-02-28 09:00 CST", "2026-03-31 09:00 CST", "2026-04-30 09:00 CST", "2026-05-31 09:00 CST"}},
		{tp(2027, time.December, 30, 23, 15), 3, []string{
			"2027-12-30 23:15 CST", "2028-01-30 23:15 CST", "2028-02-29 23:15 CST"}},
		{time.Date(2026, time.November, 15, 0, 0, 0, 0, time.UTC), 3, []string{
			"2026-11-15 00:00 UTC", "2026-12-15 00:00 UTC", "2027-01-15 00:00 UTC"}},
		{tp(2026, time.March, 1, 12, 0), 0, nil},
	}
	for _, tt := range tests {
		got := format(BillingDates(tt.first, tt.n))
		if len(got) != len(tt.want) {
			t.Errorf("BillingDates(%v, %d) = %q\nwant %q", tt.first, tt.n, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("BillingDates(%v, %d) = %q\nwant %q", tt.first, tt.n, got, tt.want)
				break
			}
		}
	}
}

func TestHolds(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		expired := make(chan string, 10)
		h := NewHolds(func(seat string) { expired <- seat })

		h.Hold("A1", 10*time.Minute)
		h.Hold("B2", 10*time.Minute)
		h.Hold("C3", 10*time.Minute)
		if h.Len() != 3 {
			t.Fatalf("保留了 3 個座位，Len = %d", h.Len())
		}

		time.Sleep(5 * time.Minute)
		h.Hold("A1", 10*time.Minute) // A1 重新計時：第 15 分鐘才逾時
		if !h.Confirm("C3") {
			t.Error("C3 在第 5 分鐘確認，應該回傳 true")
		}

		time.Sleep(5 * time.Minute)
		synctest.Wait()
		if len(expired) != 1 || h.Len() != 1 {
			t.Fatalf("第 10 分鐘：逾時了 %d 個、還保留 %d 個；want 只有 B2 逾時、還保留 A1", len(expired), h.Len())
		}
		if seat := <-expired; seat != "B2" {
			t.Errorf("第 10 分鐘逾時的是 %q，want B2", seat)
		}
		if h.Confirm("B2") || h.Confirm("C3") || h.Confirm("Z9") {
			t.Error("逾時的、已經確認過的、沒保留過的座位，Confirm 都應該回傳 false")
		}

		time.Sleep(4 * time.Minute)
		if !h.Confirm("A1") || h.Len() != 0 {
			t.Errorf("第 14 分鐘確認 A1 應該成功，之後 Len 是 0（現在 %d）", h.Len())
		}

		time.Sleep(10 * time.Minute)
		synctest.Wait()
		if len(expired) != 0 {
			t.Errorf("確認過的 A1 不應該再逾時，卻放出了 %q", <-expired)
		}
	})
}

func TestCountByDay(t *testing.T) {
	orders := []time.Time{
		tp(2026, time.November, 19, 7, 30),
		tp(2026, time.November, 19, 12, 0),
		tp(2026, time.November, 19, 23, 59),
		tp(2026, time.November, 20, 0, 1),
		time.Date(2026, time.November, 19, 20, 0, 0, 0, time.UTC), // 台北 11/20 04:00
	}
	want := map[string]int{"2026-11-19": 3, "2026-11-20": 2}
	if got := CountByDay(orders); !maps.Equal(got, want) {
		t.Errorf("CountByDay = %v，want %v", got, want)
	}
}

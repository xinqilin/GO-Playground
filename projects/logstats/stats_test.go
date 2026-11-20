package main

import (
	"fmt"
	"testing"
	"time"
)

func TestRoute(t *testing.T) {
	tests := []struct{ method, path, want string }{
		{"GET", "/menu", "GET /menu"},
		{"GET", "/menu?lang=en", "GET /menu"},
		{"GET", "/orders/17", "GET /orders/{id}"},
		{"GET", "/orders/17/items/3", "GET /orders/{id}/items/{id}"},
		{"GET", "/reports/2026/11/21", "GET /reports/{id}/{id}/{id}"},
		{"GET", "/orders/abc", "GET /orders/abc"},
		{"GET", "/", "GET /"},
		{"GET", "/v2/menu", "GET /v2/menu"},
	}
	for _, tt := range tests {
		if got := route(tt.method, tt.path); got != tt.want {
			t.Errorf("route(%q, %q) = %q，想要 %q", tt.method, tt.path, got, tt.want)
		}
	}
}

func Example_percentile() {
	ms := func(ns ...int) []time.Duration {
		var d []time.Duration
		for _, n := range ns {
			d = append(d, time.Duration(n)*time.Millisecond)
		}
		return d
	}
	orders := ms(45, 48, 52, 61, 230, 1204) // 已經排序好
	fmt.Println(percentile(orders, 50), percentile(orders, 95), percentile(orders, 100))

	many := ms(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20)
	fmt.Println(percentile(many, 50), percentile(many, 95), percentile(many, 100))
	// Output:
	// 52ms 1.204s 1.204s
	// 10ms 19ms 20ms
}

func TestPercentile(t *testing.T) {
	one := []time.Duration{7 * time.Millisecond}
	for _, p := range []int{0, 1, 50, 99, 100} {
		if got := percentile(one, p); got != one[0] {
			t.Errorf("只有一個值時，percentile(%d) = %v", p, got)
		}
	}
	hundred := make([]time.Duration, 100)
	for i := range hundred {
		hundred[i] = time.Duration(i + 1)
	}
	for _, p := range []int{1, 50, 95, 99, 100} {
		if got := percentile(hundred, p); got != time.Duration(p) {
			t.Errorf("1 到 100 的 p%d = %v，想要 %d", p, got, p)
		}
	}
}

// TestBusiestTie：兩分鐘一樣多時，取比較早的那一分鐘；跨檔案、不同的時區寫法也算同一分鐘。
func TestBusiestTie(t *testing.T) {
	taipei := time.FixedZone("", 8*60*60)
	utc := time.UTC
	s := newStats()
	for _, tm := range []time.Time{
		time.Date(2026, 11, 21, 9, 31, 10, 0, taipei),
		time.Date(2026, 11, 21, 1, 31, 50, 0, utc), // 和上一筆是同一分鐘
		time.Date(2026, 11, 21, 9, 32, 0, 0, taipei),
		time.Date(2026, 11, 21, 9, 32, 59, 0, taipei),
	} {
		s.Add(Entry{Time: tm, Method: "GET", Path: "/menu", Status: 200})
	}
	r := s.Report(10)
	if r.BusiestCount != 2 || r.BusiestMinute.Format("15:04 -07:00") != "09:31 +08:00" {
		t.Errorf("最忙的一分鐘 = %s，%d 筆", r.BusiestMinute.Format("15:04 -07:00"), r.BusiestCount)
	}
}

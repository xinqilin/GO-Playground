package main

import (
	"cmp"
	"maps"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Stats 一邊讀一邊累計。用 newStats 建立。
type Stats struct {
	Total, Bad int
	from, to   time.Time
	classes    map[string]int        // "2xx" → 筆數
	minutes    map[int64]int         // Unix 時間的第幾分鐘 → 筆數
	routes     map[string]*routeData // "GET /orders/{id}" → 這一類請求的資料
}

type routeData struct {
	count, errors int
	took          []time.Duration
}

func newStats() *Stats {
	return &Stats{classes: map[string]int{}, minutes: map[int64]int{}, routes: map[string]*routeData{}}
}

// Add 把一筆請求算進去。
func (s *Stats) Add(e Entry) {
	s.Total++
	if s.from.IsZero() || e.Time.Before(s.from) {
		s.from = e.Time
	}
	if e.Time.After(s.to) {
		s.to = e.Time
	}
	s.classes[strconv.Itoa(e.Status/100)+"xx"]++
	s.minutes[e.Time.Unix()/60]++

	key := route(e.Method, e.Path)
	rd := s.routes[key]
	if rd == nil {
		rd = &routeData{}
		s.routes[key] = rd
	}
	rd.count++
	if e.Status >= 500 {
		rd.errors++
	}
	rd.took = append(rd.took, e.Took)
}

// route 把同一類的請求歸在一起：去掉 ? 後面的參數，全是數字的路徑段換成 {id}。
func route(method, path string) string {
	path, _, _ = strings.Cut(path, "?")
	segs := strings.Split(path, "/")
	for i, seg := range segs {
		if seg != "" && strings.Trim(seg, "0123456789") == "" {
			segs[i] = "{id}"
		}
	}
	return method + " " + strings.Join(segs, "/")
}

// Report 是統計的結果。
type Report struct {
	From, To      time.Time
	Total, Bad    int
	ByClass       map[string]int
	BusiestMinute time.Time
	BusiestCount  int
	Routes        []RouteStats
}

// RouteStats 是一類請求的統計。
type RouteStats struct {
	Route         string
	Count, Errors int
	P50, P95, Max time.Duration
}

// Report 算出結果：路徑依次數由多到少，最多 top 個。
func (s *Stats) Report(top int) Report {
	r := Report{From: s.from, To: s.to, Total: s.Total, Bad: s.Bad, ByClass: s.classes}
	for _, m := range slices.Sorted(maps.Keys(s.minutes)) { // 由早到晚，同樣多時取最早的
		if s.minutes[m] > r.BusiestCount {
			r.BusiestMinute, r.BusiestCount = time.Unix(m*60, 0).In(s.to.Location()), s.minutes[m]
		}
	}
	for key, rd := range s.routes {
		slices.Sort(rd.took)
		r.Routes = append(r.Routes, RouteStats{
			Route: key, Count: rd.count, Errors: rd.errors,
			P50: percentile(rd.took, 50), P95: percentile(rd.took, 95), Max: rd.took[len(rd.took)-1],
		})
	}
	slices.SortFunc(r.Routes, func(a, b RouteStats) int {
		return cmp.Or(cmp.Compare(b.Count, a.Count), strings.Compare(a.Route, b.Route))
	})
	r.Routes = r.Routes[:min(max(top, 0), len(r.Routes))]
	return r
}

// percentile 用 nearest-rank 的方法：排序好的 n 個值裡，第 ⌈p×n/100⌉ 個。sorted 不能是空的。
func percentile(sorted []time.Duration, p int) time.Duration {
	rank := (p*len(sorted) + 99) / 100 // 整數的無條件進位
	return sorted[max(rank, 1)-1]
}

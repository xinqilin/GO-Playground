package main

import (
	"encoding/csv"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"io"
	"strconv"
	"text/tabwriter"
	"time"
)

// writers 是 -format 可以選的輸出格式。
var writers = map[string]func(io.Writer, Report) error{
	"text": writeText,
	"json": writeJSON,
	"csv":  writeCSV,
}

// writeText 給人看：摘要加上一張依欄對齊的表。
func writeText(w io.Writer, r Report) error {
	if r.Total == 0 {
		_, err := fmt.Fprintf(w, "沒有符合條件的請求（%d 行無法解析）\n", r.Bad)
		return err
	}
	const layout = "2006-01-02 15:04:05 -07:00"
	fmt.Fprintf(w, "期間：%s 到 %s\n", r.From.Format(layout), r.To.Format(layout))
	fmt.Fprintf(w, "請求：%d 筆，另有 %d 行無法解析\n", r.Total, r.Bad)
	fmt.Fprintf(w, "狀態：2xx %d、3xx %d、4xx %d、5xx %d（%.1f%%）\n",
		r.ByClass["2xx"], r.ByClass["3xx"], r.ByClass["4xx"], r.ByClass["5xx"], percent(r.ByClass["5xx"], r.Total))
	fmt.Fprintf(w, "最忙的一分鐘：%s，%d 筆\n\n", r.BusiestMinute.Format("15:04"), r.BusiestCount)

	// AlignRight：每一欄靠右對齊。最後的 route 後面沒有 \t，不參與對齊。
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintln(tw, "count\t5xx\tp50\tp95\tmax\t  route")
	for _, rt := range r.Routes {
		fmt.Fprintf(tw, "%d\t%d\t%v\t%v\t%v\t  %s\n", rt.Count, rt.Errors, rt.P50, rt.P95, rt.Max, rt.Route)
	}
	return tw.Flush()
}

func percent(n, total int) float64 { return float64(n) * 100 / float64(total) }

// ms 把時間換成毫秒的數字。
func ms(d time.Duration) float64 { return float64(d) / float64(time.Millisecond) }

// jsonRoute 是 -json 輸出裡的一列。時間用毫秒的數字：json/v2 對 time.Duration 沒有預設的格式，
// 而且讀這份 JSON 的程式不一定是 Go，數字比 "1.204s" 這種字串好處理。
type jsonRoute struct {
	Route  string  `json:"route"`
	Count  int     `json:"count"`
	Errors int     `json:"errors_5xx"`
	P50    float64 `json:"p50_ms"`
	P95    float64 `json:"p95_ms"`
	Max    float64 `json:"max_ms"`
}

// writeJSON 給其他程式讀。輸出的格式是獨立的型別，改了內部的 Report 也不會不小心改到它。
func writeJSON(w io.Writer, r Report) error {
	out := struct {
		From    time.Time      `json:"from"`
		To      time.Time      `json:"to"`
		Total   int            `json:"total"`
		Bad     int            `json:"bad_lines"`
		ByClass map[string]int `json:"by_class"`
		Routes  []jsonRoute    `json:"routes"`
	}{From: r.From, To: r.To, Total: r.Total, Bad: r.Bad, ByClass: r.ByClass}
	for _, rt := range r.Routes {
		out.Routes = append(out.Routes, jsonRoute{rt.Route, rt.Count, rt.Errors, ms(rt.P50), ms(rt.P95), ms(rt.Max)})
	}
	if err := json.MarshalWrite(w, out, json.Deterministic(true), jsontext.WithIndent("  ")); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// writeCSV 給試算表用：每個路徑一列。
func writeCSV(w io.Writer, r Report) error {
	cw := csv.NewWriter(w)
	cw.Write([]string{"route", "count", "errors_5xx", "p50_ms", "p95_ms", "max_ms"})
	f := func(d time.Duration) string { return strconv.FormatFloat(ms(d), 'f', -1, 64) }
	for _, rt := range r.Routes {
		cw.Write([]string{rt.Route, strconv.Itoa(rt.Count), strconv.Itoa(rt.Errors), f(rt.P50), f(rt.P95), f(rt.Max)})
	}
	cw.Flush()
	return cw.Error() // Write 的錯誤會留到這裡一起回報
}

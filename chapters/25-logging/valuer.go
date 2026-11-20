package main

import (
	"fmt"
	"log/slog"
	"os"
)

// Phone 是顧客的手機號碼。實作 slog.LogValuer：記錄時只留頭尾，完整的號碼不會出現在 log 裡。
type Phone string

func (p Phone) LogValue() slog.Value {
	s := string(p)
	if len(s) < 7 {
		return slog.StringValue("***")
	}
	return slog.StringValue(s[:4] + "***" + s[len(s)-3:])
}

// Customer 是外送的顧客。LogValue 只挑出可以記的欄位，名字和地址都不記。
type Customer struct {
	ID      int
	Name    string
	Phone   Phone
	Address string
}

func (c Customer) LogValue() slog.Value {
	return slog.GroupValue(slog.Int("id", c.ID), slog.Any("phone", c.Phone))
}

func valuerDemo() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, demoOptions))
	c := Customer{ID: 1024, Name: "王小明", Phone: "0912345678", Address: "台北市信義路五段 7 號"}
	logger.Info("外送訂單", "customer", c)
	logger.Info("回撥", "phone", c.Phone)
}

// report 是要花時間計算的報表。放進 LogValue，只有這筆紀錄真的要寫出去時才會計算。
type report struct{ calls *int }

func (r report) LogValue() slog.Value {
	*r.calls++
	return slog.StringValue("今日營收 12,345 元")
}

func lazyDemo() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, demoOptions)) // 沒有設定 Level，最低是 Info
	calls := 0
	logger.Debug("報表", "report", report{&calls})
	fmt.Println("計算了", calls, "次")
	logger.Info("報表", "report", report{&calls})
	fmt.Println("計算了", calls, "次")
}

package main

import (
	"log/slog"
	"os"
	"time"
)

// dropTime 拿掉每一筆紀錄的時間，範例的輸出才會每次一樣。正式的程式要保留時間。
func dropTime(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey && len(groups) == 0 {
		return slog.Attr{} // 回傳空的 Attr 就是不輸出這個欄位
	}
	return a
}

// demoOptions 是本章範例共用的 Handler 設定。
var demoOptions = &slog.HandlerOptions{ReplaceAttr: dropTime}

// textJSONDemo 同一筆紀錄，用 TextHandler 和 JSONHandler 各寫一次。
func textJSONDemo() {
	text := slog.New(slog.NewTextHandler(os.Stdout, demoOptions))
	text.Info("出餐", "table", 7, "dish", "排骨飯", "note", "少飯 不要辣", "wait", 95*time.Second)

	jsonLog := slog.New(slog.NewJSONHandler(os.Stdout, demoOptions))
	jsonLog.Info("出餐", "table", 7, "dish", "排骨飯", "note", "少飯 不要辣", "wait", 95*time.Second)
}

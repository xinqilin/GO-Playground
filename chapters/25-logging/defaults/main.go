// defaults 示範沒有任何設定時，套件層級的 log 和 slog 寫出什麼，以及 slog.SetDefault 之後的變化。
// 輸出有日期和時間，每次執行都不同，所以這個程式沒有測試；頁面上貼的是實際執行的結果。
//
// 執行：go run ./chapters/25-logging/defaults
package main

import (
	"log"
	"log/slog"
	"os"
)

func main() {
	log.Println("開火")
	slog.Info("出餐", "table", 7)

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))
	log.Println("舊的程式碼呼叫 log")
	slog.Info("出餐", "table", 8)
}

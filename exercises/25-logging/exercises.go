//go:build !solution

package journal

import (
	"context"
	"io"
	"log/slog"
	"time"
)

// NewLogger 依照設定建立 Logger，寫到 w。
//   - format 是 "text" 或 "json"，其他值回傳錯誤。
//   - level 是最低要記的層級，寫法和 slog.Level 的文字一樣：debug、INFO、warn+1 都可以，
//     不分大小寫；空字串表示 Info。看不懂時回傳錯誤。
func NewLogger(w io.Writer, format, level string) (*slog.Logger, error) {
	// TODO: 實作
	return nil, nil
}

// Redact 回傳一個 ReplaceAttr 函式：key 是 keys 其中之一（不分大小寫）的屬性，
// 不管在第幾層群組裡，值都換成字串 "[REDACTED]"；其他屬性不變。
func Redact(keys ...string) func(groups []string, a slog.Attr) slog.Attr {
	// TODO: 實作
	return nil
}

// Summarize 讀 JSONHandler 寫出的 log（每一行一筆），統計每個 level 的筆數，
// 並收集 ERROR 以上每一筆的 msg。某一筆不是合法的 JSON，或沒有 level、msg 時，
// 回傳錯誤，錯誤訊息要指出是第幾筆（從 1 開始）。
func Summarize(r io.Reader) (Summary, error) {
	// TODO: 實作
	return Summary{}, nil
}

// Timed 開始計時，回傳的函式在工作結束時呼叫（通常寫成 defer），記一筆紀錄：
//   - 平常記 Info，msg 是「完成」；花的時間到達 slow 時記 Warn，msg 是「太慢」。
//   - 兩種都帶著 op（工作名稱）和 elapsed（花了多久，time.Duration）。
func Timed(logger *slog.Logger, op string, slow time.Duration) (done func()) {
	// TODO: 實作
	return func() {}
}

// CountingHandler 包住另一個 Handler：每一筆被記下來的紀錄都照常交給它，
// 同時依照層級計數。從它 With、WithGroup 出來的 Handler 共用同一份計數。
// 可以被很多 goroutine 同時使用。
type CountingHandler struct {
	// TODO: 加上需要的欄位
}

// NewCountingHandler 建立包住 next 的 CountingHandler。
func NewCountingHandler(next slog.Handler) *CountingHandler {
	// TODO: 實作
	return &CountingHandler{}
}

func (h *CountingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// TODO: 實作
	return false
}

func (h *CountingHandler) Handle(ctx context.Context, r slog.Record) error {
	// TODO: 實作
	return nil
}

func (h *CountingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// TODO: 實作
	return h
}

func (h *CountingHandler) WithGroup(name string) slog.Handler {
	// TODO: 實作
	return h
}

// Count 回傳目前為止記了幾筆 level 層級的紀錄。被 Enabled 擋掉的不算。
func (h *CountingHandler) Count(level slog.Level) int {
	// TODO: 實作
	return 0
}

// LogItems 每個品項記一筆 Info「出餐」，帶著 item。
//
// 找錯：第一筆正常，第二筆起每一筆的 item 越來越多。
func LogItems(logger *slog.Logger, items []string) {
	for _, it := range items {
		logger = logger.With("item", it)
		logger.Info("出餐")
	}
}

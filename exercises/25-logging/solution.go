//go:build solution

package journal

import (
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"time"
)

// NewLogger 用 Level.UnmarshalText 解析層級，它本來就不分大小寫，也懂 warn+1 這種寫法。
func NewLogger(w io.Writer, format, level string) (*slog.Logger, error) {
	var lv slog.Level // 零值是 Info
	if level != "" {
		if err := lv.UnmarshalText([]byte(level)); err != nil {
			return nil, err
		}
	}
	opts := &slog.HandlerOptions{Level: lv}
	switch format {
	case "text":
		return slog.New(slog.NewTextHandler(w, opts)), nil
	case "json":
		return slog.New(slog.NewJSONHandler(w, opts)), nil
	}
	return nil, fmt.Errorf("不認得的格式 %q", format)
}

// Redact 只看 key，不用管 groups：ReplaceAttr 本來就會對每一層群組裡的屬性呼叫。
func Redact(keys ...string) func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if slices.ContainsFunc(keys, func(k string) bool { return strings.EqualFold(k, a.Key) }) {
			return slog.String(a.Key, "[REDACTED]")
		}
		return a
	}
}

// Summarize 一筆一筆解成只有 level 和 msg 的 struct，其他欄位略過。
func Summarize(r io.Reader) (Summary, error) {
	s := Summary{Counts: map[string]int{}}
	dec := jsontext.NewDecoder(r)
	for n := 1; ; n++ {
		var rec struct {
			Level string  `json:"level"`
			Msg   *string `json:"msg"`
		}
		err := json.UnmarshalDecode(dec, &rec)
		if errors.Is(err, io.EOF) {
			return s, nil
		}
		if err != nil {
			return Summary{}, fmt.Errorf("第 %d 筆：%w", n, err)
		}
		var lv slog.Level
		if err := lv.UnmarshalText([]byte(rec.Level)); err != nil || rec.Msg == nil {
			return Summary{}, fmt.Errorf("第 %d 筆：沒有 level 或 msg", n)
		}
		s.Counts[rec.Level]++
		if lv >= slog.LevelError {
			s.Errors = append(s.Errors, *rec.Msg)
		}
	}
}

// Timed 用 closure 記住開始的時間。
func Timed(logger *slog.Logger, op string, slow time.Duration) (done func()) {
	start := time.Now()
	return func() {
		elapsed := time.Since(start)
		if elapsed >= slow {
			logger.Warn("太慢", "op", op, "elapsed", elapsed)
			return
		}
		logger.Info("完成", "op", op, "elapsed", elapsed)
	}
}

// counter 是共用的計數，WithAttrs、WithGroup 出來的 Handler 都指向同一個。
type counter struct {
	mu sync.Mutex
	n  map[slog.Level]int
}

type CountingHandler struct {
	next  slog.Handler
	count *counter
}

func NewCountingHandler(next slog.Handler) *CountingHandler {
	return &CountingHandler{next: next, count: &counter{n: map[slog.Level]int{}}}
}

func (h *CountingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *CountingHandler) Handle(ctx context.Context, r slog.Record) error {
	h.count.mu.Lock()
	h.count.n[r.Level]++
	h.count.mu.Unlock()
	return h.next.Handle(ctx, r)
}

func (h *CountingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CountingHandler{next: h.next.WithAttrs(attrs), count: h.count}
}

func (h *CountingHandler) WithGroup(name string) slog.Handler {
	return &CountingHandler{next: h.next.WithGroup(name), count: h.count}
}

func (h *CountingHandler) Count(level slog.Level) int {
	h.count.mu.Lock()
	defer h.count.mu.Unlock()
	return h.count.n[level]
}

// LogItems 每一圈的 item 只屬於那一筆，直接放在 Info 的參數裡；
// 要用 With 的話，結果存到新的變數，不要蓋掉 logger。
func LogItems(logger *slog.Logger, items []string) {
	for _, it := range items {
		logger.Info("出餐", "item", it)
	}
}

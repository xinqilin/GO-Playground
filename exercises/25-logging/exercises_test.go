package journal

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"io"
	"log/slog"
	"reflect"
	"strings"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

// records 把 JSONHandler 的輸出一行一行解成 map。
func records(t *testing.T, r io.Reader) []map[string]any {
	t.Helper()
	var out []map[string]any
	dec := jsontext.NewDecoder(r)
	for {
		var m map[string]any
		err := json.UnmarshalDecode(dec, &m)
		if errors.Is(err, io.EOF) {
			return out
		}
		if err != nil {
			t.Fatalf("輸出不是合法的 JSON：%v", err)
		}
		out = append(out, m)
	}
}

// noTime 拿掉時間，TextHandler 的輸出才能直接比對。
func noTime(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey && len(groups) == 0 {
		return slog.Attr{}
	}
	return a
}

func TestNewLogger(t *testing.T) {
	var buf bytes.Buffer
	logger, err := NewLogger(&buf, "json", "DEBUG")
	if err != nil || logger == nil {
		t.Fatalf("NewLogger(json, DEBUG) = %v, %v", logger, err)
	}
	logger.Debug("備料", "step", 1)
	recs := records(t, &buf)
	if len(recs) != 1 || recs[0]["level"] != "DEBUG" || recs[0]["msg"] != "備料" || recs[0]["step"] != 1.0 {
		t.Errorf("json、DEBUG：%v", recs)
	}

	buf.Reset()
	logger, err = NewLogger(&buf, "text", "")
	if err != nil || logger == nil {
		t.Fatalf("NewLogger(text, \"\") = %v, %v", logger, err)
	}
	logger.Debug("看不到")
	logger.Info("接單", "table", 7)
	if got := buf.String(); strings.Contains(got, "看不到") || !strings.Contains(got, "level=INFO msg=接單 table=7") {
		t.Errorf("text、預設 Info：%q", got)
	}

	buf.Reset()
	logger, err = NewLogger(&buf, "text", "warn+1")
	if err != nil || logger == nil {
		t.Fatalf("NewLogger(text, warn+1) = %v, %v", logger, err)
	}
	logger.Warn("看不到")
	logger.Error("看得到")
	if got := buf.String(); strings.Contains(got, "看不到") || !strings.Contains(got, "看得到") {
		t.Errorf("warn+1：%q", got)
	}

	for _, c := range [][2]string{{"xml", "info"}, {"json", "verbose"}, {"", "info"}} {
		if _, err := NewLogger(io.Discard, c[0], c[1]); err == nil {
			t.Errorf("NewLogger(%q, %q) 應該回傳錯誤", c[0], c[1])
		}
	}
}

func TestRedact(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{ReplaceAttr: Redact("password", "token")}))
	logger.Info("登入", "user", "amei", "Password", "hunter2",
		slog.Group("auth", "token", "abc123", "scheme", "Bearer", slog.Group("refresh", "TOKEN", "xyz")))
	recs := records(t, &buf)
	if len(recs) != 1 {
		t.Fatalf("記了 %d 筆，want 1", len(recs))
	}
	want := map[string]any{
		"level": "INFO", "msg": "登入", "user": "amei", "Password": "[REDACTED]",
		"auth": map[string]any{"token": "[REDACTED]", "scheme": "Bearer", "refresh": map[string]any{"TOKEN": "[REDACTED]"}},
	}
	got := recs[0]
	if _, ok := got["time"].(string); !ok {
		t.Errorf("time 不該被拿掉或遮起來：%v", got["time"])
	}
	delete(got, "time")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Redact：\n got %v\nwant %v", got, want)
	}
}

func TestSummarize(t *testing.T) {
	log := `{"time":"2026-11-20T10:00:00+08:00","level":"INFO","msg":"開店"}
{"time":"2026-11-20T10:05:00+08:00","level":"WARN","msg":"排骨快賣完了","left":3}

{"time":"2026-11-20T10:06:00+08:00","level":"ERROR","msg":"收銀機斷線","err":"timeout"}
{"time":"2026-11-20T10:07:00+08:00","level":"INFO","msg":"接單","table":7}
{"time":"2026-11-20T10:08:00+08:00","level":"ERROR+2","msg":"冰箱溫度過高"}
{"time":"2026-11-20T10:09:00+08:00","level":"WARN+3","msg":"快要出事"}
`
	got, err := Summarize(strings.NewReader(log))
	want := Summary{
		Counts: map[string]int{"INFO": 2, "WARN": 1, "ERROR": 1, "ERROR+2": 1, "WARN+3": 1},
		Errors: []string{"收銀機斷線", "冰箱溫度過高"},
	}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("Summarize = %+v（%v），want %+v", got, err, want)
	}

	bad := `{"level":"INFO","msg":"a"}
{"level":"INFO","msg":"b"}
{"level":"INFO","msg":"c"
`
	if _, err := Summarize(strings.NewReader(bad)); err == nil || !strings.Contains(err.Error(), "第 3 筆") {
		t.Errorf("第 3 筆壞掉：err = %v，want 提到「第 3 筆」", err)
	}
	if _, err := Summarize(strings.NewReader(`{"level":"INFO"}` + "\n")); err == nil || !strings.Contains(err.Error(), "第 1 筆") {
		t.Errorf("沒有 msg：err = %v，want 提到「第 1 筆」", err)
	}
	if got, err := Summarize(strings.NewReader("")); err != nil || len(got.Counts) != 0 || got.Errors != nil {
		t.Errorf("空的 log：%+v（%v）", got, err)
	}
}

func TestTimed(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&buf, nil))

		done := Timed(logger, "炸排骨", 3*time.Second)
		time.Sleep(2 * time.Second)
		done()

		done = Timed(logger, "滷肉", 3*time.Second)
		time.Sleep(3 * time.Second)
		done()

		recs := records(t, &buf)
		if len(recs) != 2 {
			t.Fatalf("記了 %d 筆，want 2：%v", len(recs), recs)
		}
		wants := []map[string]any{
			{"level": "INFO", "msg": "完成", "op": "炸排骨", "elapsed": float64(2 * time.Second)},
			{"level": "WARN", "msg": "太慢", "op": "滷肉", "elapsed": float64(3 * time.Second)},
		}
		for i, want := range wants {
			for k, w := range want {
				if recs[i][k] != w {
					t.Errorf("第 %d 筆的 %s = %v，want %v", i+1, k, recs[i][k], w)
				}
			}
		}
	})
}

func TestCountingHandler(t *testing.T) {
	var buf bytes.Buffer
	h := NewCountingHandler(slog.NewJSONHandler(&buf, nil))
	logger := slog.New(h)

	logger.Debug("看不到")
	logger.Info("接單")
	logger.With("table", 7).Warn("等太久")
	logger.WithGroup("pay").Error("刷卡失敗", "code", 51)

	var wg sync.WaitGroup
	for range 50 {
		wg.Go(func() { logger.With("worker", 1).Info("出餐") })
	}
	wg.Wait()

	for level, want := range map[slog.Level]int{slog.LevelDebug: 0, slog.LevelInfo: 51, slog.LevelWarn: 1, slog.LevelError: 1} {
		if got := h.Count(level); got != want {
			t.Errorf("Count(%v) = %d，want %d", level, got, want)
		}
	}

	recs := records(t, &buf)
	if len(recs) != 53 {
		t.Fatalf("寫出 %d 筆，want 53", len(recs))
	}
	if recs[1]["table"] != 7.0 {
		t.Errorf("With 的屬性不見了：%v", recs[1])
	}
	if pay, ok := recs[2]["pay"].(map[string]any); !ok || pay["code"] != 51.0 {
		t.Errorf("WithGroup 的群組不見了：%v", recs[2])
	}
}

func TestLogItems(t *testing.T) {
	var buf bytes.Buffer
	LogItems(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{ReplaceAttr: noTime})), []string{"排骨飯", "滷肉飯", "雞腿飯"})
	want := "level=INFO msg=出餐 item=排骨飯\n" +
		"level=INFO msg=出餐 item=滷肉飯\n" +
		"level=INFO msg=出餐 item=雞腿飯\n"
	if got := buf.String(); got != want {
		t.Errorf("LogItems 寫出\n%s\nwant\n%s", got, want)
	}
}

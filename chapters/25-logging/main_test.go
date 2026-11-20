package main

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"io"
	"log/slog"
	"testing"
	"testing/slogtest"
)

func Example_logDemo() {
	logDemo()
	// Output:
	// [廚房] 開火
	// [廚房] 第 7 桌 排骨飯 出餐
	// [外場] 叫號 42
}

func Example_textJSONDemo() {
	textJSONDemo()
	// Output:
	// level=INFO msg=出餐 table=7 dish=排骨飯 note="少飯 不要辣" wait=1m35s
	// {"level":"INFO","msg":"出餐","table":7,"dish":"排骨飯","note":"少飯 不要辣","wait":95000000000}
}

func Example_attrsDemo() {
	attrsDemo()
	// Output:
	// {"level":"WARN","msg":"出餐延遲","table":7,"order":{"dish":"排骨飯","qty":2},"err":"排骨還在炸"}
	// {"level":"INFO","msg":"出餐","table":8,"dish":"雞腿飯"}
}

func Example_withDemo() {
	withDemo()
	// Output:
	// level=INFO msg=接單 order_id=A042 table=7
	// level=INFO msg=出餐 order_id=A042 table=7 dish=排骨飯
	// level=INFO msg=付款 order_id=A042 table=7 pay.method=悠遊卡 pay.amount=225
	// level=INFO msg=打烊
}

func Example_levelsDemo() {
	levelsDemo()
	// Output:
	// DEBUG INFO WARN ERROR
	// -4 0 4 8
	// INFO+2 DEBUG-4
	// warn → WARN
	// INFO+2 → INFO+2
	// verbose → slog: level string "verbose": unknown name
}

func Example_levelVarDemo() {
	levelVarDemo()
	// Output:
	// level=INFO msg=接單 table=7
	// level=DEBUG msg=看得到了 step=2
	// level=WARN msg=排骨快賣完了 left=3
	// Info 要不要記： false
}

func Example_valuerDemo() {
	valuerDemo()
	// Output:
	// {"level":"INFO","msg":"外送訂單","customer":{"id":1024,"phone":"0912***678"}}
	// {"level":"INFO","msg":"回撥","phone":"0912***678"}
}

func Example_lazyDemo() {
	lazyDemo()
	// Output:
	// 計算了 0 次
	// level=INFO msg=報表 report="今日營收 12,345 元"
	// 計算了 1 次
}

func Example_multiDemo() {
	multiDemo()
	// Output:
	// level=INFO msg=接單 table=7
	// level=WARN msg=排骨快賣完了 left=3
	// level=ERROR msg=收銀機斷線 err="connection refused"
	// 告警：
	// {"level":"WARN","msg":"排骨快賣完了","left":3}
	// {"level":"ERROR","msg":"收銀機斷線","err":"connection refused"}
}

func Example_contextDemo() {
	contextDemo()
	// Output:
	// main.ctxHandler：
	// level=INFO msg=接單 table=7 request_id=r-7f3a
	// level=INFO msg=出餐 table=7 request_id=r-7f3a
	// main.handleOnly：
	// level=INFO msg=接單 table=7 request_id=r-7f3a
	// level=INFO msg=出餐 table=7
}

func Example_checkoutDemo() {
	checkoutDemo()
	// Output:
	// level=ERROR msg=付款失敗 order_id=A042 total=225 err=餘額不足
	// 訂單 A042：餘額不足
}

// logLines 把 JSONHandler 寫出的每一行解成一個 map，交給測試檢查。
func logLines(t *testing.T, r io.Reader) []map[string]any {
	t.Helper()
	var lines []map[string]any
	dec := jsontext.NewDecoder(r)
	for {
		var m map[string]any
		err := json.UnmarshalDecode(dec, &m)
		if errors.Is(err, io.EOF) {
			return lines
		}
		if err != nil {
			t.Fatalf("log 不是合法的 JSON：%v", err)
		}
		lines = append(lines, m)
	}
}

// TestCheckoutLogs 付款失敗時，要記一筆 ERROR，帶著訂單編號、金額和錯誤。
func TestCheckoutLogs(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	order := Order{ID: "A042", Items: 3, Total: 225}

	err := checkout(t.Context(), logger, order, func(int) error { return errors.New("餘額不足") })
	if err == nil {
		t.Fatal("付款失敗時 checkout 應該回傳錯誤")
	}
	lines := logLines(t, &buf)
	if len(lines) != 2 {
		t.Fatalf("記了 %d 筆，want 2：%v", len(lines), lines)
	}
	want := map[string]any{"level": "ERROR", "msg": "付款失敗", "order_id": "A042", "total": 225.0, "err": "餘額不足"}
	for key, w := range want {
		if got := lines[1][key]; got != w {
			t.Errorf("%s = %v，want %v", key, got, w)
		}
	}
}

// TestCtxHandler 用 slogtest 檢查 ctxHandler 符合 Handler 的規則，再確認 With 之後還有 request id。
func TestCtxHandler(t *testing.T) {
	var buf bytes.Buffer
	h := ctxHandler{slog.NewJSONHandler(&buf, nil)}
	if err := slogtest.TestHandler(h, func() []map[string]any { return logLines(t, &buf) }); err != nil {
		t.Error(err)
	}

	buf.Reset()
	slog.New(h).With("table", 7).InfoContext(withRequestID(t.Context(), "r-1"), "出餐")
	lines := logLines(t, &buf)
	if len(lines) != 1 || lines[0]["request_id"] != "r-1" || lines[0]["table"] != 7.0 {
		t.Errorf("With 之後：%v，want request_id=r-1 table=7", lines)
	}
}

// TestPhoneLogValue 太短的號碼整個遮起來。
func TestPhoneLogValue(t *testing.T) {
	for p, want := range map[Phone]string{"0912345678": "0912***678", "123": "***", "": "***"} {
		if got := p.LogValue().String(); got != want {
			t.Errorf("Phone(%q).LogValue() = %q，want %q", p, got, want)
		}
	}
}

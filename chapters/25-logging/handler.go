package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// multiDemo MultiHandler 把每一筆紀錄交給每一個 Handler，各自決定要不要記、怎麼記。
func multiDemo() {
	var alerts strings.Builder // 實際的程式裡可能是檔案，或送到告警系統的 Writer
	logger := slog.New(slog.NewMultiHandler(
		slog.NewTextHandler(os.Stdout, demoOptions),
		slog.NewJSONHandler(&alerts, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: dropTime}),
	))
	logger.Info("接單", "table", 7)
	logger.Warn("排骨快賣完了", "left", 3)
	logger.Error("收銀機斷線", "err", errors.New("connection refused"))
	fmt.Print("告警：\n", alerts.String())
}

type requestIDKey struct{}

// withRequestID 把 request id 放進 context，之後所有 ...Context 的紀錄都會帶著它。
func withRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

// ctxHandler 包住另一個 Handler，把 context 裡的 request id 加到每一筆紀錄。
// 嵌入 slog.Handler，Enabled 就直接用裡面那個的。
type ctxHandler struct {
	slog.Handler
}

func (h ctxHandler) Handle(ctx context.Context, r slog.Record) error {
	if id, ok := ctx.Value(requestIDKey{}).(string); ok {
		r.AddAttrs(slog.String("request_id", id))
	}
	return h.Handler.Handle(ctx, r)
}

// WithAttrs 和 WithGroup 要把結果再包成 ctxHandler，logger.With 之後才不會換回裡面那個 Handler。
func (h ctxHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return ctxHandler{h.Handler.WithAttrs(attrs)}
}

func (h ctxHandler) WithGroup(name string) slog.Handler {
	return ctxHandler{h.Handler.WithGroup(name)}
}

// handleOnly 只覆寫 Handle。WithAttrs 是嵌入的 Handler 提供的，回傳的不再是 handleOnly。
type handleOnly struct {
	slog.Handler
}

func (h handleOnly) Handle(ctx context.Context, r slog.Record) error {
	return ctxHandler{h.Handler}.Handle(ctx, r)
}

// contextDemo 兩個 Handler 都直接用時一樣，差別出在 logger.With 之後。
func contextDemo() {
	ctx := withRequestID(context.Background(), "r-7f3a")
	inner := slog.NewTextHandler(os.Stdout, demoOptions)
	for _, h := range []slog.Handler{ctxHandler{inner}, handleOnly{inner}} {
		fmt.Printf("%T：\n", h)
		logger := slog.New(h)
		logger.InfoContext(ctx, "接單", "table", 7)
		logger.With("table", 7).InfoContext(ctx, "出餐")
	}
}

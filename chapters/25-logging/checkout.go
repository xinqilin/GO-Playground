package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
)

// Order 是一筆外送訂單。
type Order struct {
	ID    string
	Items int
	Total int
}

// checkout 付款並記錄過程。logger 由呼叫的一方傳進來，不用全域的 slog.Default()，
// 測試時就能換成寫到 buffer 的 Logger，檢查記了什麼。
func checkout(ctx context.Context, logger *slog.Logger, o Order, pay func(amount int) error) error {
	logger = logger.With("order_id", o.ID)
	logger.DebugContext(ctx, "開始結帳", "items", o.Items)
	if err := pay(o.Total); err != nil {
		logger.ErrorContext(ctx, "付款失敗", "total", o.Total, "err", err)
		return fmt.Errorf("訂單 %s：%w", o.ID, err)
	}
	logger.InfoContext(ctx, "結帳完成", "total", o.Total)
	return nil
}

func checkoutDemo() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, demoOptions))
	err := checkout(context.Background(), logger, Order{ID: "A042", Items: 3, Total: 225}, func(int) error {
		return errors.New("餘額不足")
	})
	fmt.Println(err)
}

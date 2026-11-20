package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
)

// attrsDemo 鍵值對可以直接寫，也可以用 slog.Int 這類函式建立 Attr：
// 型別寫在函式名稱上，key 和 value 不會錯位。Group 把幾個屬性收成一組。
func attrsDemo() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, demoOptions))
	logger.Warn("出餐延遲",
		slog.Int("table", 7),
		slog.Group("order", slog.String("dish", "排骨飯"), slog.Int("qty", 2)),
		slog.Any("err", errors.New("排骨還在炸")),
	)

	// LogAttrs 只收 Attr，不收鍵值對；在很常執行的路徑上，它的配置比較少。
	logger.LogAttrs(context.Background(), slog.LevelInfo, "出餐",
		slog.Int("table", 8), slog.String("dish", "雞腿飯"))
}

// withDemo With 回傳一個新的 Logger，之後每一筆都帶著這些屬性；
// WithGroup 讓之後加上的屬性都放進一個群組。原本的 Logger 不會被改變。
func withDemo() {
	base := slog.New(slog.NewTextHandler(os.Stdout, demoOptions))
	order := base.With("order_id", "A042", "table", 7)
	order.Info("接單")
	order.Info("出餐", "dish", "排骨飯")

	pay := order.WithGroup("pay")
	pay.Info("付款", "method", "悠遊卡", "amount", 225)

	base.Info("打烊")
}

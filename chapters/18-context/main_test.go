package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

func Example() {
	main()
	// Output:
	// == 逾時 ==
	// context deadline exceeded true
	// == 手動取消 ==
	// context canceled
	// == 取消會往下傳 ==
	// child: context canceled
	// == 取消原因 ==
	// context canceled / 使用者按了取消
	// == 帶值 ==
	// 處理請求 req-42
}

func TestSlowWork(t *testing.T) {
	if err := slowWork(context.Background(), time.Millisecond); err != nil {
		t.Errorf("沒有取消時應該正常完成，got %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := slowWork(ctx, time.Hour); !errors.Is(err, context.Canceled) {
		t.Errorf("取消後應該回傳 context.Canceled，got %v", err)
	}
}

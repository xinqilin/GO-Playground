package main

import (
	"errors"
	"testing"
)

func Example() {
	main()
	// Output:
	// == 錯誤是值 ==
	// 蘋果 2 → <nil>
	// 香蕉 1 → 下單失敗：查庫存 "香蕉"：找不到
	// 蘋果 0 → 欄位 qty 的值 0 不合法
	// 蘋果 5 → 下單失敗：蘋果 只剩 3 個
	// == errors.Is：比對特定錯誤 ==
	// true
	// == errors.As 與 errors.AsType：取出錯誤型別 ==
	// As: qty -1
	// AsType: qty -1
	// == panic 與 recover ==
	// 5 <nil>
	// 0 recover: runtime error: integer divide by zero
}

func TestOrder(t *testing.T) {
	tests := []struct {
		name         string
		item         string
		qty          int
		wantNil      bool
		wantNotFound bool
		wantInvalid  bool
	}{
		{name: "成功", item: "蘋果", qty: 3, wantNil: true},
		{name: "沒有這個商品", item: "香蕉", qty: 1, wantNotFound: true},
		{name: "數量不合法", item: "蘋果", qty: 0, wantInvalid: true},
		{name: "庫存不足", item: "蘋果", qty: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := order(tt.item, tt.qty)
			if (err == nil) != tt.wantNil {
				t.Fatalf("err = %v, wantNil %v", err, tt.wantNil)
			}
			if got := errors.Is(err, ErrNotFound); got != tt.wantNotFound {
				t.Errorf("errors.Is(err, ErrNotFound) = %v, want %v", got, tt.wantNotFound)
			}
			if _, got := errors.AsType[*ValidationError](err); got != tt.wantInvalid {
				t.Errorf("AsType ValidationError = %v, want %v", got, tt.wantInvalid)
			}
		})
	}
}

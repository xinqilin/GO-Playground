package main

import (
	"errors"
	"fmt"
	"testing"
)

func Example_basicsDemo() {
	basicsDemo()
	// Output:
	// 人數 4
	// 錯誤：strconv.Atoi: parsing "四": invalid syntax
	// 錯誤：人數 20 不在 1 到 12 之間
}

func Example_sentinelDemo() {
	sentinelDemo()
	// Output:
	// 排骨飯 85
	// 雞腿飯 今天賣完了，明天請早
	// 珍奶 錯誤：菜單上沒有
}

func Example_wrapDemo() {
	wrapDemo()
	// Output:
	// 訂單 17：報價 雞腿飯：賣完了
	// false true
	// *fmt.wrapError：訂單 17：報價 雞腿飯：賣完了
	// *fmt.wrapError：報價 雞腿飯：賣完了
	// *errors.errorString：賣完了
}

func Example_wrapVsFormat() {
	wrapVsFormat()
	// Output:
	// 報價：賣完了 true
	// 報價：賣完了 false
}

func Example_asDemo() {
	asDemo()
	// Output:
	// 訂單 18：排骨飯 只剩 3 份
	// As 排骨飯 3
	// AsType 排骨飯 3
}

func Example_joinDemo() {
	joinDemo()
	// Output:
	// 要留名字
	// 人數 20 不在 1 到 12 之間
	// 電話要 10 碼
	// true
	// <nil>
}

func Example_closeDemo() {
	closeDemo()
	// Output:
	// <nil>
	// 磁碟滿了，收據沒有存進去
}

func Example_mustDemo() {
	mustDemo()
	// Output:
	// 110
	// panic: mustPrice("珍奶")：菜單上沒有
}

func Example_layersDemo() {
	layersDemo()
	// Output:
	// 200 <nil>
	// 404 訂單 19：菜單上沒有
	// 409 訂單 20：排骨飯 只剩 3 份
	// 500 資料庫連不上
}

func TestParsePeople(t *testing.T) {
	tests := []struct {
		in      string
		want    int
		wantErr bool
	}{
		{"1", 1, false},
		{"12", 12, false},
		{"0", 0, true},
		{"13", 0, true},
		{"四", 0, true},
		{"", 0, true},
	}
	for _, tt := range tests {
		got, err := parsePeople(tt.in)
		if got != tt.want || (err != nil) != tt.wantErr {
			t.Errorf("parsePeople(%q) = %d, %v, want %d, wantErr %v", tt.in, got, err, tt.want, tt.wantErr)
		}
	}
}

func TestReserveNil(t *testing.T) {
	if err := reserve("排骨飯", 3); err != nil {
		t.Errorf("reserve(排骨飯, 3) = %v, want nil", err)
	}
}

func TestStatusFor(t *testing.T) {
	tests := []struct {
		err  error
		want int
	}{
		{nil, 200},
		{ErrNotOnMenu, 404},
		{fmt.Errorf("a：%w", fmt.Errorf("b：%w", ErrSoldOut)), 409},
		{&OrderError{No: 1, Err: &StockError{Item: "x"}}, 409},
		{errors.New("其他"), 500},
	}
	for _, tt := range tests {
		if got := statusFor(tt.err); got != tt.want {
			t.Errorf("statusFor(%v) = %d, want %d", tt.err, got, tt.want)
		}
	}
}

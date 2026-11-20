package main

import (
	"fmt"
	"testing"
)

func Example_callBasics() {
	callBasics()
	// Output:
	// 255
	// 205
}

func Example_passByValue() {
	passByValue()
	// Output:
	// 函式裡： 11
	// 呼叫端： 10
}

func Example_multiReturn() {
	multiReturn()
	// Output:
	// 3 2
	// 1
}

func Example_valueAndError() {
	valueAndError()
	// Output:
	// 價格： 85
	// 錯誤： "八十五 元" 的金額不是整數
	// 錯誤： "85" 沒有「 元」結尾
}

func Example_namedResults() {
	namedResults()
	// Output:
	// 333 1
	// 0 0
}

func Example_variadic() {
	variadic()
	// Output:
	// 0 85 190
	// 190
	// 排骨、雞腿、滷蛋
}

func Example_funcValues() {
	funcValues()
	// Output:
	// 42
	// 81
	// func(int) int
	// true
}

func Example_funcTypes() {
	funcTypes()
	// Output:
	// 400
	// 405
}

func Example_sortWithFunc() {
	sortWithFunc()
	// Output:
	// [{素食 70} {排骨 85} {雞腿 90}]
}

func Example_closureCounter() {
	closureCounter()
	// Output:
	// 1 2 3
	// 1 4
}

func Example_closureShared() {
	closureShared()
	// Output:
	// 目前價格： 85
	// 目前價格： 90
}

func Example_closureFactory() {
	closureFactory()
	// Output:
	// 歡迎光臨，小明
	// 請取餐：A07
}

func Example_recursion() {
	recursion()
	// Output:
	// 1101 11111111 0
}

func Example_deferOrder() {
	deferOrder()
	// Output:
	// 營業中
	// 收銀機結帳
	// 關燈
	// 關店
}

func Example_deferArgs() {
	deferArgs()
	// Output:
	// closure： 20
	// 參數： 10
}

func Example_deferResult() {
	fmt.Println(deferResult())
	// Output:
	// 95
}

func Example_deferInLoop() {
	deferInLoop()
	// Output:
	// 開啟 a.txt
	// 處理 a.txt
	// 開啟 b.txt
	// 處理 b.txt
	// 迴圈結束
	// 關閉 b.txt
	// 關閉 a.txt
}

func Example_deferInLoopFixed() {
	deferInLoopFixed()
	// Output:
	// 開啟 a.txt
	// 處理 a.txt
	// 關閉 a.txt
	// 開啟 b.txt
	// 處理 b.txt
	// 關閉 b.txt
	// 迴圈結束
}

func Example_panicRecover() {
	panicRecover()
	// Output:
	// 3 <nil>
	// 0 拿貨失敗：庫存不足：剩 3 個，要拿 10 個
	// 程式繼續執行
}

func TestParsePrice(t *testing.T) {
	tests := []struct {
		in      string
		want    int
		wantErr bool
	}{
		{"85 元", 85, false},
		{"0 元", 0, false},
		{"-5 元", -5, false},
		{"85", 0, true},
		{"八十五 元", 0, true},
		{"85元", 0, true},
	}
	for _, tt := range tests {
		got, err := parsePrice(tt.in)
		if got != tt.want || (err != nil) != tt.wantErr {
			t.Errorf("parsePrice(%q) = %d, %v; want %d, err=%v", tt.in, got, err, tt.want, tt.wantErr)
		}
	}
}

func TestBinary(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{0, "0"}, {1, "1"}, {2, "10"}, {5, "101"}, {13, "1101"}, {1024, "10000000000"},
	}
	for _, tt := range tests {
		if got := binary(tt.n); got != tt.want {
			t.Errorf("binary(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

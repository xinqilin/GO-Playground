package inout

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestPercent(t *testing.T) {
	tests := []struct {
		part, total int
		want        string
	}{
		{1, 3, "33.3%"},
		{2, 3, "66.7%"},
		{3, 3, "100.0%"},
		{1, 8, "12.5%"},
		{0, 5, "0.0%"},
		{1, 0, "0.0%"},
		{5, -1, "0.0%"},
	}
	for _, tt := range tests {
		if got := Percent(tt.part, tt.total); got != tt.want {
			t.Errorf("Percent(%d, %d) = %q, want %q", tt.part, tt.total, got, tt.want)
		}
	}
}

func TestHexColor(t *testing.T) {
	tests := []struct {
		r, g, b uint8
		want    string
	}{
		{255, 136, 0, "#FF8800"},
		{0, 0, 0, "#000000"},
		{18, 52, 86, "#123456"},
		{255, 255, 255, "#FFFFFF"},
		{10, 11, 12, "#0A0B0C"},
	}
	for _, tt := range tests {
		if got := HexColor(tt.r, tt.g, tt.b); got != tt.want {
			t.Errorf("HexColor(%d, %d, %d) = %q, want %q", tt.r, tt.g, tt.b, got, tt.want)
		}
	}
}

func TestMoneyString(t *testing.T) {
	tests := []struct {
		m    Money
		want string
	}{
		{0, "NT$0"},
		{85, "NT$85"},
		{999, "NT$999"},
		{1000, "NT$1,000"},
		{1234, "NT$1,234"},
		{123456, "NT$123,456"},
		{1234567, "NT$1,234,567"},
		{-1500, "-NT$1,500"},
		{-85, "-NT$85"},
	}
	for _, tt := range tests {
		// 用 fmt.Sprint 而不是直接呼叫 String，確認 fmt 真的會用到它。
		if got := fmt.Sprint(tt.m); got != tt.want {
			t.Errorf("fmt.Sprint(Money(%d)) = %q, want %q", int(tt.m), got, tt.want)
		}
	}
	if got, want := fmt.Sprintf("合計 %v，找零 %s", Money(1200), Money(15)), "合計 NT$1,200，找零 NT$15"; got != want {
		t.Errorf("Sprintf = %q, want %q", got, want)
	}
}

// failAfter 寫入 n 次之後就開始失敗。
type failAfter struct {
	n   int
	buf strings.Builder
}

var errDiskFull = errors.New("磁碟滿了")

func (f *failAfter) Write(p []byte) (int, error) {
	if f.n == 0 {
		return 0, errDiskFull
	}
	f.n--
	return f.buf.Write(p)
}

func TestWriteReceipt(t *testing.T) {
	items := []Item{{"排骨飯", 2, 85}, {"滷蛋", 3, 15}, {"紅茶", 1, 25}}
	var sb strings.Builder
	if err := WriteReceipt(&sb, items); err != nil {
		t.Fatalf("WriteReceipt 回傳錯誤：%v", err)
	}
	want := "  170  排骨飯 x2\n" +
		"   45  滷蛋 x3\n" +
		"   25  紅茶 x1\n" +
		"-----\n" +
		"  240  合計\n"
	if got := sb.String(); got != want {
		t.Errorf("WriteReceipt 寫出\n%s\nwant\n%s", got, want)
	}

	var empty strings.Builder
	if err := WriteReceipt(&empty, nil); err != nil {
		t.Fatalf("沒有品項時回傳錯誤：%v", err)
	}
	if got, want := empty.String(), "-----\n    0  合計\n"; got != want {
		t.Errorf("沒有品項時寫出 %q, want %q", got, want)
	}

	for n := range 5 {
		w := &failAfter{n: n}
		if err := WriteReceipt(w, items); !errors.Is(err, errDiskFull) {
			t.Errorf("第 %d 次寫入失敗時回傳 %v, want %v", n+1, err, errDiskFull)
		}
	}
}

type errReader struct{}

var errBroken = errors.New("讀取中斷")

func (errReader) Read(p []byte) (int, error) { return 0, errBroken }

func TestParseOrders(t *testing.T) {
	input := "# 午餐訂單\n" +
		"排骨飯 2 85\n" +
		"\n" +
		"  滷蛋   3  15  \n" +
		"紅茶 1 25\n"
	got, err := ParseOrders(strings.NewReader(input))
	if err != nil || got != 240 {
		t.Errorf("ParseOrders(正常訂單) = %d, %v; want 240, nil", got, err)
	}

	bad := []struct {
		in, wantErr string
	}{
		{"排骨飯 2 85\n雞腿飯 兩 90\n", `第 2 行格式錯誤："雞腿飯 兩 90"`},
		{"# 註解\n\n排骨飯 2\n", `第 3 行格式錯誤："排骨飯 2"`},
		{"排骨飯 2 85 加辣\n", `第 1 行格式錯誤："排骨飯 2 85 加辣"`},
		{"排骨飯 2 8.5\n", `第 1 行格式錯誤："排骨飯 2 8.5"`},
	}
	for _, tt := range bad {
		got, err := ParseOrders(strings.NewReader(tt.in))
		if got != 0 || err == nil || err.Error() != tt.wantErr {
			t.Errorf("ParseOrders(%q) = %d, %v; want 0, %s", tt.in, got, err, tt.wantErr)
		}
	}

	if got, err := ParseOrders(errReader{}); got != 0 || !errors.Is(err, errBroken) {
		t.Errorf("讀取失敗時回傳 %d, %v; want 0, %v", got, err, errBroken)
	}
}

func TestCallNumbers(t *testing.T) {
	var sb strings.Builder
	if err := CallNumbers(&sb, []string{"A01", "A02"}); err != nil {
		t.Fatal(err)
	}
	want := "請 A01 到櫃台\n請 A02 到櫃台\n== 以上 ==\n"
	if got := sb.String(); got != want {
		t.Errorf("CallNumbers 寫出 %q, want %q", got, want)
	}
}

func TestLongestLine(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"一般", "abc\n\nhello\nhi\n", 5},
		{"中文以 byte 計", "排骨飯\nGo\n", 9},
		{"最後一行沒有換行", "a\nabcdef", 6},
		{"空的", "", 0},
		{"超過 64 KB", "short\n" + strings.Repeat("x", 200_000) + "\nend\n", 200_000},
		{"剛好 1 MB", strings.Repeat("y", 1<<20) + "\n", 1 << 20},
	}
	for _, tt := range tests {
		got, err := LongestLine(strings.NewReader(tt.in))
		if err != nil || got != tt.want {
			t.Errorf("%s：LongestLine = %d, %v; want %d, nil", tt.name, got, err, tt.want)
		}
	}
	if got, err := LongestLine(errReader{}); got != 0 || !errors.Is(err, errBroken) {
		t.Errorf("讀取失敗時回傳 %d, %v; want 0, %v", got, err, errBroken)
	}
}

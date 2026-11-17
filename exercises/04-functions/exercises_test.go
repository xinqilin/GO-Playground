package functions

import (
	"slices"
	"testing"
)

func TestMinMax(t *testing.T) {
	tests := []struct {
		nums   []int
		lo, hi int
		ok     bool
	}{
		{[]int{85, 90, 15}, 15, 90, true},
		{[]int{7}, 7, 7, true},
		{[]int{-3, -8, -1}, -8, -1, true},
		{nil, 0, 0, false},
	}
	for _, tt := range tests {
		lo, hi, ok := MinMax(tt.nums...)
		if lo != tt.lo || hi != tt.hi || ok != tt.ok {
			t.Errorf("MinMax(%v) = %d, %d, %v; want %d, %d, %v", tt.nums, lo, hi, ok, tt.lo, tt.hi, tt.ok)
		}
	}
	if lo, hi, ok := MinMax(3, 1, 2); lo != 1 || hi != 3 || !ok {
		t.Errorf("MinMax(3, 1, 2) = %d, %d, %v; want 1, 3, true", lo, hi, ok)
	}
}

func TestCountIf(t *testing.T) {
	prices := []int{85, 90, 15, 120, 70}
	if got := CountIf(prices, func(p int) bool { return p >= 85 }); got != 3 {
		t.Errorf("CountIf(價格 >= 85) = %d, want 3", got)
	}
	if got := CountIf(prices, func(p int) bool { return p%2 == 0 }); got != 3 {
		t.Errorf("CountIf(偶數) = %d, want 3", got)
	}
	if got := CountIf(nil, func(int) bool { return true }); got != 0 {
		t.Errorf("CountIf(nil) = %d, want 0", got)
	}
}

func TestNewDispenser(t *testing.T) {
	next, reset := NewDispenser("A")
	for _, want := range []string{"A001", "A002", "A003"} {
		if got := next(); got != want {
			t.Fatalf("next() = %q, want %q", got, want)
		}
	}
	other, _ := NewDispenser("B")
	if got := other(); got != "B001" {
		t.Errorf("另一台機器的 next() = %q, want %q", got, "B001")
	}
	if got := next(); got != "A004" {
		t.Errorf("A 機器的 next() = %q, want %q（不應受 B 影響）", got, "A004")
	}
	reset()
	if got := next(); got != "A001" {
		t.Errorf("reset 之後 next() = %q, want %q", got, "A001")
	}
	if got := other(); got != "B002" {
		t.Errorf("A reset 之後 B 的 next() = %q, want %q", got, "B002")
	}
}

func TestPercentOff(t *testing.T) {
	tests := []struct {
		pct, price, want int
	}{
		{10, 1000, 900},
		{15, 85, 72},
		{0, 85, 85},
		{100, 85, 0},
	}
	for _, tt := range tests {
		if got := PercentOff(tt.pct)(tt.price); got != tt.want {
			t.Errorf("PercentOff(%d)(%d) = %d, want %d", tt.pct, tt.price, got, tt.want)
		}
	}
}

func TestMinusOver(t *testing.T) {
	rule := MinusOver(800, 100)
	tests := []struct{ price, want int }{
		{1000, 900}, {800, 700}, {799, 799}, {0, 0},
	}
	for _, tt := range tests {
		if got := rule(tt.price); got != tt.want {
			t.Errorf("MinusOver(800, 100)(%d) = %d, want %d", tt.price, got, tt.want)
		}
	}
}

func TestApplyRules(t *testing.T) {
	half := func(p int) int { return p / 2 }
	minus50 := func(p int) int { return p - 50 }
	tests := []struct {
		name  string
		price int
		rules []Rule
		want  int
	}{
		{"沒有規則", 500, nil, 500},
		{"先半價再減 50", 500, []Rule{half, minus50}, 200},
		{"先減 50 再半價", 500, []Rule{minus50, half}, 225},
		{"不會變成負數", 30, []Rule{minus50}, 0},
	}
	for _, tt := range tests {
		if got := ApplyRules(tt.price, tt.rules...); got != tt.want {
			t.Errorf("%s：ApplyRules(%d) = %d, want %d", tt.name, tt.price, got, tt.want)
		}
	}
}

func TestProcessAll(t *testing.T) {
	var tr Tracker
	ProcessAll(&tr, []string{"a", "b", "c"})
	want := []string{
		"open a", "use a", "close a",
		"open b", "use b", "close b",
		"open c", "use c", "close c",
	}
	if !slices.Equal(tr.Events, want) {
		t.Errorf("紀錄 = %q\nwant %q", tr.Events, want)
	}
}

func TestSafeDiv(t *testing.T) {
	if q, err := SafeDiv(17, 5); q != 3 || err != nil {
		t.Errorf("SafeDiv(17, 5) = %d, %v; want 3, nil", q, err)
	}
	q, err := SafeDiv(10, 0)
	if q != 0 || err == nil {
		t.Fatalf("SafeDiv(10, 0) = %d, %v; want 0 和一個錯誤", q, err)
	}
	if want := "除法失敗：runtime error: integer divide by zero"; err.Error() != want {
		t.Errorf("錯誤訊息 = %q, want %q", err.Error(), want)
	}
}

func TestPayWays(t *testing.T) {
	tests := []struct {
		amount int
		coins  []int
		want   int
	}{
		{10, []int{1, 5, 10}, 4},
		{6, []int{2, 3}, 2},
		{3, []int{2}, 0},
		{0, []int{1}, 1},
		{0, nil, 1},
		{7, nil, 0},
		{-1, []int{1}, 0},
		{100, []int{1, 5, 10, 50}, 158},
		{100, []int{50, 10, 5, 1}, 158},
	}
	for _, tt := range tests {
		if got := PayWays(tt.amount, tt.coins); got != tt.want {
			t.Errorf("PayWays(%d, %v) = %d, want %d", tt.amount, tt.coins, got, tt.want)
		}
	}
}

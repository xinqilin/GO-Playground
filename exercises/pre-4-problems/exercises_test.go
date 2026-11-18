package problems

import (
	"slices"
	"testing"
)

func TestCartTotal(t *testing.T) {
	tests := []struct {
		prices, qtys []int
		want         int
	}{
		{[]int{85, 90}, []int{2, 1}, 260},
		{[]int{120}, []int{3}, 360},
		{[]int{85, 90, 75}, []int{0, 1, 2}, 240},
		{[]int{}, []int{}, 0},
	}
	for _, tt := range tests {
		if got := CartTotal(tt.prices, tt.qtys); got != tt.want {
			t.Errorf("CartTotal(%v, %v) = %d, want %d", tt.prices, tt.qtys, got, tt.want)
		}
	}
}

func TestContains(t *testing.T) {
	menu := []string{"排骨飯", "雞腿飯", "素食餐"}
	tests := []struct {
		dish string
		want bool
	}{
		{"排骨飯", true},
		{"素食餐", true},
		{"牛肉麵", false},
	}
	for _, tt := range tests {
		if got := Contains(menu, tt.dish); got != tt.want {
			t.Errorf("Contains(%v, %q) = %v, want %v", menu, tt.dish, got, tt.want)
		}
	}
	if Contains([]string{}, "排骨飯") {
		t.Error("空的菜單裡不應該找得到任何菜")
	}
}

func TestBiggestJump(t *testing.T) {
	tests := []struct {
		prices []int
		want   int
	}{
		{[]int{100, 103, 98, 106, 107}, 8},
		{[]int{50, 60}, 10},
		{[]int{100, 90, 80}, 0},
		{[]int{100}, 0},
		{[]int{}, 0},
	}
	for _, tt := range tests {
		if got := BiggestJump(tt.prices); got != tt.want {
			t.Errorf("BiggestJump(%v) = %d, want %d", tt.prices, got, tt.want)
		}
	}
}

func TestAbove(t *testing.T) {
	tests := []struct {
		prices []int
		limit  int
		want   []int
	}{
		{[]int{85, 120, 90, 150}, 100, []int{120, 150}},
		{[]int{85, 100, 101}, 100, []int{101}},
		{[]int{50, 60}, 100, []int{}},
	}
	for _, tt := range tests {
		if got := Above(tt.prices, tt.limit); !slices.Equal(got, tt.want) {
			t.Errorf("Above(%v, %d) = %v, want %v", tt.prices, tt.limit, got, tt.want)
		}
	}
	prices := []int{85, 120, 90}
	Above(prices, 100)
	if !slices.Equal(prices, []int{85, 120, 90}) {
		t.Errorf("Above 不應該改到傳進來的清單，結果變成 %v", prices)
	}
}

func TestMostVoted(t *testing.T) {
	tests := []struct {
		votes []string
		want  string
	}{
		{[]string{"拉麵", "咖哩", "拉麵", "便當", "咖哩"}, "拉麵"},
		{[]string{"便當", "咖哩", "咖哩"}, "咖哩"},
		{[]string{"咖哩", "拉麵", "拉麵", "咖哩"}, "咖哩"},
		{[]string{"便當"}, "便當"},
		{[]string{}, ""},
	}
	for _, tt := range tests {
		if got := MostVoted(tt.votes); got != tt.want {
			t.Errorf("MostVoted(%v) = %q, want %q", tt.votes, got, tt.want)
		}
	}
}

func TestCheapestPrice(t *testing.T) {
	tests := []struct {
		prices []int
		want   int
	}{
		{[]int{85, 90, 75}, 75},
		{[]int{60}, 60},
		{[]int{120, 85, 85}, 85},
	}
	for _, tt := range tests {
		if got := CheapestPrice(tt.prices); got != tt.want {
			t.Errorf("CheapestPrice(%v) = %d, want %d", tt.prices, got, tt.want)
		}
	}
}

func TestLowStock(t *testing.T) {
	tests := []struct {
		stock []int
		limit int
		want  int
	}{
		{[]int{3, 12, 0, 8}, 5, 2},
		{[]int{1, 2, 3}, 5, 3},
		{[]int{10, 20}, 5, 0},
	}
	for _, tt := range tests {
		if got := LowStock(tt.stock, tt.limit); got != tt.want {
			t.Errorf("LowStock(%v, %d) = %d, want %d", tt.stock, tt.limit, got, tt.want)
		}
	}
}

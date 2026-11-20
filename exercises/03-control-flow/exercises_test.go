package controlflow

import (
	"reflect"
	"testing"
)

func TestTicketPrice(t *testing.T) {
	tests := []struct {
		km, want int
	}{
		{-3, 0}, {0, 0},
		{1, 20}, {5, 20},
		{6, 25}, {10, 25},
		{11, 35}, {20, 35},
		{21, 50}, {40, 50},
		{41, 65}, {100, 65},
	}
	for _, tt := range tests {
		if got := TicketPrice(tt.km); got != tt.want {
			t.Errorf("TicketPrice(%d) = %d, want %d", tt.km, got, tt.want)
		}
	}
}

func TestRestock(t *testing.T) {
	items := []Item{{"雞腿", 0}, {"排骨", 2}}
	Restock(items, 5)
	want := []Item{{"雞腿", 5}, {"排骨", 7}}
	if !reflect.DeepEqual(items, want) {
		t.Errorf("after Restock(items, 5): items = %v, want %v", items, want)
	}
}

func TestTable(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{0, ""},
		{1, "1×1=1"},
		{2, "1×1=1\t1×2=2\n2×1=2\t2×2=4"},
		{3, "1×1=1\t1×2=2\t1×3=3\n2×1=2\t2×2=4\t2×3=6\n3×1=3\t3×2=6\t3×3=9"},
	}
	for _, tt := range tests {
		if got := Table(tt.n); got != tt.want {
			t.Errorf("Table(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

func TestInvoicePrize(t *testing.T) {
	const winning = "12345678"
	tests := []struct {
		ticket string
		want   int
	}{
		{"12345678", 200000},
		{"92345678", 40000},
		{"99345678", 10000},
		{"99945678", 4000},
		{"99995678", 1000},
		{"99999678", 200},
		{"99999978", 0},
		{"12345670", 0},
	}
	for _, tt := range tests {
		if got := InvoicePrize(winning, tt.ticket); got != tt.want {
			t.Errorf("InvoicePrize(%q, %q) = %d, want %d", winning, tt.ticket, got, tt.want)
		}
	}
}

func TestCountKinds(t *testing.T) {
	tests := []struct {
		s                         string
		han, letter, digit, other int
	}{
		{"", 0, 0, 0, 0},
		{"ABCxyz", 0, 6, 0, 0},
		{"Go 語言 1.27 版！", 3, 2, 3, 5},
		{"台灣ＮＯ１", 2, 0, 0, 3},
	}
	for _, tt := range tests {
		han, letter, digit, other := CountKinds(tt.s)
		if han != tt.han || letter != tt.letter || digit != tt.digit || other != tt.other {
			t.Errorf("CountKinds(%q) = %d, %d, %d, %d; want %d, %d, %d, %d",
				tt.s, han, letter, digit, other, tt.han, tt.letter, tt.digit, tt.other)
		}
	}
}

func TestTotalUntilNegative(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{nil, 0},
		{[]int{5}, 5},
		{[]int{3, 4, -1, 10}, 7},
		{[]int{-2, 8}, 0},
	}
	for _, tt := range tests {
		if got := TotalUntilNegative(tt.nums); got != tt.want {
			t.Errorf("TotalUntilNegative(%v) = %d, want %d", tt.nums, got, tt.want)
		}
	}
}

func TestSpiral(t *testing.T) {
	if got := Spiral(0); len(got) != 0 {
		t.Errorf("Spiral(0) = %v, want an empty slice", got)
	}
	tests := []struct {
		n    int
		want [][]int
	}{
		{1, [][]int{{1}}},
		{2, [][]int{{1, 2}, {4, 3}}},
		{3, [][]int{{1, 2, 3}, {8, 9, 4}, {7, 6, 5}}},
		{4, [][]int{{1, 2, 3, 4}, {12, 13, 14, 5}, {11, 16, 15, 6}, {10, 9, 8, 7}}},
	}
	for _, tt := range tests {
		if got := Spiral(tt.n); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Spiral(%d) = %v, want %v", tt.n, got, tt.want)
		}
	}
}

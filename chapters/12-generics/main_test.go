package main

import (
	"slices"
	"testing"
)

func Example_whyDemo() {
	whyDemo()
	// Output:
	// 180 3.75
	// 87.5
	// 180 3.75
}

func Example_paramsDemo() {
	paramsDemo()
	// Output:
	// 12
	// 0.3
	// M
	// 5
	// 2
}

func Example_typeSetDemo() {
	typeSetDemo()
	// Output:
	// NT$210
	// main.NTD
	// 3
}

func Example_inferDemo() {
	inferDemo()
	// Output:
	// []float64 [85 95]
	// [85 95]
	// NT$180
}

func Example_setDemo() {
	setDemo()
	// Output:
	// 2 true false
	// true 1
}

func Example_methodDemo() {
	methodDemo()
	// Output:
	// main.Page[string]
	// [17 號 排骨飯 18 號 雞腿飯] 19
}

func Example_whenDemo() {
	whenDemo()
	// Output:
	// NT$85、NT$95
}

func Example_kindDemo() {
	kindDemo()
	// Output:
	// 整數 85 字串 排骨飯 其他 main.NTD
}

func TestClamp(t *testing.T) {
	tests := []struct {
		v, lo, hi, want int
	}{
		{5, 1, 12, 5}, {0, 1, 12, 1}, {20, 1, 12, 12}, {1, 1, 12, 1}, {12, 1, 12, 12},
	}
	for _, tt := range tests {
		if got := Clamp(tt.v, tt.lo, tt.hi); got != tt.want {
			t.Errorf("Clamp(%d, %d, %d) = %d, want %d", tt.v, tt.lo, tt.hi, got, tt.want)
		}
	}
	if got := Clamp("B", "A", "M"); got != "B" {
		t.Errorf("Clamp(B, A, M) = %q, want B", got)
	}
}

func TestSetZeroValue(t *testing.T) {
	var s Set[string]
	if s.Has("x") || s.Len() != 0 {
		t.Errorf("零值的 Set：Has = %v, Len = %d", s.Has("x"), s.Len())
	}
}

func TestPageMapKeepsNext(t *testing.T) {
	p := Page[int]{Items: []int{1, 2}, Next: 0}
	got := p.Map(func(n int) bool { return n%2 == 0 })
	if !slices.Equal(got.Items, []bool{false, true}) || got.Next != 0 {
		t.Errorf("Map = %+v", got)
	}
}

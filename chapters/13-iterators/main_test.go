package main

import (
	"slices"
	"testing"
)

func Example() {
	main()
	// Output:
	// == 自己寫疊代器 ==
	// 3
	// 2
	// 1
	// == 無限數列，break 就停 ==
	// [0 1 1 2 3 5 8 13 21 34]
	// == 把疊代器接起來 ==
	// [0 2 8 34 144]
	// == 標準庫的疊代器 ==
	// [a b c]
	// 0 x
	// 1 y
	// abc
}

func TestTake(t *testing.T) {
	tests := []struct {
		n    int
		want []int
	}{
		{0, nil},
		{1, []int{3}},
		{5, []int{3, 2, 1}}, // 上游只有 3 個，取不滿也要正常結束
	}
	for _, tt := range tests {
		if got := slices.Collect(Take(Countdown(3), tt.n)); !slices.Equal(got, tt.want) {
			t.Errorf("Take(Countdown(3), %d) = %v, want %v", tt.n, got, tt.want)
		}
	}
}

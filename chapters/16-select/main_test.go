package main

import (
	"slices"
	"testing"
)

func Example() {
	main()
	// Output:
	// == 逾時 ==
	// 拿到資料
	// 逾時
	// == default：不等 ==
	// 現在沒有
	// 收到 5
	// == 合併兩個 channel ==
	// [1 2 3 4 5]
}

func TestMerge(t *testing.T) {
	tests := []struct {
		a, b []int
		want []int
	}{
		{nil, nil, nil},
		{[]int{1}, nil, []int{1}},
		{[]int{3, 1}, []int{2}, []int{1, 2, 3}},
	}
	for _, tt := range tests {
		if got := merge(from(tt.a...), from(tt.b...)); !slices.Equal(got, tt.want) {
			t.Errorf("merge(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

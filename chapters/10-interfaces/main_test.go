package main

import (
	"math"
	"testing"
)

func Example() {
	main()
	// Output:
	// == 隱式實作 ==
	// Rect(3x4) 面積 12.00 周長 14.00
	// {1} 面積 3.14 周長 6.28
	// 總面積 15.14
	// == type switch 與 type assertion ==
	// 形狀，面積 2.00
	// 也是 Rect，寬 1
	// 形狀，面積 12.57
	// 整數 43
	// 其他型別 string
	// == nil interface 陷阱 ==
	// false
	// *main.MyErr true
}

func TestTotalArea(t *testing.T) {
	tests := []struct {
		name   string
		shapes []Shape
		want   float64
	}{
		{"沒有形狀", nil, 0},
		{"兩個長方形", []Shape{Rect{1, 2}, Rect{3, 4}}, 14},
		{"單位圓", []Shape{Circle{1}}, math.Pi},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := totalArea(tt.shapes); math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

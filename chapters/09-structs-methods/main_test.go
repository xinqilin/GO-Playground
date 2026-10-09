package main

import "testing"

func Example() {
	main()
	// Output:
	// == 建立 struct ==
	// {0 0} {1 2} {3 4} {5 0}
	// {X:1 Y:2}
	// == value 與 pointer receiver ==
	// 25
	// {4 5}
	// == 嵌入 ==
	// 小黑 柴犬
	// 我是小黑
	// 小黑
	// 我是小白
	// == 比較與匿名 struct ==
	// true
	// localhost 8080
}

func TestMove(t *testing.T) {
	tests := []struct {
		start  Point
		dx, dy int
		want   Point
	}{
		{Point{0, 0}, 1, 2, Point{1, 2}},
		{Point{3, 4}, -3, -4, Point{0, 0}},
	}
	for _, tt := range tests {
		p := tt.start
		p.Move(tt.dx, tt.dy)
		if p != tt.want {
			t.Errorf("%v.Move(%d, %d) = %v, want %v", tt.start, tt.dx, tt.dy, p, tt.want)
		}
	}
}

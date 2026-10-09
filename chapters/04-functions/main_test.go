package main

import "testing"

func Example() {
	main()
	// Output:
	// == 多個回傳值 ==
	// 3 2
	// == 可變參數 ==
	// 0 3 6
	// == 函式是值 ==
	// [2 4 6]
	// == closure ==
	// 1 2 3
	// 1
	// == defer ==
	// 函式本體結束，x = 20
	// 註冊 defer 時 x 是 10
	// defer 2
	// defer 1
	// defer 0
}

// 整數除法往 0 截斷，餘數的正負號跟著被除數。
func TestDivmod(t *testing.T) {
	tests := []struct {
		a, b, q, r int
	}{
		{17, 5, 3, 2},
		{-7, 2, -3, -1},
		{7, -2, -3, 1},
	}
	for _, tt := range tests {
		q, r := divmod(tt.a, tt.b)
		if q != tt.q || r != tt.r {
			t.Errorf("divmod(%d, %d) = (%d, %d), want (%d, %d)", tt.a, tt.b, q, r, tt.q, tt.r)
		}
	}
}

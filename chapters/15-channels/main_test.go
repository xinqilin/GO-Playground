package main

import "testing"

func Example() {
	main()
	// Output:
	// == 無緩衝：一手交錢一手交貨 ==
	// ping
	// == 有緩衝：先放進信箱 ==
	// 2 2
	// 1 2
	// == close 與 range ==
	// 收到 1
	// 收到 2
	// 收到 3
	// == 關閉之後再收 ==
	// 7 true
	// 0 false
	// == 管線：產生 → 平方 → 加總 ==
	// 30
}

func TestPipeline(t *testing.T) {
	tests := []struct{ n, want int }{
		{0, 0},
		{1, 1},
		{10, 385},
	}
	for _, tt := range tests {
		if got := sum(square(produce(tt.n))); got != tt.want {
			t.Errorf("n=%d: got %d, want %d", tt.n, got, tt.want)
		}
	}
}

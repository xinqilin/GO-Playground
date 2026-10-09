package main

import "testing"

func Example() {
	main()
	// Output:
	// == 每個元素一個 goroutine ==
	// [1 4 9 16]
	// == 切段平行加總 ==
	// 500500
	// == goroutine 很便宜 ==
	// 10000
}

func TestSumParallel(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6, 7}
	for _, parts := range []int{0, 1, 2, 3, 7, 20} {
		if got := sumParallel(nums, parts); got != 28 {
			t.Errorf("sumParallel(nums, %d) = %d, want 28", parts, got)
		}
	}
	if got := sumParallel(nil, 3); got != 0 {
		t.Errorf("sumParallel(nil, 3) = %d, want 0", got)
	}
}

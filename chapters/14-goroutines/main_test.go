package main

import (
	"errors"
	"slices"
	"testing"
)

func Example_kitchenDemo() {
	kitchenDemo()
	// Output:
	// 排骨飯好了、雞腿飯好了、魚排飯好了
}

func Example_kitchenAddDone() {
	kitchenAddDone()
	// Output:
	// 排骨飯好了、雞腿飯好了、魚排飯好了
}

func Example_argsNowDemo() {
	argsNowDemo()
	// Output:
	// 廚房收到 排骨飯
	// 櫃台現在的單 雞腿飯
}

func Example_sumDemo() {
	sumDemo()
	// Output:
	// 500500
}

func Example_storesDemo() {
	storesDemo()
	// Output:
	// [12 0 0 7]
	// 圓山店：沒有這家分店
	// true
}

func Example_spawnDemo() {
	spawnDemo()
	// Output:
	// 10000
}

func Example_growDemo() {
	growDemo()
	// Output:
	// 100000
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

func TestCheckStores(t *testing.T) {
	counts, err := checkStores([]string{"民權店", "中山店"})
	if err != nil || !slices.Equal(counts, []int{7, 12}) {
		t.Errorf("checkStores = %v, %v; want [7 12], nil", counts, err)
	}
	_, err = checkStores([]string{"台北店", "中山店", "圓山店"})
	if !errors.Is(err, ErrNoStore) || err.Error() != "台北店：沒有這家分店\n圓山店：沒有這家分店" {
		t.Errorf("兩家查不到時 err = %q", err)
	}
}

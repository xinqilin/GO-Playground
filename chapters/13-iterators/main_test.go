package main

import (
	"fmt"
	"iter"
	"slices"
	"testing"
)

func Example_rangeFuncDemo() {
	rangeFuncDemo()
	// Output:
	// 請 17 號
	// 請 18 號
	// 請 19 號
}

func Example_byHand() {
	byHand()
	// Output:
	// 請 17 號
	// 請 18 號
	// 請 19 號
}

func Example_earlyStopDemo() {
	earlyStopDemo()
	// Output:
	// 打開收據匣
	// 處理 A-001
	// 處理 A-002
	// 關上收據匣
	// 迴圈結束
}

func Example_ignoreYieldDemo() {
	ignoreYieldDemo()
	// Output:
	// 1
	// panic: runtime error: range function continued iteration after function for loop body returned false
}

func Example_seq2Demo() {
	seq2Demo()
	// Output:
	// 排骨飯 85
	// 雞腿飯 95
	// 紅茶 25
}

func Example_treeDemo() {
	treeDemo()
	// Output:
	// 菜單
	// 菜單/便當
	// 菜單/便當/排骨飯
	// 菜單/便當/雞腿飯
	// 菜單/飲料
	// 菜單/飲料/紅茶
	// 找到 菜單/便當/雞腿飯
}

func Example_stdlibDemo() {
	stdlibDemo()
	// Output:
	// [排骨飯 滷蛋 雞腿飯]
	// 2 19
	// 1 18
	// 0 17
	// "排骨飯 x2\n"
	// "紅茶 x1\n"
	// "合計 195\n"
	// [1 2]
	// [3 4]
	// [5]
}

func Example_pipelineDemo() {
	pipelineDemo()
	// Output:
	// [001 002 003 005 006]
}

func Example_lazyDemo() {
	lazyDemo()
	// Output:
	// 產生 1
	// 加倍 1
	// 拿到 2
	// 產生 2
	// 加倍 2
	// 拿到 4
}

func Example_pullBasics() {
	pullBasics()
	// Output:
	// A-001 true A-002 true
	// "" false
}

func Example_pullDemo() {
	pullDemo()
	// Output:
	// 17 → 1 號櫃台
	// 18 → 2 號櫃台
	// 19 → 3 號櫃台
}

func Example_reuseDemo() {
	reuseDemo()
	// Output:
	// [1 2 3] [1 2 3]
	// [排骨飯 紅茶] []
}

// counting 包住一個疊代器，記錄上游一共交出了幾個值。
func counting[T any](seq iter.Seq[T], n *int) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			*n++
			if !yield(v) {
				return
			}
		}
	}
}

func TestTakeStopsUpstream(t *testing.T) {
	produced := 0
	got := slices.Collect(Take(counting(naturals(), &produced), 3))
	if !slices.Equal(got, []int{1, 2, 3}) || produced != 3 {
		t.Errorf("Take(3) = %v，上游產生了 %d 個, want [1 2 3]，3 個", got, produced)
	}
	if got := slices.Collect(Take(naturals(), 0)); len(got) != 0 {
		t.Errorf("Take(0) = %v, want 空的", got)
	}
}

func TestFilterMap(t *testing.T) {
	odd := Filter(callNumbers(1, 6), func(n int) bool { return n%2 == 1 })
	got := slices.Collect(Map(odd, func(n int) int { return n * 10 }))
	if !slices.Equal(got, []int{10, 30, 50}) {
		t.Errorf("Filter + Map = %v, want [10 30 50]", got)
	}
}

func TestWalkStopsEarly(t *testing.T) {
	tree := &Category{Name: "a", Children: []*Category{
		{Name: "b", Children: []*Category{{Name: "c"}}},
		{Name: "d"},
	}}
	var seen []string
	for p := range tree.Walk() {
		seen = append(seen, p)
		if p == "a/b/c" {
			break
		}
	}
	if !slices.Equal(seen, []string{"a", "a/b", "a/b/c"}) {
		t.Errorf("Walk 提早結束時看到 %q", seen)
	}
}

func TestAssignStopsWhenNumbersRunOut(t *testing.T) {
	var got []string
	for n, c := range assign(slices.Values([]int{1}), slices.Values([]string{"x", "y"})) {
		got = append(got, fmt.Sprint(n, c))
	}
	if !slices.Equal(got, []string{"1x"}) {
		t.Errorf("assign = %q, want [1x]", got)
	}
}

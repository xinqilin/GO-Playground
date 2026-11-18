package main

import "testing"

func Example_greetDemo() {
	greetDemo()
	// Output:
	// 早安，小明
	// 早安，阿美
}

func Example_subtotalDemo() {
	subtotalDemo()
	// Output:
	// 255 180 435
}

func Example_memberPriceDemo() {
	memberPriceDemo()
	// Output:
	// 100
	// 90
	// 76
}

func Example_hoursMinutesDemo() {
	hoursMinutesDemo()
	// Output:
	// 2 小時 15 分
}

func Example_copyDemo() {
	copyDemo()
	// Output:
	// 函式裡： 6
	// 函式外： 5
}

func Example_listBasics() {
	listBasics()
	// Output:
	// [85 90 75]
	// 85 75
	// 3
	// [85 95 75]
	// [85 95 75 100] 4
}

func Example_stringList() {
	stringList()
	// Output:
	// 排骨飯
	// 素食餐
}

func Example_emptyList() {
	emptyList()
	// Output:
	// 0
	// [120 85] 2
}

func Example_loopIndex() {
	loopIndex()
	// Output:
	// 0 排骨飯
	// 1 雞腿飯
	// 2 素食餐
}

func Example_loopRange() {
	loopRange()
	// Output:
	// 第 1 個： 85 元
	// 第 2 個： 90 元
	// 第 3 個： 75 元
}

func Example_sumDemo() {
	sumDemo()
	// Output:
	// 19800
}

func Example_countDemo() {
	countDemo()
	// Output:
	// 3
}

func Example_maxDemo() {
	maxDemo()
	// Output:
	// 5600
}

func Example_wrongMaxDemo() {
	wrongMaxDemo()
	// Output:
	// 0 -1
}

func Example_searchDemo() {
	searchDemo()
	// Output:
	// 1
	// -1
}

func Example_weeklyReport() {
	weeklyReport()
	// Output:
	// 總營業額： 31100
	// 平均每天： 4442
	// 最好的一天：星期六 6300
	// 達到 4500 的天數： 3
}

func Example_debugDemo() {
	debugDemo()
	// Output:
	// i = 1 total = 90
	// i = 2 total = 165
	// 165
}

func TestMaxOf(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"一般", []int{3, 9, 4}, 9},
		{"最大的在第一個", []int{9, 3, 4}, 9},
		{"最大的在最後", []int{3, 4, 9}, 9},
		{"只有一個", []int{7}, 7},
		{"全是負數", []int{-3, -5, -1}, -1},
		{"有一樣大的", []int{5, 5, 2}, 5},
	}
	for _, tt := range tests {
		if got := maxOf(tt.nums); got != tt.want {
			t.Errorf("%s：maxOf(%v) = %d, want %d", tt.name, tt.nums, got, tt.want)
		}
	}
}

func TestIndexOfMax(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{[]int{3, 9, 4}, 1},
		{[]int{9, 3, 4}, 0},
		{[]int{5, 9, 9}, 1},
		{[]int{-3, -5, -1}, 2},
	}
	for _, tt := range tests {
		if got := indexOfMax(tt.nums); got != tt.want {
			t.Errorf("indexOfMax(%v) = %d, want %d", tt.nums, got, tt.want)
		}
	}
}

func TestIndexOf(t *testing.T) {
	queue := []string{"小明", "阿美", "小華", "阿美"}
	tests := []struct {
		target string
		want   int
	}{
		{"小明", 0},
		{"阿美", 1},
		{"小華", 2},
		{"小芬", -1},
	}
	for _, tt := range tests {
		if got := indexOf(queue, tt.target); got != tt.want {
			t.Errorf("indexOf(%v, %q) = %d, want %d", queue, tt.target, got, tt.want)
		}
	}
}

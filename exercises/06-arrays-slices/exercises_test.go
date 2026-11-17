package slicing

import (
	"slices"
	"testing"
)

func TestChunk(t *testing.T) {
	orders := []string{"o1", "o2", "o3", "o4", "o5"}
	tests := []struct {
		size int
		want [][]string
	}{
		{2, [][]string{{"o1", "o2"}, {"o3", "o4"}, {"o5"}}},
		{5, [][]string{{"o1", "o2", "o3", "o4", "o5"}}},
		{10, [][]string{{"o1", "o2", "o3", "o4", "o5"}}},
		{1, [][]string{{"o1"}, {"o2"}, {"o3"}, {"o4"}, {"o5"}}},
		{0, nil},
	}
	for _, tt := range tests {
		got := Chunk(orders, tt.size)
		if !slices.EqualFunc(got, tt.want, slices.Equal) {
			t.Errorf("Chunk(5 筆, %d) = %v, want %v", tt.size, got, tt.want)
		}
	}
	if got := Chunk(nil, 3); got != nil {
		t.Errorf("Chunk(nil, 3) = %v, want nil", got)
	}
}

func TestDiscounted(t *testing.T) {
	prices := []int{85, 90, 100}
	got := Discounted(prices, 10)
	if want := []int{76, 81, 90}; !slices.Equal(got, want) {
		t.Errorf("Discounted = %v, want %v", got, want)
	}
	if want := []int{85, 90, 100}; !slices.Equal(prices, want) {
		t.Errorf("原本的價格被改成 %v, want %v", prices, want)
	}
	if got := Discounted(nil, 10); len(got) != 0 {
		t.Errorf("Discounted(nil, 10) = %v, want 空的", got)
	}
}

func TestCompact(t *testing.T) {
	codes := []string{"A1", "A1", "B2", "A1", "C3", "C3", "C3"}
	got := Compact(codes)
	if want := []string{"A1", "B2", "A1", "C3"}; !slices.Equal(got, want) {
		t.Fatalf("Compact = %v, want %v", got, want)
	}
	if &got[0] != &codes[0] {
		t.Error("結果應該沿用 codes 的底層陣列，不要配置新的")
	}
	if got := Compact([]string{"A1"}); !slices.Equal(got, []string{"A1"}) {
		t.Errorf("Compact([A1]) = %v", got)
	}
	if got := Compact(nil); len(got) != 0 {
		t.Errorf("Compact(nil) = %v, want 空的", got)
	}
}

func TestSplitAt(t *testing.T) {
	orders := []int{1, 2, 3, 4, 5}
	head, tail := SplitAt(orders, 2)
	if !slices.Equal(head, []int{1, 2}) || !slices.Equal(tail, []int{3, 4, 5}) {
		t.Fatalf("SplitAt(2) = %v, %v; want [1 2], [3 4 5]", head, tail)
	}
	head = append(head, 99)
	if tail[0] != 3 || orders[2] != 3 {
		t.Errorf("對 head append 之後，tail = %v、orders = %v；兩者都不應該改變", tail, orders)
	}
	if !slices.Equal(head, []int{1, 2, 99}) {
		t.Errorf("append 之後 head = %v", head)
	}

	head, tail = SplitAt([]int{1, 2, 3}, -1)
	if len(head) != 0 || !slices.Equal(tail, []int{1, 2, 3}) {
		t.Errorf("SplitAt(-1) = %v, %v", head, tail)
	}
	head, tail = SplitAt([]int{1, 2, 3}, 10)
	if !slices.Equal(head, []int{1, 2, 3}) || len(tail) != 0 {
		t.Errorf("SplitAt(10) = %v, %v", head, tail)
	}
}

func TestRotate(t *testing.T) {
	tests := []struct {
		k    int
		want []string
	}{
		{1, []string{"B", "C", "D", "A"}},
		{3, []string{"D", "A", "B", "C"}},
		{4, []string{"A", "B", "C", "D"}},
		{5, []string{"B", "C", "D", "A"}},
		{0, []string{"A", "B", "C", "D"}},
	}
	for _, tt := range tests {
		staff := []string{"A", "B", "C", "D"}
		Rotate(staff, tt.k)
		if !slices.Equal(staff, tt.want) {
			t.Errorf("Rotate([A B C D], %d) 之後是 %v, want %v", tt.k, staff, tt.want)
		}
	}
	Rotate(nil, 3) // 不能 panic
}

func TestBestStreak(t *testing.T) {
	tests := []struct {
		sales        []int
		k            int
		start, total int
	}{
		{[]int{3, 8, 2, 9, 9, 1}, 2, 3, 18},
		{[]int{3, 8, 2, 9, 9, 1}, 1, 3, 9},
		{[]int{3, 8, 2, 9, 9, 1}, 6, 0, 32},
		{[]int{1, 4, 4, 1}, 2, 1, 8},
		{[]int{5, 1, 5, 1}, 2, 0, 6},
		{[]int{3, 8}, 3, -1, 0},
		{[]int{3, 8}, 0, -1, 0},
		{nil, 1, -1, 0},
	}
	for _, tt := range tests {
		start, total := BestStreak(tt.sales, tt.k)
		if start != tt.start || total != tt.total {
			t.Errorf("BestStreak(%v, %d) = %d, %d; want %d, %d", tt.sales, tt.k, start, total, tt.start, tt.total)
		}
	}
}

func TestInsertSorted(t *testing.T) {
	tests := []struct {
		s    []int
		v    int
		want []int
	}{
		{[]int{1, 3, 5}, 4, []int{1, 3, 4, 5}},
		{[]int{1, 3, 5}, 0, []int{0, 1, 3, 5}},
		{[]int{1, 3, 5}, 9, []int{1, 3, 5, 9}},
		{[]int{1, 3, 3, 5}, 3, []int{1, 3, 3, 3, 5}},
		{nil, 7, []int{7}},
	}
	for _, tt := range tests {
		if got := InsertSorted(slices.Clone(tt.s), tt.v); !slices.Equal(got, tt.want) {
			t.Errorf("InsertSorted(%v, %d) = %v, want %v", tt.s, tt.v, got, tt.want)
		}
	}

	roomy := make([]int, 3, 10)
	copy(roomy, []int{1, 3, 5})
	got := InsertSorted(roomy, 2)
	if !slices.Equal(got, []int{1, 2, 3, 5}) {
		t.Fatalf("InsertSorted = %v", got)
	}
	if &got[0] != &roomy[0] {
		t.Error("容量足夠時應該原地插入，沿用 s 的底層陣列")
	}

	full := []int{1, 3, 5}
	got = InsertSorted(full[:3:3], 2)
	if !slices.Equal(got, []int{1, 2, 3, 5}) || !slices.Equal(full, []int{1, 3, 5}) {
		t.Errorf("容量不夠時：結果 %v、原本的 s 變成 %v；原本的 s 不能被改", got, full)
	}
}

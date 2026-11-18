package lineup

import (
	"iter"
	"slices"
	"testing"
)

// firstN 從 seq 取前 n 個就 break。疊代器沒有遵守 yield 的約定而 panic 時，
// 把 panic 的內容放在 panicked 回傳，讓測試回報失敗，而不是整個測試程式當掉。
func firstN[T any](seq iter.Seq[T], n int) (got []T, panicked any) {
	defer func() {
		panicked = recover()
	}()
	for v := range seq {
		got = append(got, v)
		if len(got) == n {
			break
		}
	}
	return got, nil
}

func TestEvery(t *testing.T) {
	tests := []struct {
		from, to, step int
		want           []int
	}{
		{360, 420, 15, []int{360, 375, 390, 405, 420}},
		{1, 10, 4, []int{1, 5, 9}},
		{5, 5, 1, []int{5}},
		{6, 5, 1, nil},
		{1, 10, 0, nil},
		{1, 10, -1, nil},
	}
	for _, tt := range tests {
		got := slices.Collect(Every(tt.from, tt.to, tt.step))
		if !slices.Equal(got, tt.want) {
			t.Errorf("Every(%d, %d, %d) = %v, want %v", tt.from, tt.to, tt.step, got, tt.want)
		}
	}
	if got, p := firstN(Every(0, 100, 10), 2); p != nil || !slices.Equal(got, []int{0, 10}) {
		t.Errorf("取兩個就 break：拿到 %v，panic: %v", got, p)
	}
}

func TestNumbered(t *testing.T) {
	var nums []int
	var lines []string
	for i, l := range Numbered("排骨飯 x2\n紅茶 x1\n合計 195") {
		nums = append(nums, i)
		lines = append(lines, l)
	}
	if !slices.Equal(nums, []int{1, 2, 3}) || !slices.Equal(lines, []string{"排骨飯 x2", "紅茶 x1", "合計 195"}) {
		t.Errorf("Numbered = %v %q", nums, lines)
	}

	count := 0
	for range Numbered("") {
		count++
	}
	if count != 0 {
		t.Errorf("Numbered(\"\") 交出了 %d 行", count)
	}

	var got []string
	for i, l := range Numbered("a\n\nb\n") {
		got = append(got, l)
		if i == 2 {
			break
		}
	}
	if !slices.Equal(got, []string{"a", ""}) {
		t.Errorf("空行也算一行，取兩行就停：%q", got)
	}
}

func TestDedup(t *testing.T) {
	got := slices.Collect(Dedup(slices.Values([]int{17, 17, 18, 18, 18, 17})))
	if !slices.Equal(got, []int{17, 18, 17}) {
		t.Errorf("Dedup = %v, want [17 18 17]", got)
	}
	got2 := slices.Collect(Dedup(slices.Values([]string{"", "", "a"})))
	if !slices.Equal(got2, []string{"", "a"}) {
		t.Errorf("Dedup 第一個值剛好是零值時 = %q, want [\"\" \"a\"]", got2)
	}
	if got3 := slices.Collect(Dedup(slices.Values([]int{}))); len(got3) != 0 {
		t.Errorf("Dedup(空的) = %v", got3)
	}
	if got4, p := firstN(Dedup(slices.Values([]int{1, 1, 2, 3})), 2); p != nil || !slices.Equal(got4, []int{1, 2}) {
		t.Errorf("取兩個就 break：拿到 %v，panic: %v", got4, p)
	}
}

func TestSegments(t *testing.T) {
	var got []string
	for a, b := range Segments(slices.Values([]string{"台北", "中山", "雙連", "民權西路"})) {
		got = append(got, a+"-"+b)
	}
	if want := []string{"台北-中山", "中山-雙連", "雙連-民權西路"}; !slices.Equal(got, want) {
		t.Errorf("Segments = %q, want %q", got, want)
	}
	count := 0
	for range Segments(slices.Values([]int{1})) {
		count++
	}
	if count != 0 {
		t.Errorf("只有一個站時 Segments 交出了 %d 段", count)
	}
	got = nil
	for a, b := range Segments(slices.Values([]string{"a", "b", "c"})) {
		got = append(got, a+b)
		break
	}
	if !slices.Equal(got, []string{"ab"}) {
		t.Errorf("取一段就 break：%q", got)
	}
}

func TestBatch(t *testing.T) {
	got := slices.Collect(Batch(slices.Values([]int{1, 2, 3, 4, 5}), 2))
	want := [][]int{{1, 2}, {3, 4}, {5}}
	if len(got) != len(want) {
		t.Fatalf("Batch 交出 %d 批 %v, want %v", len(got), got, want)
	}
	for i := range want {
		if !slices.Equal(got[i], want[i]) {
			t.Errorf("第 %d 批 = %v, want %v（每一批要是新的 slice）", i+1, got[i], want[i])
		}
	}
	if got := slices.Collect(Batch(slices.Values([]int{1, 2}), 0)); len(got) != 0 {
		t.Errorf("n 是 0 時 Batch = %v", got)
	}
	if got := slices.Collect(Batch(slices.Values([]int{}), 3)); len(got) != 0 {
		t.Errorf("空的輸入 Batch = %v", got)
	}
	if got, p := firstN(Batch(slices.Values([]int{1, 2, 3, 4, 5}), 2), 1); p != nil || len(got) != 1 {
		t.Errorf("取一批就 break：拿到 %v，panic: %v", got, p)
	}
}

// tracked 回傳一個疊代器，結束時（不論是跑完還是被喊停）把 *closed 設成 true。
func tracked(vals []int, closed *bool) iter.Seq[int] {
	return func(yield func(int) bool) {
		defer func() { *closed = true }()
		for _, v := range vals {
			if !yield(v) {
				return
			}
		}
	}
}

func TestMergeSorted(t *testing.T) {
	var ca, cb bool
	got := slices.Collect(MergeSorted(tracked([]int{1, 4, 9}, &ca), tracked([]int{2, 3, 10, 11}, &cb)))
	if !slices.Equal(got, []int{1, 2, 3, 4, 9, 10, 11}) {
		t.Errorf("MergeSorted = %v", got)
	}
	if got := slices.Collect(MergeSorted(slices.Values([]string{}), slices.Values([]string{"a", "b"}))); !slices.Equal(got, []string{"a", "b"}) {
		t.Errorf("一邊是空的：MergeSorted = %q", got)
	}

	ca, cb = false, false
	got, p := firstN(MergeSorted(tracked([]int{1, 5, 7}, &ca), tracked([]int{2, 6}, &cb)), 2)
	if p != nil || !slices.Equal(got, []int{1, 2}) {
		t.Errorf("取兩個就 break：拿到 %v，panic: %v", got, p)
	}
	if !ca || !cb {
		t.Errorf("提早結束後，兩個來源都應該被停止（呼叫 stop）：a 停了 %v，b 停了 %v", ca, cb)
	}
}

func TestSeats(t *testing.T) {
	if got := slices.Collect(Seats("A", 3)); !slices.Equal(got, []string{"A1", "A2", "A3"}) {
		t.Errorf("Seats(A, 3) = %q", got)
	}
	got, p := firstN(Seats("B", 5), 2)
	if p != nil {
		t.Errorf("取兩個就 break 時 panic: %v", p)
	}
	if !slices.Equal(got, []string{"B1", "B2"}) {
		t.Errorf("取兩個 = %q, want [B1 B2]", got)
	}
}

package menu

import (
	"errors"
	"slices"
	"strconv"
	"testing"
)

func TestFilter(t *testing.T) {
	nums := []int{85, 120, 95, 150}
	got := Filter(nums, func(n int) bool { return n >= 100 })
	if !slices.Equal(got, []int{120, 150}) {
		t.Errorf("Filter(>= 100) = %v, want [120 150]", got)
	}
	if !slices.Equal(nums, []int{85, 120, 95, 150}) {
		t.Errorf("Filter 改到了原本的 slice：%v", nums)
	}

	words := []string{"排骨飯", "紅茶", "雞腿飯"}
	got2 := Filter(words, func(s string) bool { return len(s) > 6 })
	if !slices.Equal(got2, []string{"排骨飯", "雞腿飯"}) {
		t.Errorf("Filter(字串) = %q", got2)
	}

	if got3 := Filter([]int{1, 3}, func(n int) bool { return n%2 == 0 }); len(got3) != 0 {
		t.Errorf("沒有符合的元素時 Filter = %v, want 空的", got3)
	}
}

// ntd 是測試用的自訂型別，底層型別是 int。
type ntd int

func TestAverage(t *testing.T) {
	if got, ok := Average([]int{85, 95, 30}); got != 70 || !ok {
		t.Errorf("Average([]int) = %v, %v, want 70, true", got, ok)
	}
	if got, ok := Average([]float64{1.5, 2.5}); got != 2 || !ok {
		t.Errorf("Average([]float64) = %v, %v, want 2, true", got, ok)
	}
	if got, ok := Average([]ntd{85, 90}); got != 87.5 || !ok {
		t.Errorf("Average([]ntd) = %v, %v, want 87.5, true", got, ok)
	}
	if got, ok := Average([]int64{}); got != 0 || ok {
		t.Errorf("Average(空的) = %v, %v, want 0, false", got, ok)
	}
}

type order struct {
	Customer string
	Item     string
}

func TestGroupBy(t *testing.T) {
	orders := []order{
		{"小明", "排骨飯"}, {"阿美", "雞腿飯"}, {"小明", "紅茶"}, {"阿美", "滷蛋"}, {"阿國", "排骨飯"},
	}
	got := GroupBy(orders, func(o order) string { return o.Customer })
	want := map[string][]order{
		"小明": {{"小明", "排骨飯"}, {"小明", "紅茶"}},
		"阿美": {{"阿美", "雞腿飯"}, {"阿美", "滷蛋"}},
		"阿國": {{"阿國", "排骨飯"}},
	}
	if len(got) != len(want) {
		t.Fatalf("GroupBy 分成 %d 組, want %d 組：%v", len(got), len(want), got)
	}
	for k, w := range want {
		if !slices.Equal(got[k], w) {
			t.Errorf("GroupBy[%s] = %v, want %v", k, got[k], w)
		}
	}

	byLen := GroupBy([]string{"a", "bb", "c"}, func(s string) int { return len(s) })
	if !slices.Equal(byLen[1], []string{"a", "c"}) || !slices.Equal(byLen[2], []string{"bb"}) {
		t.Errorf("GroupBy(依長度) = %v", byLen)
	}
}

func TestTally(t *testing.T) {
	var dishes Tally[string]
	for _, d := range []string{"排骨飯", "雞腿飯", "排骨飯", "紅茶", "雞腿飯", "排骨飯", "滷蛋"} {
		dishes.Add(d)
	}
	if got := dishes.Count("排骨飯"); got != 3 {
		t.Errorf("Count(排骨飯) = %d, want 3", got)
	}
	if got := dishes.Count("魚排飯"); got != 0 {
		t.Errorf("Count(魚排飯) = %d, want 0", got)
	}
	if got, want := dishes.Top(2), []string{"排骨飯", "雞腿飯"}; !slices.Equal(got, want) {
		t.Errorf("Top(2) = %q, want %q", got, want)
	}
	if got, want := dishes.Top(10), []string{"排骨飯", "雞腿飯", "滷蛋", "紅茶"}; !slices.Equal(got, want) {
		t.Errorf("Top(10) = %q, want %q（次數一樣時 key 小的在前）", got, want)
	}

	var tables Tally[int]
	tables.Add(5)
	tables.Add(3)
	if got := tables.Top(1); !slices.Equal(got, []int{3}) {
		t.Errorf("Tally[int].Top(1) = %v, want [3]", got)
	}

	var empty Tally[string]
	if got := empty.Top(3); len(got) != 0 {
		t.Errorf("空的 Tally Top = %v", got)
	}
}

func TestResultThen(t *testing.T) {
	errRange := errors.New("人數不在 1 到 12 之間")
	checkRange := func(n int) (int, error) {
		if n < 1 || n > 12 {
			return 0, errRange
		}
		return n, nil
	}

	r := Result[string]{Val: "4"}.Then(strconv.Atoi).Then(checkRange)
	if r.Val != 4 || r.Err != nil {
		t.Errorf(`"4" 的結果 = %+v, want {Val:4 Err:<nil>}`, r)
	}

	r = Result[string]{Val: "20"}.Then(strconv.Atoi).Then(checkRange)
	if !errors.Is(r.Err, errRange) {
		t.Errorf(`"20" 的錯誤 = %v, want %v`, r.Err, errRange)
	}

	calls := 0
	errEarly := errors.New("一開始就失敗")
	r2 := Result[int]{Err: errEarly}.Then(func(n int) (string, error) {
		calls++
		return "不應該被呼叫", nil
	})
	if calls != 0 || !errors.Is(r2.Err, errEarly) || r2.Val != "" {
		t.Errorf("已經失敗時：呼叫 f %d 次，結果 %+v", calls, r2)
	}
}

func TestMemo(t *testing.T) {
	calls := 0
	fare := Memo(func(stations int) int {
		calls++
		return 20 + stations*5
	})
	if got := fare(3); got != 35 {
		t.Errorf("fare(3) = %d, want 35", got)
	}
	if got := fare(3); got != 35 {
		t.Errorf("第二次 fare(3) = %d, want 35", got)
	}
	if got := fare(10); got != 70 {
		t.Errorf("fare(10) = %d, want 70", got)
	}
	if calls != 2 {
		t.Errorf("f 被呼叫了 %d 次, want 2（同一個參數只算一次）", calls)
	}

	other := Memo(func(s string) int { calls++; return len(s) })
	other("排骨飯")
	other("排骨飯")
	if calls != 3 {
		t.Errorf("兩個 Memo 應該各自有自己的記憶：f 總共被呼叫了 %d 次, want 3", calls)
	}
}

func TestMax(t *testing.T) {
	if got, ok := Max([]int{-3, -1, -7}); got != -1 || !ok {
		t.Errorf("Max(-3, -1, -7) = %d, %v, want -1, true", got, ok)
	}
	if got, ok := Max([]int{85, 95, 30}); got != 95 || !ok {
		t.Errorf("Max(85, 95, 30) = %d, %v, want 95, true", got, ok)
	}
	if got, ok := Max([]float64{-2.5}); got != -2.5 || !ok {
		t.Errorf("Max(-2.5) = %v, %v, want -2.5, true", got, ok)
	}
	if got, ok := Max([]string{"排骨飯", "雞腿飯"}); got != "雞腿飯" || !ok {
		t.Errorf("Max(字串) = %q, %v, want 雞腿飯, true", got, ok)
	}
	if got, ok := Max([]int{}); got != 0 || ok {
		t.Errorf("Max(空的) = %d, %v, want 0, false", got, ok)
	}
}

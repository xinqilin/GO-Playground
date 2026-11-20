package testkit

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"
)

func TestSubtestNames(t *testing.T) {
	tests := []struct {
		name     string
		in, want []string
	}{
		{"空白換成底線", []string{"a b", "停 30 分鐘", "tab\there"}, []string{"a_b", "停_30_分鐘", "tab_here"}},
		{"重複加編號", []string{"x", "x", "y", "x"}, []string{"x", "x#01", "y", "x#02"}},
		{"空字串", []string{"", "", "z"}, []string{"#00", "#01", "z"}},
		{"換完才比重複", []string{"a b", "a_b"}, []string{"a_b", "a_b#01"}},
		{"沒有子測試", nil, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SubtestNames(tt.in)
			if len(got) != len(tt.want) || (len(got) > 0 && !slices.Equal(got, tt.want)) {
				t.Errorf("SubtestNames(%q) = %q，want %q", tt.in, got, tt.want)
			}
		})
	}
	// 用真的 t.Run 對照：前兩組的名稱應該和 go test 實際用的一樣。
	var real []string
	for _, n := range []string{"a b", "x", "x", ""} {
		t.Run(n, func(t *testing.T) {
			real = append(real, t.Name()[len("TestSubtestNames/"):])
		})
	}
	if got := SubtestNames([]string{"a b", "x", "x", ""}); !slices.Equal(got, real) {
		t.Errorf("SubtestNames = %q，go test 實際用的是 %q", got, real)
	}
}

func TestDiffLines(t *testing.T) {
	tests := []struct {
		name, got, want string
		diff            []string
	}{
		{"一樣", "排骨飯 90\n滷肉飯 50\n", "排骨飯 90\n滷肉飯 50\n", nil},
		{"一行不同", "排骨飯 90\n滷肉飯 50", "排骨飯 95\n滷肉飯 50", []string{`第 1 行：got "排骨飯 90"，want "排骨飯 95"`}},
		{"got 少一行", "a\nb", "a\nb\nc", []string{`第 3 行：got （沒有這一行），want "c"`}},
		{"got 多一行", "a\nb\n", "a\nb", []string{`第 3 行：got ""，want （沒有這一行）`}},
		{"好幾行", "a\nX\nc\nY", "a\nb\nc\nd", []string{`第 2 行：got "X"，want "b"`, `第 4 行：got "Y"，want "d"`}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DiffLines(tt.got, tt.want); !slices.Equal(got, tt.diff) {
				t.Errorf("DiffLines(%q, %q) =\n  %q\nwant\n  %q", tt.got, tt.want, got, tt.diff)
			}
		})
	}
}

func TestRecorder(t *testing.T) {
	r := &Recorder{FailOn: map[string]bool{"0900-000-000": true}}
	if err := r.Notify("0900-000-000", "您的車已出場"); !errors.Is(err, ErrUnreachable) {
		t.Errorf("FailOn 的號碼應該回傳 ErrUnreachable，拿到 %v", err)
	}
	if err := r.Notify("0912-345-678", "停車費 80 元"); err != nil {
		t.Fatalf("Notify：%v", err)
	}
	if got, want := r.Sent(), []string{"0912-345-678：停車費 80 元"}; !slices.Equal(got, want) {
		t.Fatalf("Sent() = %q，want %q", got, want)
	}

	// Sent 回傳的是複本。
	r.Sent()[0] = "被改掉了"
	if got := r.Sent()[0]; got != "0912-345-678：停車費 80 元" {
		t.Errorf("改了 Sent() 的結果，Recorder 裡也跟著變成 %q", got)
	}

	// 同時送 100 則（用 -race 跑可以抓到沒有加鎖）。
	var wg sync.WaitGroup
	for i := range 100 {
		wg.Go(func() { r.Notify(fmt.Sprintf("09%08d", i), "通知") })
	}
	wg.Wait()
	sent := r.Sent()
	if len(sent) != 101 {
		t.Fatalf("同時送 100 則之後共 %d 則，want 101", len(sent))
	}
	slices.Sort(sent[1:])
	if got, want := sent[1], "0900000000：通知"; got != want {
		t.Errorf("排序後第一則 = %q，want %q", got, want)
	}
}

// spyTB 包住真正的 testing.TB，記下 Golden 回報了什麼。
// Fatalf 和真的一樣會呼叫 runtime.Goexit 結束 goroutine，所以要在另一個 goroutine 裡跑。
type spyTB struct {
	testing.TB
	helper         bool
	errors, fatals []string
}

func (s *spyTB) Helper() { s.helper = true }
func (s *spyTB) Errorf(format string, args ...any) {
	s.errors = append(s.errors, fmt.Sprintf(format, args...))
}
func (s *spyTB) Fatalf(format string, args ...any) {
	s.fatals = append(s.fatals, fmt.Sprintf(format, args...))
	runtime.Goexit()
}

func runGolden(t *testing.T, path string, got []byte, update bool) *spyTB {
	spy := &spyTB{TB: t}
	done := make(chan struct{})
	go func() {
		defer close(done)
		Golden(spy, path, got, update)
	}()
	<-done
	return spy
}

func TestGolden(t *testing.T) {
	path := filepath.Join(t.TempDir(), "testdata", "report.golden")

	spy := runGolden(t, path, []byte("合計 80 元\n"), false)
	if len(spy.fatals) != 1 || !strings.Contains(spy.fatals[0], path) || !strings.Contains(spy.fatals[0], "-update") {
		t.Errorf("golden file 不存在時應該 Fatalf 一次，訊息有路徑和 -update，拿到 %q", spy.fatals)
	}

	spy = runGolden(t, path, []byte("合計 80 元\n"), true)
	if data, err := os.ReadFile(path); err != nil || string(data) != "合計 80 元\n" {
		t.Fatalf("update 之後 %s 的內容 = %q, %v", path, data, err)
	}
	if len(spy.errors)+len(spy.fatals) != 0 {
		t.Errorf("update 不應該回報失敗，拿到 %q %q", spy.errors, spy.fatals)
	}

	spy = runGolden(t, path, []byte("合計 80 元\n"), false)
	if len(spy.errors)+len(spy.fatals) != 0 {
		t.Errorf("內容一樣時不應該回報失敗，拿到 %q %q", spy.errors, spy.fatals)
	}

	spy = runGolden(t, path, []byte("合計 90 元\n"), false)
	if len(spy.errors) != 1 || len(spy.fatals) != 0 || !strings.Contains(spy.errors[0], path) {
		t.Errorf("內容不同時應該 Errorf 一次，訊息有路徑，拿到 %q %q", spy.errors, spy.fatals)
	}
	if !spy.helper {
		t.Error("Golden 沒有呼叫 tb.Helper()")
	}
}

func TestShrink(t *testing.T) {
	sum := func(xs []int) int {
		s := 0
		for _, x := range xs {
			s += x
		}
		return s
	}
	tests := []struct {
		name  string
		input []int
		fails func([]int) bool
	}{
		{"同時有 3 和 7", []int{1, 3, 5, 7, 9, 3}, func(xs []int) bool {
			return slices.Contains(xs, 3) && slices.Contains(xs, 7)
		}},
		{"總和至少 10", []int{1, 2, 3, 4, 5, 6}, func(xs []int) bool { return sum(xs) >= 10 }},
		{"有重複", []int{4, 8, 15, 16, 23, 42, 15, 8}, func(xs []int) bool {
			return len(slices.Compact(slices.Sorted(slices.Values(xs)))) < len(xs)
		}},
		{"有遞減的相鄰兩個", []int{1, 2, 9, 3, 4, 8, 5}, func(xs []int) bool { return !slices.IsSorted(xs) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := slices.Clone(tt.input)
			got := Shrink(tt.input, tt.fails)
			if !slices.Equal(tt.input, before) {
				t.Fatalf("Shrink 改了 input：%v → %v", before, tt.input)
			}
			if !tt.fails(got) {
				t.Fatalf("Shrink = %v，但它不會失敗", got)
			}
			if !isSubsequence(got, tt.input) {
				t.Errorf("Shrink = %v，不是 %v 依原順序挑出來的元素", got, tt.input)
			}
			for i := range got {
				smaller := slices.Delete(slices.Clone(got), i, i+1)
				if tt.fails(smaller) {
					t.Errorf("Shrink = %v，還可以再拿掉第 %d 個（%v 也會失敗）", got, i, smaller)
					break
				}
			}
		})
	}
}

// isSubsequence 回傳 sub 是不是從 xs 依順序挑出來的元素。
func isSubsequence(sub, xs []int) bool {
	i := 0
	for _, x := range xs {
		if i < len(sub) && sub[i] == x {
			i++
		}
	}
	return i == len(sub)
}

func TestCleanupsRun(t *testing.T) {
	var c Cleanups
	var order []string
	c.Add(func() { order = append(order, "關資料庫") })
	c.Add(func() { order = append(order, "關檔案") })
	c.Add(func() { order = append(order, "刪暫存目錄") })
	c.Run()
	c.Run()
	want := []string{"刪暫存目錄", "關檔案", "關資料庫"}
	if !slices.Equal(order, want) {
		t.Errorf("執行順序 = %q，want %q", order, want)
	}

	c.Add(func() { order = append(order, "再加一個") })
	c.Run()
	if got := order[len(order)-1]; len(order) != 4 || got != "再加一個" {
		t.Errorf("Run 之後再 Add，下一次 Run 只該執行新的那一個，順序 = %q", order)
	}
}

//go:build solution

package testkit

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"unicode"
)

func SubtestNames(names []string) []string {
	seen := map[string]int{}
	out := make([]string, len(names))
	for i, name := range names {
		name = strings.Map(func(r rune) rune {
			if unicode.IsSpace(r) {
				return '_'
			}
			return r
		}, name)
		n := seen[name]
		seen[name] = n + 1
		if n == 0 && name != "" {
			out[i] = name
		} else {
			out[i] = fmt.Sprintf("%s#%02d", name, n)
		}
	}
	return out
}

func DiffLines(got, want string) []string {
	g, w := strings.Split(got, "\n"), strings.Split(want, "\n")
	var out []string
	for i := range max(len(g), len(w)) {
		gl, gok := line(g, i)
		wl, wok := line(w, i)
		if gok == wok && gl == wl {
			continue
		}
		out = append(out, fmt.Sprintf("第 %d 行：got %s，want %s", i+1, show(gl, gok), show(wl, wok)))
	}
	return out
}

func line(lines []string, i int) (string, bool) {
	if i < len(lines) {
		return lines[i], true
	}
	return "", false
}

func show(s string, ok bool) string {
	if !ok {
		return "（沒有這一行）"
	}
	return strconv.Quote(s)
}

func (r *Recorder) Notify(phone, msg string) error {
	if r.FailOn[phone] {
		return ErrUnreachable
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sent = append(r.sent, phone+"："+msg)
	return nil
}

func (r *Recorder) Sent() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return slices.Clone(r.sent)
}

func Golden(tb testing.TB, path string, got []byte, update bool) {
	tb.Helper()
	if update {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			tb.Fatalf("建立 %s 的目錄：%v", path, err)
		}
		if err := os.WriteFile(path, got, 0o644); err != nil {
			tb.Fatalf("寫入 %s：%v", path, err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		tb.Fatalf("%s 不存在，請加上 -update 產生它", path)
	}
	if err != nil {
		tb.Fatalf("讀取 %s：%v（可以加上 -update 重新產生）", path, err)
	}
	if diff := DiffLines(string(got), string(want)); diff != nil {
		tb.Errorf("輸出和 %s 不一樣：\n%s", path, strings.Join(diff, "\n"))
	}
}

// Shrink 一次試著拿掉一個元素，還會失敗就保留拿掉後的版本。
// 一整輪都拿不掉任何元素時，就代表拿掉哪一個都不會失敗了。
func Shrink(input []int, fails func([]int) bool) []int {
	cur := slices.Clone(input)
	for changed := true; changed; {
		changed = false
		for i := 0; i < len(cur); {
			candidate := slices.Delete(slices.Clone(cur), i, i+1)
			if fails(candidate) {
				cur, changed = candidate, true // 同一個 i 現在指向下一個元素
				continue
			}
			i++
		}
	}
	return cur
}

func (c *Cleanups) Run() {
	fns := c.fns
	c.fns = nil
	for _, f := range slices.Backward(fns) {
		f()
	}
}

//go:build solution

package lineup

import (
	"cmp"
	"fmt"
	"iter"
	"strings"
)

func Every(from, to, step int) iter.Seq[int] {
	return func(yield func(int) bool) {
		if step <= 0 {
			return
		}
		for n := from; n <= to; n += step {
			if !yield(n) {
				return
			}
		}
	}
}

// Numbered 用 strings.Lines 拿到每一行（含結尾的換行），再把換行去掉。
func Numbered(s string) iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		i := 0
		for line := range strings.Lines(s) {
			i++
			if !yield(i, strings.TrimSuffix(line, "\n")) {
				return
			}
		}
	}
}

// Dedup 記住上一個值；first 用來處理「還沒有上一個值」的情況。
func Dedup[T comparable](seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		var prev T
		first := true
		for v := range seq {
			if !first && v == prev {
				continue
			}
			first = false
			prev = v
			if !yield(v) {
				return
			}
		}
	}
}

func Segments[T any](seq iter.Seq[T]) iter.Seq2[T, T] {
	return func(yield func(T, T) bool) {
		var prev T
		has := false
		for v := range seq {
			if has && !yield(prev, v) {
				return
			}
			prev, has = v, true
		}
	}
}

// Batch 每交出一批，就重新 make 一個 slice，不重用舊的底層陣列。
func Batch[T any](seq iter.Seq[T], n int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		if n < 1 {
			return
		}
		buf := make([]T, 0, n)
		for v := range seq {
			buf = append(buf, v)
			if len(buf) == n {
				if !yield(buf) {
					return
				}
				buf = make([]T, 0, n)
			}
		}
		if len(buf) > 0 {
			yield(buf)
		}
	}
}

// MergeSorted 把兩邊都轉成 next 函式，每次比較兩邊目前的值，交出小的那個。
func MergeSorted[T cmp.Ordered](a, b iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		nextA, stopA := iter.Pull(a)
		defer stopA()
		nextB, stopB := iter.Pull(b)
		defer stopB()

		va, okA := nextA()
		vb, okB := nextB()
		for okA || okB {
			if okA && (!okB || va <= vb) {
				if !yield(va) {
					return
				}
				va, okA = nextA()
			} else {
				if !yield(vb) {
					return
				}
				vb, okB = nextB()
			}
		}
	}
}

// Seats 檢查 yield 的回傳值：使用端 break 之後就停。
func Seats(row string, n int) iter.Seq[string] {
	return func(yield func(string) bool) {
		for i := 1; i <= n; i++ {
			if !yield(fmt.Sprintf("%s%d", row, i)) {
				return
			}
		}
	}
}

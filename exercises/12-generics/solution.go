//go:build solution

package menu

import (
	"cmp"
	"maps"
	"slices"
)

// Filter 建立新的 slice 來放結果，不重用 s 的底層陣列，s 才不會被改到。
func Filter[T any](s []T, keep func(T) bool) []T {
	var out []T
	for _, v := range s {
		if keep(v) {
			out = append(out, v)
		}
	}
	return out
}

// Average 把每個值轉成 float64 再加總：N 的型別集合裡都是數字，所以可以轉換。
func Average[N Number](s []N) (float64, bool) {
	if len(s) == 0 {
		return 0, false
	}
	total := 0.0
	for _, v := range s {
		total += float64(v)
	}
	return total / float64(len(s)), true
}

func GroupBy[T any, K comparable](s []T, key func(T) K) map[K][]T {
	out := map[K][]T{}
	for _, v := range s {
		k := key(v)
		out[k] = append(out[k], v)
	}
	return out
}

// Add 第一次用到時才建立 map，所以零值的 Tally 也能用。
func (t *Tally[K]) Add(k K) {
	if t.counts == nil {
		t.counts = map[K]int{}
	}
	t.counts[k]++
}

func (t *Tally[K]) Count(k K) int {
	return t.counts[k]
}

// Top 先依次數由多到少排，次數一樣時用 cmp.Compare 比 key；K 是 cmp.Ordered 才能這樣比。
func (t *Tally[K]) Top(n int) []K {
	keys := slices.Collect(maps.Keys(t.counts))
	slices.SortFunc(keys, func(a, b K) int {
		if c := cmp.Compare(t.counts[b], t.counts[a]); c != 0 {
			return c
		}
		return cmp.Compare(a, b)
	})
	return keys[:min(n, len(keys))]
}

// Then 是 generic method：U 由 f 的回傳型別推論出來。
func (r Result[T]) Then[U any](f func(T) (U, error)) Result[U] {
	if r.Err != nil {
		return Result[U]{Err: r.Err}
	}
	v, err := f(r.Val)
	return Result[U]{Val: v, Err: err}
}

// Memo 回傳的 closure 抓住了 cache，每次呼叫都共用同一個 map。
func Memo[K comparable, V any](f func(K) V) func(K) V {
	cache := map[K]V{}
	return func(k K) V {
		if v, ok := cache[k]; ok {
			return v
		}
		v := f(k)
		cache[k] = v
		return v
	}
}

// Max 從第一個元素開始比，而不是從 T 的零值開始。
func Max[T cmp.Ordered](s []T) (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}
	best := s[0]
	for _, v := range s[1:] {
		if v > best {
			best = v
		}
	}
	return best, true
}

//go:build !solution

package menu

import "cmp"

// Filter 回傳 s 裡 keep 回傳 true 的元素，順序和原本一樣。不可以修改 s。
func Filter[T any](s []T, keep func(T) bool) []T {
	// TODO: 實作
	return nil
}

// Average 回傳 s 的平均值。s 是空的時回傳 0, false。
func Average[N Number](s []N) (float64, bool) {
	// TODO: 實作
	return 0, false
}

// GroupBy 用 key 算出每個元素的分組，把同一組的元素放在同一個 slice 裡，
// 每一組裡的順序和原本一樣。
func GroupBy[T any, K comparable](s []T, key func(T) K) map[K][]T {
	// TODO: 實作
	return nil
}

// Add 把 k 出現的次數加 1。
func (t *Tally[K]) Add(k K) {
	// TODO: 實作
}

// Count 回傳 k 出現的次數。
func (t *Tally[K]) Count(k K) int {
	// TODO: 實作
	return 0
}

// Top 回傳出現最多次的前 n 個 key：次數多的在前，次數一樣時 key 小的在前。
// n 比 key 的數量還多時，全部回傳。
func (t *Tally[K]) Top(n int) []K {
	// TODO: 實作
	return nil
}

// Then 在 r 成功時呼叫 f(r.Val)，把 f 的結果包成 Result[U]；
// r 已經失敗時不呼叫 f，直接把 r.Err 帶到新的 Result[U]。
func (r Result[T]) Then[U any](f func(T) (U, error)) Result[U] {
	// TODO: 實作
	return Result[U]{}
}

// Memo 回傳一個和 f 結果一樣的函式，但會記住算過的結果：
// 同一個參數第二次以後的呼叫，直接回傳記住的值，不再呼叫 f。
func Memo[K comparable, V any](f func(K) V) func(K) V {
	// TODO: 實作
	return f
}

// Max 回傳 s 裡最大的值；s 是空的時回傳零值和 false。
//
// 找錯：Max([]int{-3, -1, -7}) 回傳 0，可是 0 根本不在 slice 裡。
func Max[T cmp.Ordered](s []T) (T, bool) {
	var best T
	for _, v := range s {
		if v > best {
			best = v
		}
	}
	return best, len(s) > 0
}

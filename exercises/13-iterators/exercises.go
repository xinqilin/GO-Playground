//go:build !solution

package lineup

import (
	"cmp"
	"fmt"
	"iter"
)

// Every 依序交出 from、from+step、from+2*step……，直到超過 to 為止（to 本身也算）。
// 例如 Every(360, 420, 15) 交出 360 375 390 405 420。step 小於等於 0 時什麼都不交出。
func Every(from, to, step int) iter.Seq[int] {
	// TODO: 實作
	return func(yield func(int) bool) {}
}

// Numbered 一行一行交出 s 的內容：第一個值是行號（從 1 開始），第二個值是那一行，
// 不含結尾的 "\n"。最後一行沒有換行也要交出；s 是空字串時什麼都不交出。
func Numbered(s string) iter.Seq2[int, string] {
	// TODO: 實作
	return func(yield func(int, string) bool) {}
}

// Dedup 交出 seq 的值，但跳過和前一個值相同的值。
// 例如 17 17 18 18 18 17 會變成 17 18 17。
func Dedup[T comparable](seq iter.Seq[T]) iter.Seq[T] {
	// TODO: 實作
	return func(yield func(T) bool) {}
}

// Segments 把一串站名變成一段一段的區間：交出每兩個相鄰的值。
// 例如 台北 中山 雙連 會交出 (台北, 中山) 和 (中山, 雙連)。
// 少於兩個值時什麼都不交出。
func Segments[T any](seq iter.Seq[T]) iter.Seq2[T, T] {
	// TODO: 實作
	return func(yield func(T, T) bool) {}
}

// Batch 把 seq 的值每 n 個包成一個 slice 交出去，最後不滿 n 個的也要交出。
// 每次交出的都要是新的 slice，因為使用端可能把它留下來。n 小於 1 時什麼都不交出。
func Batch[T any](seq iter.Seq[T], n int) iter.Seq[[]T] {
	// TODO: 實作
	return func(yield func([]T) bool) {}
}

// MergeSorted 把兩個已經由小到大排好的疊代器合併成一個，結果也由小到大；
// 一樣大時先交出 a 的。用 iter.Pull 一次從兩邊各拿一個來比，而且一定要呼叫 stop。
func MergeSorted[T cmp.Ordered](a, b iter.Seq[T]) iter.Seq[T] {
	// TODO: 實作
	return func(yield func(T) bool) {}
}

// Seats 依序交出 row 這一排的座位編號，例如 Seats("A", 3) 交出 A1 A2 A3。
//
// 找錯：全部取完沒問題，但使用端只取兩個就 break 時，程式 panic。
func Seats(row string, n int) iter.Seq[string] {
	return func(yield func(string) bool) {
		for i := 1; i <= n; i++ {
			yield(fmt.Sprintf("%s%d", row, i))
		}
	}
}

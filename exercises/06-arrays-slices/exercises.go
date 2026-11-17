//go:build !solution

package slicing

// Chunk 把訂單依序分批，每批 size 筆，最後一批可以不滿。
// size 小於 1 或沒有訂單時回傳 nil。
func Chunk(orders []string, size int) [][]string {
	// TODO: 實作
	return nil
}

// Discounted 回傳每個價格打 pct% 折扣後的新 slice（無條件捨去），prices 本身不能被修改。
//
// 找錯：回傳的價格是對的，但呼叫端原本的價格也被改掉了。
func Discounted(prices []int, pct int) []int {
	out := prices
	for i := range out {
		out[i] = out[i] * (100 - pct) / 100
	}
	return out
}

// Compact 移除相鄰的重複條碼，只保留連續相同條碼中的第一個，順序不變。
// 要原地處理：結果沿用 codes 的底層陣列，不配置新的陣列。
func Compact(codes []string) []string {
	// TODO: 實作
	return nil
}

// SplitAt 把訂單分成前 n 筆和其餘的。n 小於 0 時當成 0，大於長度時當成長度。
// 呼叫端之後可能會對 head 做 append，這不能影響 tail，也不能改到 orders。
//
// 找錯：對 head append 之後，tail 的第一筆被蓋掉了。
func SplitAt(orders []int, n int) (head, tail []int) {
	n = max(0, min(n, len(orders)))
	return orders[:n], orders[n:]
}

// Rotate 把值班表原地往左轉 k 格，例如 [A B C D] 轉 1 格變成 [B C D A]。
// k 不會是負數，但可能大於人數。空的值班表什麼都不做。
func Rotate(staff []string, k int) {
	// TODO: 實作
}

// BestStreak 找出連續 k 天營業額總和最高的區間，回傳起始索引和總和；
// 有好幾個區間一樣高時取最早的。k 小於 1 或大於天數時回傳 -1, 0。
func BestStreak(sales []int, k int) (start, total int) {
	// TODO: 實作
	return 0, 0
}

// InsertSorted 把 v 插入由小到大排好的 s，回傳的結果仍然排好。
// s 的容量足夠時要原地插入（結果沿用 s 的底層陣列）；容量不夠時才配置新的陣列，而且不能改到 s。
func InsertSorted(s []int, v int) []int {
	// TODO: 實作
	return s
}

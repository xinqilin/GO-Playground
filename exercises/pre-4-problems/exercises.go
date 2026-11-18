//go:build !solution

package problems

// CartTotal 回傳購物車的總價：prices[i] 是第 i 項的單價，qtys[i] 是它的數量。
// 兩份清單一樣長。例如單價 [85 90]、數量 [2 1]，總價是 260。
func CartTotal(prices, qtys []int) int {
	// TODO: 實作
	return 0
}

// Contains 回傳 menu 裡有沒有 dish 這道菜。
func Contains(menu []string, dish string) bool {
	// TODO: 實作
	return false
}

// BiggestJump 回傳相鄰兩天之間最大的漲幅（後一天減前一天）。
// 一次都沒漲、或不到兩天的資料時回傳 0。
// 例如 [100 103 98 106 107] 的漲幅是 3、-5、8、1，最大是 8。
func BiggestJump(prices []int) int {
	// TODO: 實作
	return 0
}

// Above 回傳一份新的清單，只留下 prices 裡大於 limit 的價格，順序不變。
// 例如 Above([85 120 90 150], 100) 是 [120 150]。
func Above(prices []int, limit int) []int {
	// TODO: 實作
	return nil
}

// MostVoted 回傳得票最多的選項。票數一樣多時，回傳最早出現在 votes 裡的那一個。
// votes 是空的時回傳 ""。例如 [拉麵 咖哩 拉麵 便當 咖哩] 是 拉麵（2 票，比咖哩早出現）。
func MostVoted(votes []string) string {
	// TODO: 實作
	return ""
}

// CheapestPrice 回傳 prices 裡最低的價格。prices 至少有一個價格，而且都大於 0。
//
// 找錯：CheapestPrice([85 90 75]) 應該是 75，結果是 0。
func CheapestPrice(prices []int) int {
	best := 0
	for _, p := range prices {
		if p < best {
			best = p
		}
	}
	return best
}

// LowStock 回傳庫存量低於 limit 的品項有幾個。
//
// 找錯：LowStock([3 12 0 8], 5) 應該是 2，結果卻是 0。
func LowStock(stock []int, limit int) int {
	count := 0
	for _, s := range stock {
		count := 0
		if s < limit {
			count++
		}
	}
	return count
}

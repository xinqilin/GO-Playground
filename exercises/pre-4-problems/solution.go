//go:build solution

package problems

func CartTotal(prices, qtys []int) int {
	total := 0
	for i := range prices {
		total += prices[i] * qtys[i]
	}
	return total
}

// Contains 找到就立刻 return true；整份清單找完都沒有，才回傳 false。
func Contains(menu []string, dish string) bool {
	for _, d := range menu {
		if d == dish {
			return true
		}
	}
	return false
}

// BiggestJump 從索引 1 開始，每一天都和前一天比；best 從 0 開始，正好代表「沒有漲」。
func BiggestJump(prices []int) int {
	best := 0
	for i := 1; i < len(prices); i++ {
		jump := prices[i] - prices[i-1]
		if jump > best {
			best = jump
		}
	}
	return best
}

func Above(prices []int, limit int) []int {
	result := []int{}
	for _, p := range prices {
		if p > limit {
			result = append(result, p)
		}
	}
	return result
}

// MostVoted 對每一張票，用內層迴圈數出它的票數；只有票數「大於」目前最多時才換，
// 所以票數一樣時會留下比較早出現的那一個。
func MostVoted(votes []string) string {
	best := ""
	bestCount := 0
	for _, v := range votes {
		count := 0
		for _, w := range votes {
			if w == v {
				count++
			}
		}
		if count > bestCount {
			best = v
			bestCount = count
		}
	}
	return best
}

// CheapestPrice 要從第一個價格開始比；從 0 開始的話，沒有任何價格會比 0 小。
func CheapestPrice(prices []int) int {
	best := prices[0]
	for _, p := range prices {
		if p < best {
			best = p
		}
	}
	return best
}

// LowStock 的計數器只能在迴圈之前建立一次。迴圈裡的 count := 0 建立的是另一個同名的新變數，
// 每圈都從 0 開始，而且迴圈結束就消失了，外面的 count 從來沒有被加過。
func LowStock(stock []int, limit int) int {
	count := 0
	for _, s := range stock {
		if s < limit {
			count++
		}
	}
	return count
}

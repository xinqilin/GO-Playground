//go:build solution

package slicing

func Chunk(orders []string, size int) [][]string {
	if size < 1 || len(orders) == 0 {
		return nil
	}
	var batches [][]string
	for start := 0; start < len(orders); start += size {
		end := min(start+size, len(orders))
		batches = append(batches, orders[start:end])
	}
	return batches
}

// Discounted 要配置新的 slice；out := prices 只複製了 slice header，兩者看向同一個陣列。
func Discounted(prices []int, pct int) []int {
	out := make([]int, len(prices))
	for i, p := range prices {
		out[i] = p * (100 - pct) / 100
	}
	return out
}

func Compact(codes []string) []string {
	if len(codes) == 0 {
		return codes
	}
	kept := codes[:1]
	for _, c := range codes[1:] {
		if c != kept[len(kept)-1] {
			kept = append(kept, c)
		}
	}
	return kept
}

// SplitAt 用完整切片運算式把 head 的容量限制在 n，之後的 append 就會換到新陣列。
func SplitAt(orders []int, n int) (head, tail []int) {
	n = max(0, min(n, len(orders)))
	return orders[:n:n], orders[n:]
}

// Rotate 先把前 k 格存起來，後面的往前搬，再把存起來的放到最後。
func Rotate(staff []string, k int) {
	if len(staff) == 0 {
		return
	}
	k %= len(staff)
	front := append([]string(nil), staff[:k]...)
	copy(staff, staff[k:])
	copy(staff[len(staff)-k:], front)
}

// BestStreak 用滑動視窗：每往右移一天，加上新的一天、減掉離開的那天，不必每次重新加總。
func BestStreak(sales []int, k int) (start, total int) {
	if k < 1 || k > len(sales) {
		return -1, 0
	}
	sum := 0
	for _, s := range sales[:k] {
		sum += s
	}
	start, total = 0, sum
	for i := k; i < len(sales); i++ {
		sum += sales[i] - sales[i-k]
		if sum > total {
			start, total = i-k+1, sum
		}
	}
	return start, total
}

// InsertSorted 先用 append 多要一格（容量不夠時 append 會換新陣列，s 不受影響），
// 再把插入點之後的元素往後搬一格。
func InsertSorted(s []int, v int) []int {
	i := len(s)
	for i > 0 && s[i-1] > v {
		i--
	}
	s = append(s, 0)
	copy(s[i+1:], s[i:])
	s[i] = v
	return s
}

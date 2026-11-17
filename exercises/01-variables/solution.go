//go:build solution

package variables

import "fmt"

func SplitBill(total, people int) (each, rest int) {
	if people < 1 {
		return 0, 0
	}
	return total / people, total % people
}

func Tables(guests, perTable int) int {
	if guests < 1 || perTable < 1 {
		return 0
	}
	return (guests + perTable - 1) / perTable
}

// Average 先把總和轉成 float64 再除，小數才不會在整數除法時被捨去。
func Average(a, b, c int) float64 {
	return float64(a+b+c) / 3
}

const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
)

func FormatSize(bytes int64) string {
	if bytes < 0 {
		bytes = 0
	}
	switch {
	case bytes < KB:
		return fmt.Sprintf("%d B", bytes)
	case bytes < MB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/KB)
	case bytes < GB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/MB)
	default:
		return fmt.Sprintf("%.1f GB", float64(bytes)/GB)
	}
}

func Grant(p, q Perm) Perm {
	return p | q
}

func Revoke(p, q Perm) Perm {
	return p &^ q
}

// Has 要檢查 q 的每一個位元都在 p 裡，所以是 p&q == q，不是 p&q != 0。
func Has(p, q Perm) bool {
	return p&q == q
}

// AddStock 利用溢位會繞回來的特性：相加後反而變小，就表示超過了。
func AddStock(stock, n uint8) (uint8, bool) {
	sum := stock + n
	if sum < stock {
		return 255, false
	}
	return sum, true
}

// Ticket 先換成從 0 開始的位置再取餘數。Go 的 % 結果正負號跟著被除數，
// 往回數時可能得到負數，所以再加 999 並取一次餘數，把它拉回 0 到 998。
func Ticket(current, delta int) int {
	pos := (current - 1 + delta) % 999
	pos = (pos + 999) % 999
	return pos + 1
}

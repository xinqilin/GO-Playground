//go:build !solution

package controlflow

// TicketPrice 依搭乘公里數回傳票價（元）。
//
//	1–5 公里：20；6–10：25；11–20：35；21–40：50；超過 40：65
//	0 或負數不是合法的距離，回傳 0。
func TicketPrice(km int) int {
	// TODO: 實作
	return 0
}

// Restock 幫每一項庫存補 n 個。
//
// 找錯：這個版本跑完之後庫存完全沒變，請修好它。
func Restock(items []Item, n int) {
	for _, it := range items {
		it.Stock += n
	}
}

// Table 回傳 n×n 的乘法表。第 i 列（從 1 開始）依序是 i×1、i×2…i×n，
// 每格寫成 "i×j=積"，同一列的格子用 tab 分開，列和列之間用換行分開，最後一列後面沒有換行。
// n 小於 1 時回傳空字串。
//
//	Table(2) == "1×1=1\t1×2=2\n2×1=2\t2×2=4"
func Table(n int) string {
	// TODO: 實作
	return ""
}

// InvoicePrize 用簡化的發票規則對獎，winning 和 ticket 都是 8 位數字的字串。
// 從末 8 碼開始往下比，看末幾碼完全相同，回傳最高的獎金：
//
//	末 8 碼 200000、末 7 碼 40000、末 6 碼 10000、末 5 碼 4000、末 4 碼 1000、末 3 碼 200
//
// 連末 3 碼都不同就回傳 0。
func InvoicePrize(winning, ticket string) int {
	// TODO: 實作
	return 0
}

// CountKinds 統計字串裡四種字元的數量：
//
//	han：漢字（unicode.Han）
//	letter：ASCII 英文字母 a–z、A–Z
//	digit：ASCII 數字 0–9
//	other：其他所有字元，包含空白、標點和全形英數
func CountKinds(s string) (han, letter, digit, other int) {
	// TODO: 實作
	return 0, 0, 0, 0
}

// TotalUntilNegative 從頭加總 nums，遇到第一個負數就停止（負數本身不加）。
//
// 找錯：遇到負數之後還是一直加下去，請修好它。
func TotalUntilNegative(nums []int) int {
	total := 0
	for _, n := range nums {
		switch {
		case n < 0:
			break
		default:
			total += n
		}
	}
	return total
}

// Spiral 回傳 n×n 的螺旋矩陣：從左上角的 1 開始，順時針往內繞，填到 n×n 為止。
// n 小於 1 時回傳長度為 0 的 slice。
//
//	Spiral(3) == [][]int{{1, 2, 3}, {8, 9, 4}, {7, 6, 5}}
func Spiral(n int) [][]int {
	// TODO: 實作
	return nil
}

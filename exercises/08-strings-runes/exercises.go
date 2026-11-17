//go:build !solution

package text

// MaskName 遮住姓名中間的字：保留第一個和最後一個字，中間每個字換成 *。
// 只有兩個字時保留第一個字，第二個字換成 *；一個字或空字串原樣回傳。
func MaskName(name string) string {
	// TODO: 實作
	return ""
}

// NormalizePhone 整理手機號碼：去掉空白和 -，開頭的 +886 換成 0。
// 整理後必須是 09 開頭的 10 位數字，才回傳 (號碼, true)；否則回傳 ("", false)。
func NormalizePhone(s string) (string, bool) {
	// TODO: 實作
	return "", false
}

// ParseOrder 解析「品名x數量」用逗號分隔的訂單，例如 "排骨x2, 雞腿x1"，同一個品名的數量要加總。
// 每一項前後的空白要去掉。品名不能是空的，數量必須是正整數，否則回傳 nil 和一個 error。
// 空字串代表沒有點餐，回傳空的 map 和 nil。
func ParseOrder(s string) (map[string]int, error) {
	// TODO: 實作
	return nil, nil
}

// Truncate 讓商品名稱最多顯示 n 個字（n 至少是 1）：超過時保留前 n 個字，再加上 "…"。
//
// 找錯：英文名稱沒問題，中文名稱卻變成亂碼。
func Truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

// IsPalindrome 回報 s 正著讀和倒著讀是否一樣（以字為單位，不忽略任何字元）。
//
// 找錯：英文的 "level" 是對的，中文的回文卻都回傳 false。
func IsPalindrome(s string) bool {
	for i := 0; i < len(s)/2; i++ {
		if s[i] != s[len(s)-1-i] {
			return false
		}
	}
	return true
}

// FormatMenu 把菜單排成兩欄，每個品項一行，行與行之間用 "\n" 連接，最後一行後面沒有換行。
// 品名欄的寬度是所有品名中最寬的顯示寬度（用 RuneWidth 計算），不足的補空白；
// 接著是兩個空白，然後是價格，價格靠右對齊到最長的價格位數。沒有品項時回傳 ""。
func FormatMenu(items []MenuItem) string {
	// TODO: 實作
	return ""
}

//go:build !solution

package hello

// Greeting 回傳 "哈囉，<name>！"；name 是空字串時回傳 "哈囉，陌生人！"。
func Greeting(name string) string {
	// TODO: 實作
	return ""
}

// BentoTotal 回傳買 qty 個便當的總價：一個 85 元，買 10 個以上每個算 80 元。
// qty 小於 1 時回傳 0。
func BentoTotal(qty int) int {
	// TODO: 實作
	return 0
}

// Clock 把秒數轉成 "時:分:秒"，每一段都是兩位數，例如 3725 秒是 "01:02:05"。
// 負數當成 0。時數超過 99 也照實寫出，例如 360000 秒是 "100:00:00"。
func Clock(seconds int) string {
	// TODO: 實作
	return ""
}

// Box 用 ASCII 字元把 text 框起來（text 只會有 ASCII 字元），三列之間用換行分開：
//
//	+----+
//	| Go |
//	+----+
//
// 框的內寬是 len(text)+2，文字左右各留一個空白。
func Box(text string) string {
	// TODO: 實作
	return ""
}

// GoVersion 讀取 go.mod 裡的 go 那一行，例如 "go 1.27.4" 或 "go 1.27"，回傳主版號和次版號。
// 格式不對時（不是 "go " 開頭、少了次版號、版號不是數字）回傳 0, 0, false。
func GoVersion(line string) (major, minor int, ok bool) {
	// TODO: 實作
	return 0, 0, false
}

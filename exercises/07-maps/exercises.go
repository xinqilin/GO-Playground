//go:build !solution

package mapping

import (
	"fmt"
	"strings"
)

// Tally 統計每個品項被點了幾次，並回傳被點最多次的品項；次數相同時取名稱排序較前面的。
// 沒有任何訂單時 top 是 ""。
func Tally(orders []string) (counts map[string]int, top string) {
	// TODO: 實作
	return nil, ""
}

// Invert 把「代碼 → 品名」反轉成「品名 → 代碼」。
// 如果有兩個代碼對到同一個品名，無法反轉，回傳 nil, false。
func Invert(codes map[string]string) (map[string]string, bool) {
	// TODO: 實作
	return nil, false
}

// GroupByLine 依路線把車站分組，每條路線的站名由小到大排序。
func GroupByLine(stations []Station) map[string][]string {
	// TODO: 實作
	return nil
}

// Restock 替 name 補 n 個庫存，回傳補貨後的數量。name 不存在時新增一筆，價格是 0。
//
// 找錯：回傳的數量是對的，但 stock 裡的庫存沒有變。
func Restock(stock map[string]Item, name string, n int) int {
	item := stock[name]
	item.Qty += n
	return item.Qty
}

// Summary 依品名排序，把價目表寫成 "品名=價格"，再用 ", " 連起來，例如 "排骨=85, 雞腿=90"。
//
// 找錯：內容是對的，但每次呼叫的順序都可能不一樣。
func Summary(menu map[string]int) string {
	var parts []string
	for name, price := range menu {
		parts = append(parts, fmt.Sprintf("%s=%d", name, price))
	}
	return strings.Join(parts, ", ")
}

// Common 回傳同時出現在 a 和 b 裡的品項，由小到大排序，不重複。沒有共同的品項時回傳 nil。
func Common(a, b []string) []string {
	// TODO: 實作
	return nil
}

// Hops 計算從 from 到 to 最少要坐幾站。links 記錄每一站可以直接到達的下一站，
// 雙向的路線兩邊都會記。from 和 to 相同時是 0；到不了時回傳 -1。
func Hops(links map[string][]string, from, to string) int {
	// TODO: 實作
	return 0
}

// Package oldstyle 用比較舊的寫法整理菜單，用來看 go fix 會建議怎麼改。
package oldstyle

import (
	"sort"
	"strings"
	"sync"
)

// Lookup 在 menu 裡找菜名，找不到回傳 nil。
func Lookup(menu map[string]interface{}, name string) interface{} {
	return menu[name]
}

// Sorted 回傳排序好的價格。
func Sorted(prices []int) []int {
	sort.Slice(prices, func(i, j int) bool { return prices[i] < prices[j] })
	return prices
}

// HasDish 回傳 menu 裡有沒有 name。
func HasDish(menu []string, name string) bool {
	for _, d := range menu {
		if d == name {
			return true
		}
	}
	return false
}

// Note 取出「備註：」後面的文字。
func Note(line string) string {
	if strings.HasPrefix(line, "備註：") {
		return strings.TrimPrefix(line, "備註：")
	}
	return ""
}

// CookAll 同時處理每一張單。
func CookAll(orders []string, cook func(string)) {
	var wg sync.WaitGroup
	for _, o := range orders {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cook(o)
		}()
	}
	wg.Wait()
}

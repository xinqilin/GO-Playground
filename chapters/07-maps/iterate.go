package main

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
)

// rangeMap 示範走訪 map。順序每次都可能不同，所以這裡只做和順序無關的加總。
func rangeMap() {
	menu := map[string]int{"排骨": 85, "雞腿": 90, "素食": 70}
	total := 0
	for _, price := range menu {
		total += price
	}
	fmt.Println(total)
}

// sortedRange 示範需要固定順序時，先取出所有 key 排序，再依序讀值。
func sortedRange() {
	menu := map[string]int{"排骨": 85, "雞腿": 90, "素食": 70, "魚排": 100}
	for _, name := range slices.Sorted(maps.Keys(menu)) {
		fmt.Println(name, menu[name])
	}

	// Go 1.23 以前的寫法：自己收集 key 再排序
	names := make([]string, 0, len(menu))
	for name := range menu {
		names = append(names, name)
	}
	slices.Sort(names)
	fmt.Println(names)
}

// sortByValue 示範依照值排序：先取出 key，排序時比較 key 對應的值；值一樣時再比 key，讓結果固定。
func sortByValue() {
	sold := map[string]int{"排骨": 12, "雞腿": 30, "素食": 12, "魚排": 7}
	names := slices.Collect(maps.Keys(sold))
	slices.SortFunc(names, func(a, b string) int {
		return cmp.Or(
			cmp.Compare(sold[b], sold[a]), // 賣得多的排前面
			cmp.Compare(a, b),
		)
	})
	for _, name := range names {
		fmt.Println(name, sold[name])
	}
}

// deleteWhileRange 示範在 range 裡刪除目前的 key 是安全的，以及 maps.DeleteFunc 和 clear。
func deleteWhileRange() {
	stock := map[string]int{"排骨": 0, "雞腿": 3, "素食": 0, "魚排": 5}
	for name, n := range stock {
		if n == 0 {
			delete(stock, name)
		}
	}
	fmt.Println(stock)

	maps.DeleteFunc(stock, func(name string, n int) bool { return n < 4 })
	fmt.Println(stock)

	clear(stock)
	fmt.Println(stock, len(stock))
}

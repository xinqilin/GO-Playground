package main

import (
	"fmt"
	"maps"
	"slices"
)

// mapAsSet 示範用 map 當集合：只在乎 key 在不在。
func mapAsSet() {
	soldOut := map[string]struct{}{} // struct{} 不佔空間，只用來表示「有這個 key」
	soldOut["雞腿"] = struct{}{}
	soldOut["魚排"] = struct{}{}
	_, out := soldOut["雞腿"]
	fmt.Println(out, len(soldOut))
	fmt.Println(slices.Sorted(maps.Keys(soldOut)))

	seen := make(map[string]bool) // 用 bool 寫起來比較直覺：不存在的 key 讀到 false
	var dup []string
	for _, t := range []string{"A01", "A02", "A01", "A03", "A02"} {
		if seen[t] {
			dup = append(dup, t)
		}
		seen[t] = true
	}
	fmt.Println(dup)
}

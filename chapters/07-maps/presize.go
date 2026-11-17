package main

import "fmt"

// buildIndex 建立「名稱 → 位置」的索引。事先知道 key 的數量，就用 make 的第二個參數告訴 map，
// 它可以一次配置好空間，不必邊放邊擴充。
func buildIndex(names []string) map[string]int {
	idx := make(map[string]int, len(names))
	for i, name := range names {
		idx[name] = i
	}
	return idx
}

// indexByName 示範使用 buildIndex 查位置。
func indexByName() {
	idx := buildIndex([]string{"排骨", "雞腿", "素食"})
	fmt.Println(idx["素食"], len(idx))
}

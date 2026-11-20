package main

import (
	"fmt"
	"maps"
)

// raisePrices 收到的 map 和呼叫端的 map 是同一份資料。
func raisePrices(menu map[string]int, n int) {
	for name := range menu {
		menu[name] += n
	}
}

// mapIsReference 示範 map 變數指向同一份資料：指定和傳參數都不會複製內容。要複本就用 maps.Clone。
func mapIsReference() {
	menu := map[string]int{"排骨": 85, "雞腿": 90}
	alias := menu
	alias["素食"] = 70
	fmt.Println(len(menu))

	raisePrices(menu, 5)
	fmt.Println(menu)

	snapshot := maps.Clone(menu)
	snapshot["排骨"] = 0
	fmt.Println(menu["排骨"], snapshot["排骨"])
	fmt.Println(maps.Equal(menu, snapshot))
}

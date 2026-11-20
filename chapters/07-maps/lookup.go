package main

import "fmt"

// commaOK 示範找不到 key 時拿到零值，以及用第二個回傳值分辨「不存在」和「值剛好是零」。
func commaOK() {
	stock := map[string]int{"排骨": 0, "雞腿": 3}
	fmt.Println(stock["排骨"], stock["魚排"]) // 都是 0，分不出來

	n, ok := stock["排骨"]
	fmt.Println(n, ok)
	n, ok = stock["魚排"]
	fmt.Println(n, ok)

	if _, ok := stock["雞腿"]; ok {
		fmt.Println("有賣雞腿")
	}
}

// countOrders 示範利用零值計數：key 不存在時讀到 0，加一之後再存回去。
func countOrders() {
	orders := []string{"排骨", "雞腿", "排骨", "素食", "排骨"}
	counts := make(map[string]int)
	for _, o := range orders {
		counts[o]++
	}
	fmt.Println(counts)
}

// groupOrders 示範 map 的值是 slice：key 不存在時讀到 nil slice，直接 append 就好。
func groupOrders() {
	type order struct {
		table int
		item  string
	}
	orders := []order{{1, "排骨"}, {3, "雞腿"}, {1, "素食"}, {2, "排骨"}, {3, "滷蛋"}}
	byTable := make(map[int][]string)
	for _, o := range orders {
		byTable[o.table] = append(byTable[o.table], o.item)
	}
	fmt.Println(byTable)
	fmt.Println(byTable[1], len(byTable[9]))
}

package main

import "fmt"

// mapBasics 示範用字面值建立 map，以及新增、覆寫、刪除和長度。
func mapBasics() {
	menu := map[string]int{
		"排骨": 85,
		"雞腿": 90,
	}
	menu["素食"] = 70 // 新增
	menu["排骨"] = 95 // key 已經存在：覆寫
	delete(menu, "雞腿")
	delete(menu, "魚排") // 刪除不存在的 key，什麼都不會發生
	fmt.Println(menu, len(menu))
	fmt.Println(menu["排骨"])
}

// makeMap 示範用 make 建立一個空的 map，之後再放資料。
func makeMap() {
	stock := make(map[string]int)
	stock["排骨"] = 10
	stock["雞腿"] = 6
	fmt.Println(stock, len(stock))
}

// nilMap 示範 nil map：可以讀、可以取長度，但寫入會 panic。這裡用 recover 接住才能印出訊息。
func nilMap() {
	var stock map[string]int
	fmt.Println(stock == nil, len(stock), stock["排骨"])
	defer func() {
		fmt.Println("panic:", recover())
	}()
	stock["排骨"] = 10
}

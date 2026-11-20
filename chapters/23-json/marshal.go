package main

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
)

// MenuItem 是菜單上的一道菜。只有大寫開頭（匯出）的欄位會變成 JSON，
// 反引號裡的 json:"name" 是 struct tag，指定 JSON 裡的名稱。
type MenuItem struct {
	Name  string   `json:"name"`
	Price int      `json:"price"`
	Tags  []string `json:"tags"`
	cost  int      // 成本：小寫開頭，不會出現在 JSON 裡
}

// marshalDemo 把 Go 的值編碼成 JSON。
func marshalDemo() {
	item := MenuItem{Name: "排骨飯", Price: 90, Tags: []string{"招牌", "炸物"}, cost: 55}
	b, err := json.Marshal(item)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

	b, err = json.Marshal(MenuItem{Name: "白飯", Price: 10}) // Tags 是 nil
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

	b, err = json.Marshal(item, jsontext.WithIndent("  ")) // 給人看的時候加上縮排
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

// unmarshalDemo 把 JSON 解碼到 Go 的值：要傳指標，Unmarshal 才改得到。
func unmarshalDemo() {
	data := []byte(`{"name":"雞腿飯","price":100,"tags":["招牌"],"calories":780}`)
	var item MenuItem
	if err := json.Unmarshal(data, &item); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", item)

	stock := map[string]int{"雞腿": 12, "排骨": 0, "滷蛋": 30}
	b, err := json.Marshal(stock, json.Deterministic(true)) // map 的順序不固定，要固定就加這個選項
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

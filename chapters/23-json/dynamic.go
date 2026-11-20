package main

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
)

// lookupDemo 不知道結構時解碼到 map[string]any，每往下一層都要做一次型別斷言。
func lookupDemo() {
	data := []byte(`{"shop":{"name":"阿嬤便當","rating":4.7},"open":true}`)
	var v map[string]any
	if err := json.Unmarshal(data, &v); err != nil {
		fmt.Println(err)
		return
	}

	shop, ok1 := v["shop"].(map[string]any)
	name, ok2 := shop["name"].(string)
	rating, ok3 := shop["rating"].(float64)
	fmt.Println(name, rating, ok1 && ok2 && ok3)

	_, ok := v["closed_on"].(string)
	fmt.Println("沒有的欄位：", ok)
}

// Event 是外送平台送來的通知：type 決定 data 長什麼樣子。
type Event struct {
	Type string         `json:"type"`
	Data jsontext.Value `json:"data"` // 先原封不動留著，看了 type 再決定怎麼解
}

// OrderPlaced 是 type 為 order.placed 時的 data。
type OrderPlaced struct {
	OrderID string `json:"order_id"`
	Total   int    `json:"total"`
}

// RiderAssigned 是 type 為 rider.assigned 時的 data。
type RiderAssigned struct {
	OrderID string `json:"order_id"`
	Rider   string `json:"rider"`
	ETA     int    `json:"eta_minutes"`
}

// describe 解兩次：第一次只看外層的 type，第二次才把 data 解成對應的型別。
func describe(msg []byte) (string, error) {
	var ev Event
	if err := json.Unmarshal(msg, &ev); err != nil {
		return "", err
	}
	switch ev.Type {
	case "order.placed":
		var p OrderPlaced
		if err := json.Unmarshal(ev.Data, &p); err != nil {
			return "", err
		}
		return fmt.Sprintf("新訂單 %s，共 %d 元", p.OrderID, p.Total), nil
	case "rider.assigned":
		var r RiderAssigned
		if err := json.Unmarshal(ev.Data, &r); err != nil {
			return "", err
		}
		return fmt.Sprintf("%s 由 %s 外送，約 %d 分鐘", r.OrderID, r.Rider, r.ETA), nil
	default:
		return "", fmt.Errorf("不認得的通知類型 %q", ev.Type)
	}
}

func eventDemo() {
	for _, msg := range []string{
		`{"type":"order.placed","data":{"order_id":"A-1024","total":280}}`,
		`{"type":"rider.assigned","data":{"order_id":"A-1024","rider":"王大哥","eta_minutes":25}}`,
		`{"type":"order.refunded","data":{"order_id":"A-0999"}}`,
	} {
		s, err := describe([]byte(msg))
		if err != nil {
			fmt.Println("錯誤：", err)
			continue
		}
		fmt.Println(s)
	}
}

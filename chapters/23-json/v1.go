package main

import (
	jsonv1 "encoding/json"
	"encoding/json/v2"
	"fmt"
	"strings"
)

// Dish 用來比較 v1 和 v2。
type Dish struct {
	Name  string   `json:"name"`
	Price int      `json:"price,omitempty"`
	Tags  []string `json:"tags"`
}

// v1v2Demo 同樣的 struct、同樣的輸入，v1 和 v2 的結果不一樣的地方。
func v1v2Demo() {
	d := Dish{Name: "蚵仔煎<招牌>"}
	b1, err1 := jsonv1.Marshal(d)
	b2, err2 := json.Marshal(d)
	fmt.Println("v1", string(b1), err1)
	fmt.Println("v2", string(b2), err2)

	b3, err3 := json.Marshal(d, jsonv1.DefaultOptionsV1()) // v2 的函式加上 v1 的選項，結果和 v1 一樣
	fmt.Println("v2+V1選項", string(b3), err3)

	in := []byte(`{"NAME":"大腸麵線","price":60}`)
	var a, b Dish
	err1 = jsonv1.Unmarshal(in, &a)
	err2 = json.Unmarshal(in, &b)
	fmt.Printf("v1 %q %v\n", a.Name, err1)
	fmt.Printf("v2 %q %v\n", b.Name, err2)

	dup := []byte(`{"name":"甲","name":"乙"}`)
	err1 = jsonv1.Unmarshal(dup, &a)
	err2 = json.Unmarshal(dup, &b)
	fmt.Println("v1", a.Name, err1)
	fmt.Println("v2", err2)
}

// v1APIDemo 現有程式碼裡常見的 v1 寫法：Decoder、DisallowUnknownFields、MarshalIndent。
func v1APIDemo() {
	dec := jsonv1.NewDecoder(strings.NewReader(`{"name":"肉圓","price":50,"size":"大"}`))
	dec.DisallowUnknownFields()
	var d Dish
	fmt.Println(dec.Decode(&d))

	b, err := jsonv1.MarshalIndent(Dish{Name: "肉圓", Price: 50, Tags: []string{"彰化"}}, "", "  ")
	fmt.Println(string(b), err)
}

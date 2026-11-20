package main

import (
	"encoding/json/v2"
	"fmt"
	"time"
)

// Order 是一筆外送訂單。每個欄位的 tag 示範一種常用的選項。
type Order struct {
	ID       int64     `json:"id,string"`         // 很大的整數用字串傳，JavaScript 那邊才不會失真
	Customer string    `json:"customer"`          // 一定會輸出
	Items    []string  `json:"items,omitempty"`   // JSON 是空的（[]）就省略
	Note     string    `json:"note,omitempty"`    // 空字串就省略
	Discount int       `json:"discount,omitzero"` // Go 的零值（0）就省略
	PaidAt   time.Time `json:"paid_at,omitzero"`  // 還沒付款（零值）就省略；time.Time 有 IsZero 方法
	Phone    string    `json:"-"`                 // 個資，永遠不輸出
}

// tagsDemo 同一筆訂單，付款前和付款後輸出的欄位不一樣。
func tagsDemo() {
	o := Order{ID: 2026111900000000123, Customer: "林小姐", Phone: "0912-345-678"}
	b, err := json.Marshal(o)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

	o.Items = []string{"排骨飯", "滷蛋"}
	o.Note = "不要香菜"
	o.Discount = 10
	o.PaidAt = time.Date(2026, time.November, 19, 12, 5, 0, 0, time.UTC)
	b, err = json.Marshal(o)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

// omitDemo omitzero 看的是 Go 的零值，omitempty 看的是編出來的 JSON 是不是空的。
func omitDemo() {
	type flags struct {
		CountZ int      `json:",omitzero"`
		CountE int      `json:",omitempty"` // 0 編成 JSON 的 0，不是「空的」，所以不會省略
		ListZ  []string `json:",omitzero"`  // 只有 nil 才省略
		ListE  []string `json:",omitempty"` // nil 和長度 0 都編成 []，都省略
	}
	for _, f := range []flags{{}, {ListZ: []string{}, ListE: []string{}}} {
		b, err := json.Marshal(f)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))
	}
}

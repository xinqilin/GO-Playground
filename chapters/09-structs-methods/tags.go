package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Member 每個欄位後面反引號裡的字串是 struct tag，給 encoding/json 這類套件讀。
type Member struct {
	Name   string `json:"name"`
	Phone  string `json:"phone,omitempty"`
	Points int    `json:"points"`
	note   string // 小寫開頭：其他 package（包括 encoding/json）看不到
}

func tagsDemo() {
	m := Member{Name: "小明", Points: 120, note: "常客"}
	data, err := json.Marshal(m)
	fmt.Println(string(data), err)

	f, _ := reflect.TypeFor[Member]().FieldByName("Phone")
	fmt.Println(f.Tag.Get("json"))
}

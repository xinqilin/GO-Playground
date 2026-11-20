package main

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
)

// explain 回傳 err 的訊息，但其中 json.SemanticError 的部分換成固定的寫法：
// 出錯的位置、JSON 的種類、Go 的型別和原因。v2 的錯誤訊息文字每次執行可能不同，
// 程式要判斷是哪一種錯誤時，看的是這些欄位，不是訊息文字。
func explain(err error) string {
	se, ok := errors.AsType[*json.SemanticError](err)
	if !ok {
		return fmt.Sprint(err)
	}
	var parts []string
	if se.JSONPointer != "" {
		parts = append(parts, "位置 "+string(se.JSONPointer))
	}
	if se.JSONKind != 0 {
		parts = append(parts, "JSON "+se.JSONKind.String())
	}
	if se.GoType != nil {
		parts = append(parts, "Go "+se.GoType.String())
	}
	if se.Err != nil {
		parts = append(parts, "原因 "+se.Err.Error())
	}
	return strings.Replace(err.Error(), se.Error(), "["+strings.Join(parts, "，")+"]", 1)
}

// anyDemo 解碼到 any 時，JSON 的每一種值會變成哪一種 Go 型別。
func anyDemo() {
	data := []byte(`{"name":"滷肉飯","price":35,"tags":["古早味"],"vegan":false,"note":null,"shop":{"id":7}}`)
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		fmt.Println(err)
		return
	}
	m := v.(map[string]any)
	for _, k := range slices.Sorted(maps.Keys(m)) {
		fmt.Printf("%-5s %T\n", k, m[k])
	}

	var big map[string]any
	if err := json.Unmarshal([]byte(`{"id":2026111900000000123}`), &big); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%.0f\n", big["id"]) // float64 只有 15 到 17 位有效數字，後面的位數變了
}

// strictDemo 大小寫不同的名稱、不認得的欄位，預設都被安靜地略過；RejectUnknownMembers 讓它們變成錯誤。
func strictDemo() {
	data := []byte(`{"Name":"雞腿飯","price":100,"spicy":true}`)

	var item MenuItem
	err := json.Unmarshal(data, &item)
	fmt.Printf("%q %d %v\n", item.Name, item.Price, err)

	err = json.Unmarshal(data, &item, json.RejectUnknownMembers(true))
	fmt.Println(explain(err), errors.Is(err, json.ErrUnknownName))
}

// decodeErrorsDemo 常見的解碼錯誤。前兩個是語意錯誤，用 JSON Pointer 指出是哪一個值；後兩個是語法錯誤。
func decodeErrorsDemo() {
	inputs := []string{
		`[{"name":"排骨飯","price":90},{"name":"雞腿飯","price":"一百"}]`, // 型別不對
		`[{"name":"排骨飯","price":90.5}]`,                           // int 放不下小數
		`[{"name":"排骨飯","price":90},]`,                            // 語法錯誤：多一個逗號
		`[{"name":"排骨飯","name":"雞腿飯"}]`,                           // 名稱重複
	}
	for _, in := range inputs {
		var menu []MenuItem
		fmt.Println(explain(json.Unmarshal([]byte(in), &menu)))
	}
}

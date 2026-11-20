package main

import (
	"bytes"
	"encoding/json/v2"
	"slices"
	"strings"
	"testing"
)

func Example_marshalDemo() {
	marshalDemo()
	// Output:
	// {"name":"排骨飯","price":90,"tags":["招牌","炸物"]}
	// {"name":"白飯","price":10,"tags":[]}
	// {
	//   "name": "排骨飯",
	//   "price": 90,
	//   "tags": [
	//     "招牌",
	//     "炸物"
	//   ]
	// }
}

func Example_unmarshalDemo() {
	unmarshalDemo()
	// Output:
	// {Name:雞腿飯 Price:100 Tags:[招牌] cost:0}
	// {"排骨":0,"滷蛋":30,"雞腿":12}
}

func Example_tagsDemo() {
	tagsDemo()
	// Output:
	// {"id":"2026111900000000123","customer":"林小姐"}
	// {"id":"2026111900000000123","customer":"林小姐","items":["排骨飯","滷蛋"],"note":"不要香菜","discount":10,"paid_at":"2026-11-19T12:05:00Z"}
}

func Example_omitDemo() {
	omitDemo()
	// Output:
	// {"CountE":0}
	// {"CountE":0,"ListZ":[]}
}

func Example_anyDemo() {
	anyDemo()
	// Output:
	// name  string
	// note  <nil>
	// price float64
	// shop  map[string]interface {}
	// tags  []interface {}
	// vegan bool
	// 2026111900000000000
}

func Example_strictDemo() {
	strictDemo()
	// Output:
	// "" 100 <nil>
	// [位置 /Name，JSON string，Go main.MenuItem，原因 unknown object member name] true
}

func Example_decodeErrorsDemo() {
	decodeErrorsDemo()
	// Output:
	// [位置 /1/price，JSON string，Go int]
	// [位置 /0/price，JSON number，Go int，原因 invalid syntax]
	// jsontext: invalid character ',' at start of value after offset 32
	// jsontext: duplicate object member name "name" within "/0"
}

func Example_lookupDemo() {
	lookupDemo()
	// Output:
	// 阿嬤便當 4.7 true
	// 沒有的欄位： false
}

func Example_eventDemo() {
	eventDemo()
	// Output:
	// 新訂單 A-1024，共 280 元
	// A-1024 由 王大哥 外送，約 25 分鐘
	// 錯誤： 不認得的通知類型 "order.refunded"
}

func Example_textDemo() {
	textDemo()
	// Output:
	// [{"dish":"麻辣鴨血","spice":"大辣"},{"dish":"清燉牛肉麵","spice":"不辣"}] <nil>
	// 酸辣湯 2 true <nil>
	// [位置 /spice，JSON string，Go main.Spice，原因 不認得的辣度 "超辣"]
	// {"大辣":3,"小辣":5} <nil>
}

func Example_uuidDemo() {
	uuidDemo()
	// Output:
	// {"id":"019a9b1c-3f00-7d2e-8a41-5c6d7e8f9012","table":7} <nil>
	// true <nil>
	// [位置 /id，JSON string，Go uuid.UUID，原因 invalid uuid]
	// 36
}

func Example_durationDemo() {
	durationDemo()
	// Output:
	// [位置 /prep，Go time.Duration，原因 no default representation]
	// {"dish":"紅燒牛肉","prep":"3h30m0s"} <nil>
}

func Example_linesDemo() {
	linesDemo()
	// Output:
	// 台中店 960
	// 台北店 1730
	// 第 4 筆：[位置 /amount，JSON string，Go int]
}

func Example_arrayDemo() {
	arrayDemo()
	// Output:
	// 2 <nil>
	// 0 要的是陣列，開頭卻是 {
}

func Example_encodeDemo() {
	encodeDemo()
	// Output:
	// {"shop":"台北店","amount":1730}
	// {"shop":"台中店","amount":960}
}

func Example_v1v2Demo() {
	v1v2Demo()
	// Output:
	// v1 {"name":"蚵仔煎\u003c招牌\u003e","tags":null} <nil>
	// v2 {"name":"蚵仔煎<招牌>","price":0,"tags":[]} <nil>
	// v2+V1選項 {"name":"蚵仔煎\u003c招牌\u003e","tags":null} <nil>
	// v1 "大腸麵線" <nil>
	// v2 "" <nil>
	// v1 乙 <nil>
	// v2 jsontext: duplicate object member name "name"
}

func Example_v1APIDemo() {
	v1APIDemo()
	// Output:
	// json: unknown field "size"
	// {
	//   "name": "肉圓",
	//   "price": 50,
	//   "tags": [
	//     "彰化"
	//   ]
	// } <nil>
}

func Example_csvDemo() {
	csvDemo()
	// Output:
	// 發票號碼,品項,金額
	// AB-12345678,排骨飯,90
	// AB-12345678,"珍珠奶茶, 微糖去冰",55
	// AB-12345679,"""招牌""滷肉飯",35
	// 3 行 "珍珠奶茶, 微糖去冰" <nil>
	// 第 3 行：金額 "二十" 不是整數
	// record on line 2: wrong number of fields
}

// TestSpiceRoundTrip 每一種辣度編碼後再解碼，都要回到原本的值；超出範圍的不能編碼。
func TestSpiceRoundTrip(t *testing.T) {
	for s := NotSpicy; s <= Hot; s++ {
		b, err := json.Marshal(s)
		if err != nil {
			t.Fatalf("Marshal(%d)：%v", s, err)
		}
		var got Spice
		if err := json.Unmarshal(b, &got); err != nil || got != s {
			t.Errorf("%s 解碼成 %d（%v），want %d", b, got, err, s)
		}
	}
	if _, err := json.Marshal(Spice(9)); err == nil {
		t.Error("Spice(9) 應該編碼失敗")
	}
}

// TestInvoicesRoundTrip 有逗號、引號、換行的品項，寫出去再讀回來要一模一樣。
func TestInvoicesRoundTrip(t *testing.T) {
	want := []InvoiceLine{
		{"AB-00000001", "排骨飯", 90},
		{"AB-00000001", "珍奶, 半糖", 55},
		{"AB-00000002", "備註：\"大\"份\n多飯", 120},
	}
	var buf bytes.Buffer
	if err := writeInvoices(&buf, want); err != nil {
		t.Fatal(err)
	}
	got, err := readInvoices(&buf)
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("讀回 %+v（%v）\nwant %+v", got, err, want)
	}
}

// TestExplain explain 的結果裡不會出現每次執行都可能不同的字，而且保留外層包裝的訊息。
func TestExplain(t *testing.T) {
	_, err := sumJSONLines(strings.NewReader(`{"shop":"台北店","amount":"九十"}`))
	got := explain(err)
	if strings.Contains(got, "cannot") || strings.Contains(got, "unable to") ||
		!strings.HasPrefix(got, "第 1 筆：") || !strings.Contains(got, "/amount") {
		t.Errorf("explain = %q", got)
	}
}

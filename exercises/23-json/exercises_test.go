package takeout

import (
	"encoding/json/v2"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestDecodeOrder(t *testing.T) {
	got, err := DecodeOrder([]byte(`{"id":"A-1","customer":"林小姐","items":[{"dish":"排骨飯","qty":2},{"dish":"滷蛋","qty":1}],"note":"不要香菜"}`))
	want := Order{ID: "A-1", Customer: "林小姐", Items: []Line{{"排骨飯", 2}, {"滷蛋", 1}}, Note: "不要香菜"}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("合法的訂單：got %+v（%v），want %+v", got, err, want)
	}

	bad := map[string]string{
		"欄位名稱打錯":     `{"id":"A-1","itmes":[{"dish":"排骨飯","qty":1}]}`,
		"沒有 id":      `{"items":[{"dish":"排骨飯","qty":1}]}`,
		"沒有品項":       `{"id":"A-1","items":[]}`,
		"數量是 0":      `{"id":"A-1","items":[{"dish":"排骨飯","qty":0}]}`,
		"品項名稱是空的":    `{"id":"A-1","items":[{"dish":"","qty":1}]}`,
		"數量是字串":      `{"id":"A-1","items":[{"dish":"排骨飯","qty":"一"}]}`,
		"不是合法的 JSON": `{"id":"A-1",`,
	}
	for name, in := range bad {
		if _, err := DecodeOrder([]byte(in)); err == nil {
			t.Errorf("%s：%s 應該回傳錯誤", name, in)
		}
	}
}

func TestReadOrders(t *testing.T) {
	in := `{"id":"A-1","items":[{"dish":"排骨飯","qty":1}]}

{"id":"A-2","customer":"王先生","items":[{"dish":"雞腿飯","qty":2}]}
{"id":"A-3","items":[{"dish":"滷肉飯","qty":3}]}
`
	got, err := ReadOrders(strings.NewReader(in))
	if err != nil || len(got) != 3 || got[1].Customer != "王先生" || got[2].Items[0].Qty != 3 {
		t.Errorf("三筆合法的訂單：got %+v（%v）", got, err)
	}

	in = `{"id":"A-1","items":[{"dish":"排骨飯","qty":1}]}
{"id":"A-2","items":[{"dish":"雞腿飯","qty":0}]}
{"id":"A-3","items":[{"dish":"滷肉飯","qty":3}]}
`
	if _, err := ReadOrders(strings.NewReader(in)); err == nil || !strings.HasPrefix(err.Error(), "第 2 筆：") {
		t.Errorf("第 2 筆數量是 0：err = %v，want 以「第 2 筆：」開頭的錯誤", err)
	}
	in = `{"id":"A-1","items":[{"dish":"排骨飯","qty":1}]}
{"id":"A-2",`
	if _, err := ReadOrders(strings.NewReader(in)); err == nil || !strings.HasPrefix(err.Error(), "第 2 筆：") {
		t.Errorf("第 2 筆不完整：err = %v，want 以「第 2 筆：」開頭的錯誤", err)
	}
	if got, err := ReadOrders(strings.NewReader("")); err != nil || len(got) != 0 {
		t.Errorf("空的輸入：got %v（%v），want 沒有訂單、沒有錯誤", got, err)
	}
}

func TestSumPrices(t *testing.T) {
	good := []struct {
		in   string
		want int
	}{
		{`{"price":90}`, 90},
		{`{"shop":"阿嬤便當","menu":[{"name":"滷肉飯","price":35},{"name":"雞腿飯","price":100,"addons":[{"name":"加蛋","price":10}]}]}`, 145},
		{`[{"price":1},[{"price":2}],{"x":{"y":{"price":3}}}]`, 6},
		{`{"prices":[1,2,3],"name":"price"}`, 0},
		{`[]`, 0},
	}
	for _, tt := range good {
		if got, err := SumPrices([]byte(tt.in)); err != nil || got != tt.want {
			t.Errorf("SumPrices(%s) = %d（%v），want %d", tt.in, got, err, tt.want)
		}
	}
	for _, in := range []string{`{"price":"90"}`, `{"price":9.5}`, `[{"price":1},{"price":null}]`, `{"price":`} {
		if _, err := SumPrices([]byte(in)); err == nil {
			t.Errorf("SumPrices(%s) 應該回傳錯誤", in)
		}
	}
}

func TestMoney(t *testing.T) {
	for m, want := range map[Money]string{
		0: `"NT$0"`, 90: `"NT$90"`, 999: `"NT$999"`, 1000: `"NT$1,000"`, 1280: `"NT$1,280"`,
		1234567: `"NT$1,234,567"`, -350: `"-NT$350"`, -1000: `"-NT$1,000"`,
	} {
		b, err := json.Marshal(m)
		if err != nil || string(b) != want {
			t.Errorf("Marshal(%d) = %s（%v），want %s", int(m), b, err, want)
		}
	}
	b, err := json.Marshal(Receipt{Order: "A-1", Total: 1280})
	if want := `{"order":"A-1","total":"NT$1,280"}`; err != nil || string(b) != want {
		t.Errorf("Marshal(Receipt) = %s（%v），want %s", b, err, want)
	}

	for in, want := range map[string]Money{
		`"NT$1,280"`: 1280, `"NT$1280"`: 1280, `"NT$0"`: 0, `"-NT$350"`: -350, `"NT$1,234,567"`: 1234567,
	} {
		var got Money
		if err := json.Unmarshal([]byte(in), &got); err != nil || got != want {
			t.Errorf("Unmarshal(%s) = %d（%v），want %d", in, int(got), err, int(want))
		}
	}
	for _, in := range []string{`"1280"`, `"$1,280"`, `"NT$"`, `"NT$12a"`, `"NT$-5"`, `"NT$+5"`, `"NT$,"`} {
		var got Money
		if err := json.Unmarshal([]byte(in), &got); err == nil {
			t.Errorf("Unmarshal(%s) 應該回傳錯誤，卻得到 %d", in, int(got))
		}
	}
}

func TestReadMenu(t *testing.T) {
	in := "價格,品名,備註,分類\n" +
		"35,滷肉飯,招牌,飯類\n" +
		"100,\"雞腿飯, 大\",,飯類\n" +
		"60,大腸麵線,,麵類\n"
	got, err := ReadMenu(strings.NewReader(in))
	want := []MenuRow{{"滷肉飯", "飯類", 35}, {"雞腿飯, 大", "飯類", 100}, {"大腸麵線", "麵類", 60}}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("ReadMenu = %+v（%v）\nwant %+v", got, err, want)
	}

	if _, err := ReadMenu(strings.NewReader("品名,價格\n滷肉飯,35\n")); err == nil {
		t.Error("缺少「分類」欄，應該回傳錯誤")
	}

	in = "品名,分類,價格\n" +
		"\"招牌\n便當\",飯類,120\n" + // 這個欄位跨兩行：第 2、3 行
		"紅茶,飲料,二十\n" // 第 4 行
	if _, err := ReadMenu(strings.NewReader(in)); err == nil || !strings.Contains(err.Error(), "第 4 行") {
		t.Errorf("價格不是整數：err = %v，want 包含「第 4 行」的錯誤", err)
	}
}

type placed struct {
	OrderID string `json:"order_id"`
	Total   int    `json:"total"`
}

type assigned struct {
	OrderID string `json:"order_id"`
	Rider   string `json:"rider"`
}

func TestRouter(t *testing.T) {
	r := NewRouter()
	var log []string
	r.Handle("order.placed", func(p placed) error {
		log = append(log, p.OrderID+" 下單")
		if p.Total <= 0 {
			return errors.New("金額不正確")
		}
		return nil
	})
	r.Handle("rider.assigned", func(a assigned) error {
		log = append(log, a.OrderID+" 由 "+a.Rider+" 外送")
		return nil
	})

	if err := r.Dispatch([]byte(`{"type":"order.placed","data":{"order_id":"A-1","total":280}}`)); err != nil {
		t.Errorf("order.placed：%v", err)
	}
	if err := r.Dispatch([]byte(`{"type":"rider.assigned","data":{"order_id":"A-1","rider":"王大哥"}}`)); err != nil {
		t.Errorf("rider.assigned：%v", err)
	}
	if err := r.Dispatch([]byte(`{"type":"order.placed","data":{"order_id":"A-2","total":0}}`)); err == nil {
		t.Error("處理函式回傳的錯誤，Dispatch 要回傳出來")
	}
	for _, msg := range []string{
		`{"type":"order.refunded","data":{"order_id":"A-1"}}`,            // 沒有登記
		`{"type":"rider.assigned","data":{"order_id":"A-3","rider":7}}`,  // 型別不對
		`{"type":"rider.assigned","data":{"order_id":"A-3","ridr":"林"}}`, // 不認得的欄位
		`{"type":"order.placed",`,                                        // 不是合法的 JSON
	} {
		if err := r.Dispatch([]byte(msg)); err == nil {
			t.Errorf("Dispatch(%s) 應該回傳錯誤", msg)
		}
	}
	want := []string{"A-1 下單", "A-1 由 王大哥 外送", "A-2 下單"}
	if !reflect.DeepEqual(log, want) {
		t.Errorf("處理函式被呼叫的紀錄 %q，want %q", log, want)
	}

	r.Handle("order.placed", func(p placed) error { // 取代前面的
		log = append(log, "新的："+p.OrderID)
		return nil
	})
	if err := r.Dispatch([]byte(`{"type":"order.placed","data":{"order_id":"A-9","total":1}}`)); err != nil || len(log) == 0 || log[len(log)-1] != "新的：A-9" {
		t.Errorf("重新登記之後應該呼叫新的處理函式，紀錄 %q（%v）", log, err)
	}
}

func TestCupSize(t *testing.T) {
	for in, want := range map[string]CupSize{
		`{"name":"紅茶","size":"大杯"}`: LargeCup,
		`{"name":"綠茶","size":"中杯"}`: MediumCup,
	} {
		var d Drink
		if err := json.Unmarshal([]byte(in), &d); err != nil || d.Size != want {
			t.Errorf("Unmarshal(%s)：Size = %d（%v），want %d", in, int(d.Size), err, int(want))
		}
	}
	var d Drink
	if err := json.Unmarshal([]byte(`{"name":"奶茶","size":"超大杯"}`), &d); err == nil {
		t.Error("超大杯應該回傳錯誤")
	}
	b, err := json.Marshal(Drink{Name: "紅茶", Size: LargeCup})
	if err != nil {
		t.Fatal(err)
	}
	var back Drink
	if err := json.Unmarshal(b, &back); err != nil || back.Size != LargeCup {
		t.Errorf("%s 讀回來 Size = %d（%v），want %d", b, int(back.Size), err, int(LargeCup))
	}
}

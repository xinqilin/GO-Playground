//go:build solution

package takeout

import (
	"encoding/csv"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

// DecodeOrder 先嚴格地解碼，再檢查內容。
func DecodeOrder(data []byte) (Order, error) {
	var o Order
	if err := json.Unmarshal(data, &o, json.RejectUnknownMembers(true)); err != nil {
		return Order{}, err
	}
	if o.ID == "" {
		return Order{}, errors.New("缺少 id")
	}
	if len(o.Items) == 0 {
		return Order{}, errors.New("沒有任何品項")
	}
	for i, l := range o.Items {
		if l.Dish == "" || l.Qty < 1 {
			return Order{}, fmt.Errorf("第 %d 個品項不正確：%+v", i+1, l)
		}
	}
	return o, nil
}

// ReadOrders 用 ReadValue 一次取出一筆完整的 JSON，交給 DecodeOrder，規則只寫一次。
func ReadOrders(r io.Reader) ([]Order, error) {
	dec := jsontext.NewDecoder(r)
	var orders []Order
	for n := 1; ; n++ {
		v, err := dec.ReadValue()
		if errors.Is(err, io.EOF) {
			return orders, nil
		}
		if err != nil {
			return nil, fmt.Errorf("第 %d 筆：%w", n, err)
		}
		o, err := DecodeOrder(v)
		if err != nil {
			return nil, fmt.Errorf("第 %d 筆：%w", n, err)
		}
		orders = append(orders, o)
	}
}

// SumPrices 解碼到 any，再遞迴走過每一個物件和陣列。
func SumPrices(data []byte) (int, error) {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return 0, err
	}
	return sumPrices(v)
}

func sumPrices(v any) (int, error) {
	sum := 0
	switch v := v.(type) {
	case map[string]any:
		for k, x := range v {
			if k == "price" {
				f, ok := x.(float64)
				if !ok || f != math.Trunc(f) {
					return 0, fmt.Errorf("price 不是整數：%v", x)
				}
				sum += int(f)
				continue
			}
			n, err := sumPrices(x)
			if err != nil {
				return 0, err
			}
			sum += n
		}
	case []any:
		for _, x := range v {
			n, err := sumPrices(x)
			if err != nil {
				return 0, err
			}
			sum += n
		}
	}
	return sum, nil
}

// MarshalText 從個位數往前，每三位插一個逗號。
func (m Money) MarshalText() ([]byte, error) {
	sign, n := "", int64(m)
	if n < 0 {
		sign, n = "-", -n
	}
	digits := strconv.FormatInt(n, 10)
	var b strings.Builder
	for i, d := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(d)
	}
	return []byte(sign + "NT$" + b.String()), nil
}

// UnmarshalText 先拿掉負號和 NT$，再確認剩下的只有數字和逗號。
func (m *Money) UnmarshalText(b []byte) error {
	s := string(b)
	neg := strings.HasPrefix(s, "-")
	rest, ok := strings.CutPrefix(strings.TrimPrefix(s, "-"), "NT$")
	if !ok {
		return fmt.Errorf("金額 %q 要以 NT$ 開頭", s)
	}
	digits := strings.ReplaceAll(rest, ",", "")
	if digits == "" || strings.Trim(digits, "0123456789") != "" {
		return fmt.Errorf("金額 %q 的格式不正確", s)
	}
	n, err := strconv.Atoi(digits)
	if err != nil {
		return fmt.Errorf("金額 %q：%w", s, err)
	}
	if neg {
		n = -n
	}
	*m = Money(n)
	return nil
}

// ReadMenu 先從標題找出三個欄位的位置，再一行一行讀；行號用 FieldPos，
// 欄位裡有換行時也是對的。
func ReadMenu(r io.Reader) ([]MenuRow, error) {
	cr := csv.NewReader(r)
	header, err := cr.Read()
	if err != nil {
		return nil, err
	}
	col := map[string]int{}
	for i, h := range header {
		col[h] = i
	}
	for _, name := range []string{"品名", "分類", "價格"} {
		if _, ok := col[name]; !ok {
			return nil, fmt.Errorf("缺少「%s」欄", name)
		}
	}
	var rows []MenuRow
	for {
		rec, err := cr.Read()
		if errors.Is(err, io.EOF) {
			return rows, nil
		}
		if err != nil {
			return nil, err
		}
		price, err := strconv.Atoi(rec[col["價格"]])
		if err != nil {
			line, _ := cr.FieldPos(col["價格"])
			return nil, fmt.Errorf("第 %d 行：價格 %q 不是整數", line, rec[col["價格"]])
		}
		rows = append(rows, MenuRow{Name: rec[col["品名"]], Category: rec[col["分類"]], Price: price})
	}
}

// Router 把每個處理函式包成「收 jsontext.Value」的函式存起來，型別參數只在 Handle 裡用到。
type Router struct {
	handlers map[string]func(jsontext.Value) error
}

func NewRouter() *Router {
	return &Router{handlers: map[string]func(jsontext.Value) error{}}
}

func (r *Router) Handle[T any](typ string, fn func(T) error) {
	r.handlers[typ] = func(data jsontext.Value) error {
		var v T
		if err := json.Unmarshal(data, &v, json.RejectUnknownMembers(true)); err != nil {
			return fmt.Errorf("%s 的 data：%w", typ, err)
		}
		return fn(v)
	}
}

func (r *Router) Dispatch(msg []byte) error {
	var env struct {
		Type string         `json:"type"`
		Data jsontext.Value `json:"data"`
	}
	if err := json.Unmarshal(msg, &env); err != nil {
		return err
	}
	h, ok := r.handlers[env.Type]
	if !ok {
		return fmt.Errorf("沒有登記 %q 的處理函式", env.Type)
	}
	return h(env.Data)
}

// UnmarshalText 的 receiver 要是指標：值 receiver 改到的只是一份複本。
func (c *CupSize) UnmarshalText(b []byte) error {
	switch string(b) {
	case "中杯":
		*c = MediumCup
	case "大杯":
		*c = LargeCup
	default:
		return fmt.Errorf("不認得的杯型 %q", b)
	}
	return nil
}

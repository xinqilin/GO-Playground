//go:build solution

package kiosk

import (
	"fmt"
	"strconv"
)

// AsFloat 在每個 case 裡，x 已經是那個具體型別，可以直接轉換。
func AsFloat(v any) (float64, bool) {
	switch x := v.(type) {
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	case float32:
		return float64(x), true
	case float64:
		return x, true
	case string:
		f, err := strconv.ParseFloat(x, 64)
		if err != nil {
			return 0, false
		}
		return f, true
	default:
		return 0, false
	}
}

// Fee 用 (extra + 999) / 1000 做「無條件進位」的整數除法。
func (HomeDelivery) Fee(grams int) int {
	if grams <= 5000 {
		return 100
	}
	extra := grams - 5000
	return 100 + (extra+999)/1000*30
}

// Cheapest 只透過 Shipper 呼叫 Fee，不管背後是哪一種寄送方式。
func Cheapest(options []Shipper, grams int) int {
	best, bestFee := -1, 0
	for i, s := range options {
		fee := s.Fee(grams)
		if fee < 0 {
			continue
		}
		if best == -1 || fee < bestFee {
			best, bestFee = i, fee
		}
	}
	return best
}

// Write 先複製一份再修改，呼叫者的 p 才不會被改到。
func (m Masker) Write(p []byte) (int, error) {
	masked := make([]byte, len(p))
	for i, b := range p {
		if b >= '0' && b <= '9' {
			b = '*'
		}
		masked[i] = b
	}
	if _, err := m.W.Write(masked); err != nil {
		return 0, err
	}
	return len(p), nil
}

// Total 用 type assertion 問每個品項「你也是 Discounter 嗎？」。
func Total(items []Item) int {
	total := 0
	for _, it := range items {
		price := it.Price()
		if d, ok := it.(Discounter); ok {
			price = max(price-d.Discount(), 0)
		}
		total += price
	}
	return total
}

// Send 讓 Fallback 本身也滿足 Sender，所以 Fallback 裡面還可以再放 Fallback。
func (f Fallback) Send(msg string) error {
	if len(f) == 0 {
		return ErrNoSender
	}
	var last error
	for _, s := range f {
		err := s.Send(msg)
		if err == nil {
			return nil
		}
		last = err
	}
	return last
}

// CheckOrder 沒有錯誤時明確地回傳 nil，而不是回傳一個值為 nil 的 *OrderError。
func CheckOrder(o Order) error {
	if o.Items == 0 {
		return &OrderError{Reason: "沒有點東西"}
	}
	if !o.Paid {
		return &OrderError{Reason: "還沒付款"}
	}
	return nil
}

// Kind 的 case 由上往下比對，第一個符合的就勝出。Price 也是 fmt.Stringer，
// 所以具體的 Price 要寫在 fmt.Stringer 前面。
func Kind(v any) string {
	switch v.(type) {
	case Price:
		return "價格"
	case fmt.Stringer:
		return "可以印的"
	case error:
		return "錯誤"
	default:
		return "其他"
	}
}

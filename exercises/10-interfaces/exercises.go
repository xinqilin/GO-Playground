//go:build !solution

package kiosk

import "fmt"

// AsFloat 把 v 轉成 float64。int、int64、float32、float64 直接轉換；
// string 用 strconv.ParseFloat 解析，例如 "85.5"，解析失敗回傳 0, false。
// 其他型別（包括 nil 和 bool）都回傳 0, false。
func AsFloat(v any) (float64, bool) {
	// TODO: 實作
	return 0, false
}

// Fee 是宅配的運費：5000 公克以內 100 元；超過的部分每 1000 公克加 30 元，
// 不足 1000 公克也算 1000 公克。例如 5001 公克是 130 元，8000 公克是 190 元。
func (HomeDelivery) Fee(grams int) int {
	// TODO: 實作
	return 0
}

// Cheapest 回傳 options 裡運費最便宜的那一個的索引。寄不了的（Fee 回傳 -1）不算；
// 一樣便宜時取索引小的；全部都寄不了，或 options 是空的，回傳 -1。
func Cheapest(options []Shipper, grams int) int {
	// TODO: 實作
	return 0
}

// Write 讓 Masker 成為 io.Writer：把 p 裡的 ASCII 數字換成 '*' 之後寫到 m.W。
// 不可以改到呼叫者的 p。成功時回傳 len(p), nil；m.W 寫入失敗時回傳它的錯誤。
func (m Masker) Write(p []byte) (int, error) {
	// TODO: 實作
	return 0, nil
}

// Total 回傳購物車的總價。有實作 Discounter 的品項要扣掉折扣，
// 但每一樣最低是 0 元，不會變成負的。
func Total(items []Item) int {
	// TODO: 實作
	return 0
}

// Send 依序嘗試 f 裡的每一種 Sender：有一種成功就回傳 nil，後面的不再嘗試。
// 全部失敗時回傳最後一個錯誤；f 是空的時回傳 ErrNoSender。
func (f Fallback) Send(msg string) error {
	// TODO: 實作
	return nil
}

// CheckOrder 檢查訂單：至少要點一樣東西，而且要付過錢。沒有問題時回傳 nil。
//
// 找錯：明明是正常的訂單，呼叫的人檢查 err != nil 卻成立。
func CheckOrder(o Order) error {
	var err *OrderError
	if o.Items == 0 {
		err = &OrderError{Reason: "沒有點東西"}
	} else if !o.Paid {
		err = &OrderError{Reason: "還沒付款"}
	}
	return err
}

// Kind 回傳 v 屬於哪一類：Price 是「價格」，其他有 String 方法的是「可以印的」，
// error 是「錯誤」，其餘都是「其他」。
//
// 找錯：Kind(Price(85)) 回傳的是「可以印的」。
func Kind(v any) string {
	switch v.(type) {
	case fmt.Stringer:
		return "可以印的"
	case Price:
		return "價格"
	case error:
		return "錯誤"
	default:
		return "其他"
	}
}

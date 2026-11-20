//go:build !solution

package takeout

import (
	"fmt"
	"io"
)

// DecodeOrder 解析一筆訂單的 JSON，並檢查：
//   - 不認得的欄位是錯誤（客戶端把欄位名稱打錯時才會知道）
//   - id 不能是空字串
//   - 至少要有一個品項；每個品項的 dish 不能是空字串，qty 至少是 1
//
// 有任何問題時回傳錯誤。
func DecodeOrder(data []byte) (Order, error) {
	// TODO: 實作
	return Order{}, nil
}

// ReadOrders 從 r 讀 JSON Lines：一行一筆訂單，規則和 DecodeOrder 一樣，空行略過。
// 遇到第一筆有問題的訂單就停下來，錯誤訊息要以「第 N 筆：」開頭（N 從 1 起算）。
// 讀到結尾時回傳全部的訂單。
func ReadOrders(r io.Reader) ([]Order, error) {
	// TODO: 實作
	return nil, nil
}

// SumPrices 把 JSON 裡所有名稱是 "price" 的值加起來，不管它在第幾層、在物件裡還是陣列裡。
// price 的值一定要是整數，不是數字或有小數時回傳錯誤；JSON 本身有錯也回傳錯誤。
func SumPrices(data []byte) (int, error) {
	// TODO: 實作
	return 0, nil
}

// MarshalText 把金額寫成 "NT$1,280" 的樣子：前面是 NT$，每三位數一個逗號；
// 負數在最前面加上 -，例如 "-NT$350"。
func (m Money) MarshalText() ([]byte, error) {
	// TODO: 實作
	return nil, nil
}

// UnmarshalText 解析 MarshalText 寫出來的格式。逗號可有可無（"NT$1280" 也可以），
// 但一定要有 NT$，而且後面只能有數字和逗號，至少一個數字；不符合時回傳錯誤。
func (m *Money) UnmarshalText(b []byte) error {
	// TODO: 實作
	return nil
}

// ReadMenu 讀菜單的 CSV。第一行是標題，欄位的順序不固定，但一定要有「品名」「分類」「價格」三欄，
// 其他欄位略過；缺少任何一欄時回傳錯誤。價格不是整數時回傳錯誤，
// 訊息要包含那個價格在檔案的第幾行，例如「第 4 行」（標題是第 1 行）。
func ReadMenu(r io.Reader) ([]MenuRow, error) {
	// TODO: 實作
	return nil, nil
}

// Router 依照通知的 type，把 data 解成對應的型別，再交給登記的處理函式。
// 通知的格式是 {"type": "...", "data": {...}}。
type Router struct {
	// TODO: 加上需要的欄位
}

// NewRouter 建立一個沒有任何處理函式的 Router。
func NewRouter() *Router {
	// TODO: 實作
	return &Router{}
}

// Handle 登記 typ 的處理函式：data 會解成 T（不認得的欄位算錯誤），再呼叫 fn。
// 同一個 typ 登記兩次時，後面的取代前面的。
func (r *Router) Handle[T any](typ string, fn func(T) error) {
	// TODO: 實作
}

// Dispatch 解析一則通知並交給對應的處理函式，回傳處理函式的錯誤。
// 通知不是合法的 JSON、type 沒有登記、data 解不成登記的型別時，回傳錯誤，而且不呼叫任何處理函式。
func (r *Router) Dispatch(msg []byte) error {
	// TODO: 實作
	return nil
}

// UnmarshalText 把「中杯」「大杯」解析回 CupSize，其他的字回傳錯誤。
//
// 找錯：解碼 {"name":"紅茶","size":"大杯"} 沒有錯誤，Size 卻還是 UnknownSize。
func (c CupSize) UnmarshalText(b []byte) error {
	switch string(b) {
	case "中杯":
		c = MediumCup
	case "大杯":
		c = LargeCup
	default:
		return fmt.Errorf("不認得的杯型 %q", b)
	}
	return nil
}

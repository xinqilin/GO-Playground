// 第 11 章：錯誤處理。
//
// 在 Go 裡錯誤是一個普通的值：函式把它當成最後一個回傳值交給你，你決定怎麼處理。
// 執行：go run ./chapters/11-errors
package main

import (
	"errors"
	"fmt"
)

// ErrNotFound 是 sentinel error：一個固定的錯誤值，讓呼叫端用 errors.Is 比對。
var ErrNotFound = errors.New("找不到")

// ValidationError 是帶有欄位資訊的自訂錯誤型別，呼叫端用 errors.As 取出細節。
type ValidationError struct {
	Field string
	Value int
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("欄位 %s 的值 %d 不合法", e.Field, e.Value)
}

var stock = map[string]int{"蘋果": 3}

// findStock 找不到時，用 %w 把 ErrNotFound 包起來，加上這一層的說明。
func findStock(name string) (int, error) {
	n, ok := stock[name]
	if !ok {
		return 0, fmt.Errorf("查庫存 %q：%w", name, ErrNotFound)
	}
	return n, nil
}

// order 檢查數量、查庫存；每一層都只補充自己知道的資訊再往上傳。
func order(name string, qty int) error {
	if qty <= 0 {
		return &ValidationError{Field: "qty", Value: qty}
	}
	n, err := findStock(name)
	if err != nil {
		return fmt.Errorf("下單失敗：%w", err)
	}
	if qty > n {
		return fmt.Errorf("下單失敗：%s 只剩 %d 個", name, n)
	}
	return nil
}

func main() {
	fmt.Println("== 錯誤是值 ==")
	orders := []struct {
		name string
		qty  int
	}{{"蘋果", 2}, {"香蕉", 1}, {"蘋果", 0}, {"蘋果", 5}}
	for _, o := range orders {
		fmt.Println(o.name, o.qty, "→", order(o.name, o.qty))
	}

	fmt.Println("== errors.Is：比對特定錯誤 ==")
	err := order("香蕉", 1)
	fmt.Println(errors.Is(err, ErrNotFound)) // 一層一層拆開包裝去找

	fmt.Println("== errors.As 與 errors.AsType：取出錯誤型別 ==")
	err = order("蘋果", -1)
	var ve *ValidationError
	if errors.As(err, &ve) {
		fmt.Println("As:", ve.Field, ve.Value)
	}
	if ve, ok := errors.AsType[*ValidationError](err); ok { // Go 1.26 起的泛型寫法
		fmt.Println("AsType:", ve.Field, ve.Value)
	}

	fmt.Println("== panic 與 recover ==")
	fmt.Println(safeDivide(10, 2))
	fmt.Println(safeDivide(1, 0))
}

// safeDivide 在 defer 裡用 recover 把 panic 轉成 error。
// 只在程式邊界（例如 HTTP handler 最外層）這樣做，不要拿 panic 當一般的錯誤處理。
func safeDivide(a, b int) (q int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recover: %v", r)
		}
	}()
	return a / b, nil
}

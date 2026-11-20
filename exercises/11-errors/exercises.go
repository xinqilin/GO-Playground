//go:build !solution

package orders

import (
	"fmt"
	"io"
)

// ParsePickup 把取餐時間 "HH:MM" 轉成從午夜算起的分鐘數，例如 "18:30" 是 1110。
//
// 用冒號分成時和分，兩邊都要是數字，時在 0 到 23、分在 0 到 59，否則回傳
// fmt.Errorf("取餐時間 %q：%w", s, ErrTimeFormat)。
// 格式正確但不在 11:00 到 20:00 之間（兩端都算營業時間）時，回傳
// fmt.Errorf("取餐時間 %s：%w", s, ErrClosed)。
func ParsePickup(s string) (int, error) {
	// TODO: 實作
	return 0, nil
}

// Error 回傳像「排骨飯 每人限購 2 份」這樣的訊息。
func (e *LimitError) Error() string {
	// TODO: 實作
	return ""
}

// Buy 檢查 qty 有沒有超過 limits 裡 item 的限購數量。超過時回傳 *LimitError；
// 沒有超過，或 item 不在 limits 裡（表示不限購），回傳 nil。
func Buy(item string, qty int, limits map[string]int) error {
	// TODO: 實作
	return nil
}

// CheckCart 用 Buy 檢查 cart 裡的每一樣東西（品名對應數量），依品名排序，
// 把所有超過限購的錯誤用 errors.Join 合在一起回傳。全部沒問題時回傳 nil。
func CheckCart(cart map[string]int, limits map[string]int) error {
	// TODO: 實作
	return nil
}

// Retry 最多呼叫 op times 次，直到它回傳 nil；成功就回傳 nil。times 小於 1 時當成 1。
// op 回傳的錯誤如果包著 ErrPermanent，就不要再試，直接回傳那個錯誤。
// 試了 times 次都失敗時，回傳 fmt.Errorf("試了 %d 次：%w", times, 最後一次的錯誤)。
func Retry(times int, op func() error) error {
	// TODO: 實作
	return nil
}

// Causes 回傳 err 以及它包著的每一個錯誤的訊息，深度優先：先放 err 自己，
// 再依序走進它包著的每一個錯誤。包著的錯誤可能來自 Unwrap() error（一個），
// 也可能來自 Unwrap() []error（好幾個，例如 errors.Join 和用了好幾個 %w 的 fmt.Errorf）。
// err 是 nil 時回傳 nil。
func Causes(err error) []string {
	// TODO: 實作
	return nil
}

// IsSoldOut 回報 err 是不是因為賣完了。
//
// 找錯：直接拿到 ErrSoldOut 時是對的，被 fmt.Errorf 包過一層就認不出來；
// 而另一個剛好也寫著「賣完了」的錯誤，卻被誤認成 ErrSoldOut。
func IsSoldOut(err error) bool {
	return err != nil && err.Error() == ErrSoldOut.Error()
}

// SaveLines 把每一行寫到 w，最後關閉 w。寫入或關閉失敗都要回報。
//
// 找錯：關閉失敗時（例如磁碟滿了，資料其實沒有存進去），SaveLines 還是回傳 nil。
func SaveLines(w io.WriteCloser, lines []string) error {
	defer w.Close()
	for _, l := range lines {
		if _, err := fmt.Fprintln(w, l); err != nil {
			return err
		}
	}
	return nil
}

//go:build !solution

package testkit

import "testing"

// SubtestNames 回傳依序呼叫 t.Run(names[i], …) 時，go test 顯示的子測試名稱（簡化版的規則）：
//   - 名稱裡的空白字元（unicode.IsSpace）都換成 _。
//   - 換完之後和前面重複的名稱，第 n 次重複加上 #01、#02… 這種兩位數的編號。
//   - 空字串第一次就變成 #00，之後是 #01、#02…
//
// 例如 {"a b", "x", "x", ""} 變成 {"a_b", "x", "x#01", "#00"}。
func SubtestNames(names []string) []string {
	// TODO: 實作
	return nil
}

// DiffLines 逐行比較 got 和 want（用 "\n" 切開），每一個不同的行回傳一行說明：
//
//	第 2 行：got "排骨飯 90"，want "排骨飯 95"
//
// 其中一邊比較短時，少掉的那一邊寫成 （沒有這一行）。完全相同時回傳 nil。
func DiffLines(got, want string) []string {
	// TODO: 實作
	return nil
}

// Notify 送出一則簡訊：phone 在 FailOn 裡時回傳 ErrUnreachable，
// 否則把 "phone：msg" 記下來並回傳 nil。會被好幾個 goroutine 同時呼叫。
func (r *Recorder) Notify(phone, msg string) error {
	// TODO: 實作
	return nil
}

// Sent 依記下的順序回傳所有送出的簡訊。回傳的是複本，呼叫的人改了它也不會影響 Recorder。
func (r *Recorder) Sent() []string {
	// TODO: 實作
	return nil
}

// Golden 比對 got 和 path 這個 golden file 的內容：
//   - update 是 true 時，把 got 寫進 path（需要的目錄也要建立），然後結束。
//   - 讀不到 path 時，用 tb.Fatalf 回報，訊息要包含 path 和「-update」，提醒要怎麼產生它。
//   - 內容不同時，用 tb.Errorf 回報一次，訊息要包含 path（可以加上練習 2 的 DiffLines）。
//
// 它是測試輔助函式，記得呼叫 tb.Helper()。
func Golden(tb testing.TB, path string, got []byte, update bool) {
	// TODO: 實作
}

// Shrink 把一個會讓 fails 回傳 true 的 input 縮小，就像 fuzz 縮小失敗的輸入：
// 回傳的 slice 仍然讓 fails 回傳 true，而且再拿掉任何一個元素，fails 就會回傳 false。
// 元素的先後順序要保持不變，也不能修改 input 本身。
func Shrink(input []int, fails func([]int) bool) []int {
	// TODO: 實作
	return nil
}

// Run 執行所有清理函式，規則和 t.Cleanup 一樣：後加入的先執行（先開的檔案最後關），
// 而且每個函式只執行一次：再呼叫一次 Run 什麼都不做。
//
// 找錯：資料庫比檔案先關掉了，而且呼叫兩次 Run 會關兩次。
func (c *Cleanups) Run() {
	for _, f := range c.fns {
		f()
	}
}

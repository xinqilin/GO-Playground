//go:build !solution

package functions

import "fmt"

// MinMax 回傳 nums 裡最小和最大的值。沒有傳任何數字時回傳 0, 0, false。
func MinMax(nums ...int) (lo, hi int, ok bool) {
	// TODO: 實作
	return 0, 0, false
}

// CountIf 回傳 nums 裡有幾個數字讓 keep 回傳 true。
func CountIf(nums []int, keep func(int) bool) int {
	// TODO: 實作
	return 0
}

// NewDispenser 做出一台號碼牌機，回傳兩個函式：
// next 每呼叫一次就發出下一張號碼牌，從 prefix 加上 001 開始，例如 "A001"、"A002"；
// reset 讓號碼從 001 重新開始。每台機器的號碼各自獨立。
func NewDispenser(prefix string) (next func() string, reset func()) {
	// TODO: 實作
	return func() string { return "" }, func() {}
}

// PercentOff 回傳一條「打折」的規則：減掉 pct%，小數無條件捨去。
// 例如 PercentOff(10) 會把 1000 變成 900，PercentOff(15) 會把 85 變成 72。
func PercentOff(pct int) Rule {
	// TODO: 實作
	return func(price int) int { return price }
}

// MinusOver 回傳一條「滿額折抵」的規則：價格大於等於 threshold 時減掉 off，否則不變。
func MinusOver(threshold, off int) Rule {
	// TODO: 實作
	return func(price int) int { return price }
}

// ApplyRules 依照順序把每條規則套用到 price 上，回傳最後的價格。結果小於 0 時回傳 0。
func ApplyRules(price int, rules ...Rule) int {
	// TODO: 實作
	return 0
}

// ProcessAll 依序處理每個資源：開啟、使用、關閉，處理完一個才開下一個。
// 正確的紀錄是 open a、use a、close a、open b、use b、close b。
//
// 找錯：資源都等到最後才關閉，而且關閉的順序是反的。
func ProcessAll(t *Tracker, names []string) {
	for _, name := range names {
		t.Open(name)
		defer t.Close(name)
		t.Use(name)
	}
}

// SafeDiv 回傳 a / b。b 是 0 時不能讓程式當掉，要回傳 0 和錯誤，
// 錯誤訊息是「除法失敗：」接上 panic 的內容。
//
// 找錯：b 是 0 時，它回傳的錯誤是 nil。
func SafeDiv(a, b int) (int, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("除法失敗：%v", r)
		}
	}()
	return a / b, err
}

// PayWays 回傳用 coins 裡的面額湊出 amount 元有幾種方法。每種面額可以用任意多個，不計順序，
// 例如用 1、5、10 元湊 10 元有 4 種：10、5+5、5+1+1+1+1+1、十個 1。
// amount 是 0 時有 1 種（什麼都不用付），amount 小於 0 時是 0 種。coins 裡沒有重複、0 或負數。
func PayWays(amount int, coins []int) int {
	// TODO: 實作
	return 0
}

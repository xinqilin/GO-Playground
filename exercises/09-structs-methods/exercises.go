//go:build !solution

package shop

// DoubleBooked 依序檢查訂位紀錄，回傳第一個「被訂第二次」的座位和 true。
// 例如 A5、B3、B3、A5 要回傳 B3，因為 B3 比 A5 先被訂第二次。
// 沒有重複時回傳 Seat{} 和 false。
func DoubleBooked(orders []Seat) (Seat, bool) {
	// TODO: 實作
	return Seat{}, false
}

// Add 把 qty 個 name（單價 price）加進購物車。
// 購物車裡已經有同名的品項時，只增加那一行的數量，不新增一行。
func (c *Cart) Add(name string, price, qty int) {
	// TODO: 實作
}

// Total 回傳購物車的總價。
func (c *Cart) Total() int {
	// TODO: 實作
	return 0
}

// TopUp 儲值 amount 元。amount 要大於 0 而且是 100 的倍數，
// 否則回傳 ErrBadAmount，餘額不變。
func (c *Card) TopUp(amount int) error {
	// TODO: 實作
	return nil
}

// Pay 扣 fare 元車資。餘額不夠時回傳 ErrInsufficient，餘額不變。
func (c *Card) Pay(fare int) error {
	// TODO: 實作
	return nil
}

// Add 把 n 個 name 加進庫存。零值的 Inventory 也要能直接用。
func (inv *Inventory) Add(name string, n int) {
	// TODO: 實作
}

// Remove 從庫存取出 n 個 name。庫存不夠時回傳 ErrInsufficient，庫存不變。
func (inv *Inventory) Remove(name string, n int) error {
	// TODO: 實作
	return nil
}

// Count 回傳 name 的庫存量；沒有這個品項時回傳 0。
func (inv *Inventory) Count(name string) int {
	// TODO: 實作
	return 0
}

// Add 和 Cart.Add 一樣把品項加進購物車，並在 Log 最後記一行「加入 <name> x<qty>」，
// 例如「加入 排骨飯 x2」。
func (c *LoggedCart) Add(name string, price, qty int) {
	// TODO: 實作
}

// Inc 把計數加 1。
//
// 找錯：呼叫了三次 Inc，Value 還是 0。
func (c Counter) Inc() {
	c.n++
}

// RaiseAll 把每位員工的薪水調高 pct%（小數無條件捨去）。
//
// 找錯：呼叫之後，薪水一個都沒變。
func RaiseAll(staff []Employee, pct int) {
	for _, e := range staff {
		e.Salary = e.Salary * (100 + pct) / 100
	}
}

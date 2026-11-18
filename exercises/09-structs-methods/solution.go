//go:build solution

package shop

import "fmt"

// DoubleBooked 用 Seat 當 map 的 key：Seat 的欄位都能比較，所以 Seat 也能比較。
func DoubleBooked(orders []Seat) (Seat, bool) {
	seen := map[Seat]bool{}
	for _, s := range orders {
		if seen[s] {
			return s, true
		}
		seen[s] = true
	}
	return Seat{}, false
}

// Add 要修改 c.Lines 裡的元素，所以用索引 c.Lines[i]，不能用 range 拿到的複本。
func (c *Cart) Add(name string, price, qty int) {
	for i := range c.Lines {
		if c.Lines[i].Name == name {
			c.Lines[i].Qty += qty
			return
		}
	}
	c.Lines = append(c.Lines, Line{Name: name, Price: price, Qty: qty})
}

func (c *Cart) Total() int {
	total := 0
	for _, l := range c.Lines {
		total += l.Price * l.Qty
	}
	return total
}

func (c *Card) TopUp(amount int) error {
	if amount <= 0 || amount%100 != 0 {
		return ErrBadAmount
	}
	c.Balance += amount
	return nil
}

func (c *Card) Pay(fare int) error {
	if fare > c.Balance {
		return ErrInsufficient
	}
	c.Balance -= fare
	return nil
}

// Add 第一次用到時才建立 map，所以零值的 Inventory 也能用。
func (inv *Inventory) Add(name string, n int) {
	if inv.stock == nil {
		inv.stock = map[string]int{}
	}
	inv.stock[name] += n
}

// Remove 讀取 nil map 不會出錯，所以不必先檢查 stock 是不是 nil。
func (inv *Inventory) Remove(name string, n int) error {
	if inv.stock[name] < n {
		return ErrInsufficient
	}
	inv.stock[name] -= n
	return nil
}

func (inv *Inventory) Count(name string) int {
	return inv.stock[name]
}

// Add 遮住了被嵌入的 Cart.Add，所以要寫完整的 c.Cart.Add 才會呼叫到它；
// 寫成 c.Add 會呼叫自己，永遠停不下來。
func (c *LoggedCart) Add(name string, price, qty int) {
	c.Cart.Add(name, price, qty)
	c.Log = append(c.Log, fmt.Sprintf("加入 %s x%d", name, qty))
}

// Inc 要用指標 receiver，才會改到呼叫者手上的 Counter。
func (c *Counter) Inc() {
	c.n++
}

// RaiseAll 用索引改 slice 裡的元素；range 的第二個變數只是複本。
func RaiseAll(staff []Employee, pct int) {
	for i := range staff {
		staff[i].Salary = staff[i].Salary * (100 + pct) / 100
	}
}

package main

import (
	"fmt"
	"time"
)

// Coupon 是一張有期限的折價券。
type Coupon struct {
	Code    string
	Expires time.Time
}

// Checkout 是結帳櫃台。現在幾點由 now 決定：正式執行時是 time.Now，測試時換成固定的時間。
type Checkout struct {
	now func() time.Time
}

// Accept 回傳折價券現在還能不能用：到期的那一刻起就不能用了。
func (c Checkout) Accept(cp Coupon) bool {
	return c.now().Before(cp.Expires)
}

// checkoutDemo 把櫃台的時鐘固定在兩個時間點，看同一張券的結果。
func checkoutDemo() {
	cp := Coupon{Code: "雙11", Expires: time.Date(2026, time.November, 12, 0, 0, 0, 0, taipei)}

	for _, at := range []time.Time{
		time.Date(2026, time.November, 11, 23, 59, 0, 0, taipei),
		time.Date(2026, time.November, 12, 0, 0, 0, 0, taipei),
	} {
		c := Checkout{now: func() time.Time { return at }}
		fmt.Println(at.Format("01/02 15:04"), cp.Code, "可以用？", c.Accept(cp))
	}
}

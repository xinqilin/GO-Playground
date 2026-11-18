package main

import (
	"errors"
	"fmt"
)

// StockError 記錄哪一樣東西不夠、還剩幾份，呼叫的人可以取出這些欄位。
type StockError struct {
	Item string
	Left int
}

func (e *StockError) Error() string {
	return fmt.Sprintf("%s 只剩 %d 份", e.Item, e.Left)
}

var stock = map[string]int{"排骨飯": 3, "雞腿飯": 0}

// reserve 檢查庫存夠不夠。不夠時回傳 *StockError；夠的話明確地回傳 nil。
func reserve(item string, qty int) error {
	if left := stock[item]; qty > left {
		return &StockError{Item: item, Left: left}
	}
	return nil
}

// OrderError 把錯誤和訂單號碼綁在一起。Unwrap 讓 errors.Is 和 errors.As 看得到裡面的 Err。
type OrderError struct {
	No  int
	Err error
}

func (e *OrderError) Error() string {
	return fmt.Sprintf("訂單 %d：%v", e.No, e.Err)
}

func (e *OrderError) Unwrap() error {
	return e.Err
}

// asDemo：從錯誤鏈裡找出 *StockError，取出它的欄位。
func asDemo() {
	err := &OrderError{No: 18, Err: reserve("排骨飯", 5)}
	fmt.Println(err)

	var se *StockError
	if errors.As(err, &se) {
		fmt.Println("As", se.Item, se.Left)
	}

	if se, ok := errors.AsType[*StockError](err); ok {
		fmt.Println("AsType", se.Item, se.Left)
	}
}

package main

import (
	"errors"
	"fmt"
)

// quote 報價：查價失敗時，用 %w 把原本的錯誤包起來，前面加上這一層知道的資訊。
func quote(name string, qty int) (int, error) {
	p, err := price(name)
	if err != nil {
		return 0, fmt.Errorf("報價 %s：%w", name, err)
	}
	return p * qty, nil
}

// placeOrder 下單：再包一層，加上訂單號碼。
func placeOrder(no int, name string, qty int) error {
	total, err := quote(name, qty)
	if err != nil {
		return fmt.Errorf("訂單 %d：%w", no, err)
	}
	fmt.Println("訂單", no, "合計", total)
	return nil
}

// wrapDemo：包了兩層之後，== 比不到，errors.Is 會一層一層拆開去找。
func wrapDemo() {
	err := placeOrder(17, "雞腿飯", 2)
	fmt.Println(err)
	fmt.Println(err == ErrSoldOut, errors.Is(err, ErrSoldOut))

	for e := err; e != nil; e = errors.Unwrap(e) {
		fmt.Printf("%T：%v\n", e, e)
	}
}

// wrapVsFormat：%w 和 %v 印出來一模一樣，但只有 %w 會保留錯誤鏈。
func wrapVsFormat() {
	_, err := price("雞腿飯")
	w := fmt.Errorf("報價：%w", err)
	v := fmt.Errorf("報價：%v", err)
	fmt.Println(w, errors.Is(w, ErrSoldOut))
	fmt.Println(v, errors.Is(v, ErrSoldOut))
}

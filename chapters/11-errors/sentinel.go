package main

import (
	"errors"
	"fmt"
)

// ErrSoldOut 和 ErrNotOnMenu 是 sentinel error：事先建立好的錯誤值，讓呼叫的人比對。
var (
	ErrSoldOut   = errors.New("賣完了")
	ErrNotOnMenu = errors.New("菜單上沒有")
)

type dish struct {
	price int
	left  int
}

var menu = map[string]dish{
	"排骨飯": {price: 85, left: 10},
	"雞腿飯": {price: 95, left: 0},
}

// price 查價。找不到或賣完時，回傳對應的 sentinel error。
func price(name string) (int, error) {
	d, ok := menu[name]
	if !ok {
		return 0, ErrNotOnMenu
	}
	if d.left == 0 {
		return 0, ErrSoldOut
	}
	return d.price, nil
}

// sentinelDemo：用 errors.Is 判斷是哪一種錯誤，對不同的錯誤做不同的處理。
func sentinelDemo() {
	for _, name := range []string{"排骨飯", "雞腿飯", "珍奶"} {
		p, err := price(name)
		switch {
		case err == nil:
			fmt.Println(name, p)
		case errors.Is(err, ErrSoldOut):
			fmt.Println(name, "今天賣完了，明天請早")
		default:
			fmt.Printf("%s 錯誤：%v\n", name, err)
		}
	}
}

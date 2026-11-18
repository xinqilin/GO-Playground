package main

import (
	"errors"
	"fmt"
)

// ErrBadPhone 表示電話號碼的格式不對。
var ErrBadPhone = errors.New("電話要 10 碼")

// Booking 是一筆訂位。
type Booking struct {
	Name   string
	People int
	Phone  string
}

// validate 一次檢查所有欄位，把每個問題都收集起來，而不是遇到第一個就回傳。
func validate(b Booking) error {
	var errs []error
	if b.Name == "" {
		errs = append(errs, errors.New("要留名字"))
	}
	if b.People < 1 || b.People > 12 {
		errs = append(errs, fmt.Errorf("人數 %d 不在 1 到 12 之間", b.People))
	}
	if len(b.Phone) != 10 {
		errs = append(errs, ErrBadPhone)
	}
	return errors.Join(errs...)
}

func joinDemo() {
	err := validate(Booking{People: 20, Phone: "0912"})
	fmt.Println(err)
	fmt.Println(errors.Is(err, ErrBadPhone))
	fmt.Println(validate(Booking{Name: "王小明", People: 4, Phone: "0912345678"}))
}

// Package bento 管理便當店的菜單。
//
// 對這一章的 main 來說，它是「別人寫的 package」：main 只能用這裡大寫開頭的名稱。
package bento

import (
	"errors"
	"maps"
	"slices"
)

// ErrBadPrice 表示價格不合理。
var ErrBadPrice = errors.New("價格要大於 0")

// Menu 是一份菜單。欄位沒有匯出，外面只能透過方法讀寫，所以價格不可能被改成負的。
type Menu struct {
	prices map[string]int
}

// New 建立一份空的菜單。
func New() *Menu {
	return &Menu{prices: make(map[string]int)}
}

// Add 加入一道菜；價格要大於 0。
func (m *Menu) Add(dish string, price int) error {
	if !valid(price) {
		return ErrBadPrice
	}
	m.prices[dish] = price
	return nil
}

// Price 回傳 dish 的價格；菜單上沒有時 ok 是 false。
func (m *Menu) Price(dish string) (price int, ok bool) {
	price, ok = m.prices[dish]
	return price, ok
}

// Dishes 回傳排序好的菜名。
func (m *Menu) Dishes() []string {
	return slices.Sorted(maps.Keys(m.prices))
}

// valid 沒有匯出：它是 bento 自己的實作細節，以後要改規則也不會影響別人。
func valid(price int) bool {
	return price > 0
}

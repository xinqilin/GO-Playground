package main

import (
	"errors"
	"fmt"
	"slices"
	"sync"
)

// 店裡的錯誤。handler 用 errors.Is 判斷，決定要回哪個狀態碼。
var (
	ErrNotFound    = errors.New("找不到這筆訂單")
	ErrSoldOut     = errors.New("賣完了")
	ErrUnknownDish = errors.New("沒有這道菜")
)

// Order 是一筆訂單。
type Order struct {
	ID    int    `json:"id"`
	Dish  string `json:"dish"`
	Qty   int    `json:"qty"`
	Total int    `json:"total"`
}

// MenuItem 是菜單上的一項。
type MenuItem struct {
	Dish  string `json:"dish"`
	Price int    `json:"price"`
	Left  int    `json:"left"`
}

// shop 把菜單、庫存和訂單放在記憶體裡。handler 會被很多 goroutine 同時呼叫，所以全部用一把鎖保護。
type shop struct {
	mu     sync.Mutex
	prices map[string]int
	stock  map[string]int
	orders map[int]Order
	nextID int
}

func newShop() *shop {
	return &shop{
		prices: map[string]int{"排骨飯": 90, "雞腿飯": 100, "滷肉飯": 50},
		stock:  map[string]int{"排骨飯": 5, "雞腿飯": 0, "滷肉飯": 20},
		orders: map[int]Order{},
		nextID: 1,
	}
}

func (s *shop) menu() []MenuItem {
	s.mu.Lock()
	defer s.mu.Unlock()
	var items []MenuItem
	for dish, price := range s.prices {
		items = append(items, MenuItem{Dish: dish, Price: price, Left: s.stock[dish]})
	}
	slices.SortFunc(items, func(a, b MenuItem) int { return a.Price - b.Price })
	return items
}

func (s *shop) place(dish string, qty int) (Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	price, ok := s.prices[dish]
	if !ok {
		return Order{}, fmt.Errorf("%w：%s", ErrUnknownDish, dish)
	}
	if s.stock[dish] < qty {
		return Order{}, fmt.Errorf("%s %w，剩 %d 份", dish, ErrSoldOut, s.stock[dish])
	}
	s.stock[dish] -= qty
	o := Order{ID: s.nextID, Dish: dish, Qty: qty, Total: price * qty}
	s.orders[o.ID] = o
	s.nextID++
	return o, nil
}

func (s *shop) find(id int) (Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.orders[id]
	if !ok {
		return Order{}, ErrNotFound
	}
	return o, nil
}

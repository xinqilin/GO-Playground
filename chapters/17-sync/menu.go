package main

import (
	"fmt"
	"sync"
)

// Menu 是價目表：查價的人很多、改價很少，所以用 RWMutex。
type Menu struct {
	mu     sync.RWMutex
	prices map[string]int
}

func NewMenu() *Menu {
	return &Menu{prices: make(map[string]int)}
}

// Price 只讀，用 RLock：好幾個 goroutine 可以同時查價。
func (m *Menu) Price(item string) (int, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.prices[item]
	return p, ok
}

// SetPrice 要寫，用 Lock：改價的時候，別人不能查，也不能同時改。
func (m *Menu) SetPrice(item string, price int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.prices[item] = price
}

func menuDemo() {
	m := NewMenu()
	m.SetPrice("排骨飯", 90)

	seen := make([]int, 50)
	var wg sync.WaitGroup
	for i := range seen {
		wg.Go(func() { seen[i], _ = m.Price("排骨飯") })
	}
	wg.Go(func() { m.SetPrice("排骨飯", 95) })
	wg.Wait()

	ok := true
	for _, p := range seen {
		if p != 90 && p != 95 {
			ok = false
		}
	}
	p, _ := m.Price("排骨飯")
	fmt.Printf("50 個人查到的不是 90 就是 95：%t\n", ok)
	fmt.Println("現在的價格", p)
}

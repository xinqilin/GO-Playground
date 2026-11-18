package main

import (
	"fmt"
	"iter"
)

// Menu 依照加入的順序記住品項和價格。零值就能用。
type Menu struct {
	names  []string
	prices map[string]int
}

func (m *Menu) Add(name string, price int) {
	if m.prices == nil {
		m.prices = map[string]int{}
	}
	if _, ok := m.prices[name]; !ok {
		m.names = append(m.names, name)
	}
	m.prices[name] = price
}

// All 依加入的順序交出品名和價格，是一個 iter.Seq2。
func (m *Menu) All() iter.Seq2[string, int] {
	return func(yield func(string, int) bool) {
		for _, n := range m.names {
			if !yield(n, m.prices[n]) {
				return
			}
		}
	}
}

func seq2Demo() {
	var m Menu
	m.Add("排骨飯", 85)
	m.Add("雞腿飯", 95)
	m.Add("紅茶", 25)
	for name, price := range m.All() {
		fmt.Println(name, price)
	}
}

// Category 是菜單的分類，底下可以再有子分類。
type Category struct {
	Name     string
	Children []*Category
}

// Walk 深度優先交出每一個分類的完整路徑。
func (c *Category) Walk() iter.Seq[string] {
	return func(yield func(string) bool) {
		c.walk("", yield)
	}
}

// walk 回傳 false 表示使用端已經喊停：每一層都要把 false 往上傳，立刻停下來。
func (c *Category) walk(prefix string, yield func(string) bool) bool {
	path := prefix + c.Name
	if !yield(path) {
		return false
	}
	for _, child := range c.Children {
		if !child.walk(path+"/", yield) {
			return false
		}
	}
	return true
}

func treeDemo() {
	menu := &Category{Name: "菜單", Children: []*Category{
		{Name: "便當", Children: []*Category{{Name: "排骨飯"}, {Name: "雞腿飯"}}},
		{Name: "飲料", Children: []*Category{{Name: "紅茶"}}},
	}}
	for path := range menu.Walk() {
		fmt.Println(path)
	}
	for path := range menu.Walk() {
		if path == "菜單/便當/雞腿飯" {
			fmt.Println("找到", path)
			break
		}
	}
}

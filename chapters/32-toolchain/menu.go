package main

import "fmt"

// 改了 menu.csv 之後，在這個目錄執行 go generate，menu_gen.go 就會重新產生。
//go:generate go run ./genmenu -in menu.csv -out menu_gen.go

// dish 是菜單上的一道菜。menu 這個變數在 menu_gen.go，由 genmenu 產生。
type dish struct {
	name  string
	price int
}

// menuDemo 印出產生出來的菜單。
func menuDemo() {
	for _, d := range menu {
		fmt.Printf("%s %d 元\n", d.name, d.price)
	}
}

package main

import (
	"fmt"
	"go/version"
)

// goVersionDemo 比較 Go 的版本字串，也就是 go.mod 裡 go 那一行的寫法。
func goVersionDemo() {
	pairs := [][2]string{
		{"go1.21.0", "go1.22rc1"},
		{"go1.22rc1", "go1.22.0"},
		{"go1.9", "go1.10"},
	}
	for _, p := range pairs {
		fmt.Printf("%s 比 %s：%d\n", p[0], p[1], version.Compare(p[0], p[1]))
	}
	fmt.Println(version.Lang("go1.27.4"), version.IsValid("1.27"), version.IsValid("go1.27"))
}

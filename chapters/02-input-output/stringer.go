package main

import (
	"fmt"
	"strconv"
)

// Size 是便當的份量。
type Size int

const (
	Small Size = iota
	Medium
	Large
)

// String 讓 fmt 印出 Size 時顯示中文。型別只要有這個方法，fmt 就會自動使用它。
func (s Size) String() string {
	switch s {
	case Small:
		return "小份"
	case Medium:
		return "中份"
	case Large:
		return "大份"
	}
	return "Size(" + strconv.Itoa(int(s)) + ")"
}

// stringerDemo 示範 Println、%v、%s 都會呼叫 String，%d 則印出底層的整數。
func stringerDemo() {
	s := Large
	fmt.Println(s)
	fmt.Printf("%v %s %d\n", s, s, s)
	fmt.Println(Size(7))
	fmt.Println([]Size{Small, Large})
}

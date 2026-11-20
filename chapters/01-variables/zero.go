package main

import "fmt"

// zeroValues 印出沒有給初始值時，各型別的零值。%q 會替字串加上引號，空字串才看得出來。
func zeroValues() {
	var i int
	var f float64
	var ok bool
	var s string
	var p *int
	fmt.Printf("int=%d float64=%g bool=%t string=%q pointer=%v\n", i, f, ok, s, p)
}

// zeroUsable 示範零值可以直接拿來用，不需要先「初始化」。
func zeroUsable() {
	var count int    // 0
	var order string // ""
	count++
	order += "排骨"
	count++
	order += "、雞腿"
	fmt.Printf("%d 個：%s\n", count, order)
}

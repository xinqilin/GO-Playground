package main

import "fmt"

// Gram（公克）和 Jin（台斤）是兩個新的型別，底層型別都是 float64。
type Gram float64

type Jin float64

// Grams 把台斤換算成公克：1 台斤是 600 公克。
func (j Jin) Grams() Gram {
	return Gram(j * 600)
}

// String 讓 Gram 被 fmt 印出來時帶著單位。
func (g Gram) String() string {
	return fmt.Sprintf("%.0f 公克", float64(g))
}

// definedTypes 示範新型別：可以有自己的方法，和底層型別之間要明確轉換。
func definedTypes() {
	pork := Jin(1.5)
	fmt.Println(pork)
	fmt.Println(pork.Grams())

	var raw float64 = 300
	rice := Gram(raw)
	fmt.Println(rice + 50)
	fmt.Println(float64(rice))
}

// Score 是 int 的別名：它就是 int，不是新的型別。
type Score = int

// Points 是新的型別，底層型別是 int。
type Points int

// aliasDemo 比較別名和新型別：%T 印出的型別名稱不同。
func aliasDemo() {
	var s Score = 90
	var p Points = 90
	var n int = s
	fmt.Printf("%T %T %d\n", s, p, n)
}

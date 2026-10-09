// 第 01 章：變數、常數與型別。
//
// 變數是一個貼了名字、只能放某一種東西的盒子；沒放東西時，盒子裡是那個型別的「零值」。
// 執行：go run ./chapters/01-variables
package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("== 宣告 ==")
	declare()

	fmt.Println("== 零值 ==")
	zeroValues()

	fmt.Println("== 常數與 iota ==")
	constants()

	fmt.Println("== 型別轉換 ==")
	conversions()

	fmt.Println("== 底線 _ ==")
	blankIdentifier()
}

// declare 示範三種宣告方式。
func declare() {
	var a int = 10 // 完整寫法：名稱、型別、初始值
	var b = 20     // 省略型別，由右邊的值推斷
	c := 30        // 簡短宣告，只能用在函式裡
	a, b = b, a    // 一行交換兩個變數
	fmt.Println(a, b, c)
}

// zeroValues 示範沒有給初始值時，各型別的零值。
func zeroValues() {
	var i int
	var f float64
	var ok bool
	var s string
	var p *int
	fmt.Printf("int=%d float64=%g bool=%t string=%q pointer=%v\n", i, f, ok, s, p)
}

// Weekday 是一個自訂型別。Go 沒有 enum 關鍵字，慣例是用 const 搭配 iota。
type Weekday int

const (
	Sunday  Weekday = iota // 0
	Monday                 // 1：沒寫的話沿用上一行的寫法，iota 自動加 1
	Tuesday                // 2
)

// constants 示範常數與 iota。
func constants() {
	const greeting = "你好" // 常數在編譯時就決定，之後不能修改
	const big = 1 << 100  // 沒有型別的常數精度很高，只要最後的結果放得進去就好
	fmt.Println(greeting, big>>98)
	fmt.Println(Sunday, Monday, Tuesday)
}

// conversions 示範型別轉換：Go 不會自動轉，一定要明確寫出來。
func conversions() {
	i := 65
	f := float64(i) / 2 // int 和 float64 不能直接混著算
	fmt.Println(f)
	fmt.Println(string(rune(i)), strconv.Itoa(i)) // rune 轉字串是 "A"，數字轉文字要用 strconv

	var small int8 = 127
	small++ // 超出範圍會繞回去，不會報錯
	fmt.Println(small)
}

// blankIdentifier 示範用 _ 丟掉不需要的值。
func blankIdentifier() {
	_, value, _ := strings.Cut("lang=go", "=") // 只要等號後面的部分
	fmt.Println(value)
}

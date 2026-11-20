package main

import (
	"fmt"
	"math"
)

// intRanges 印出幾種整數型別的範圍。math 套件有每一種型別的上下限常數。
func intRanges() {
	fmt.Printf("%-6s %d ~ %d\n", "int8", math.MinInt8, math.MaxInt8)
	fmt.Printf("%-6s %d ~ %d\n", "uint8", 0, math.MaxUint8)
	fmt.Printf("%-6s %d ~ %d\n", "int16", math.MinInt16, math.MaxInt16)
	fmt.Printf("%-6s %d ~ %d\n", "int32", math.MinInt32, math.MaxInt32)
	fmt.Printf("%-6s %d ~ %d\n", "int64", int64(math.MinInt64), int64(math.MaxInt64))
	fmt.Printf("%-6s %d ~ %d\n", "uint64", 0, uint64(math.MaxUint64))
}

// overflowWrap 示範整數在執行時溢位會繞一圈，不會報錯。
func overflowWrap() {
	var level int8 = 127
	level++
	fmt.Println(level)

	var stock uint8 = 0
	stock--
	fmt.Println(stock)
}

// integerDivision 示範整數除法會捨去小數（往 0 的方向），% 的正負號跟著被除數。
func integerDivision() {
	fmt.Println(7/2, 7%2)
	fmt.Println(-7/2, -7%2)
	fmt.Println(7.0 / 2)

	total, people := 1000, 3
	fmt.Println("每人", total/people, "元，剩", total%people, "元")
}

// divide 回傳 a / b。b 是 0 時會 panic。
func divide(a, b int) int {
	return a / b
}

// divideByZero 示範整數除以 0 會在執行時 panic。
// 這裡用 defer 和 recover 把 panic 接住才能印出訊息，兩者在第 04 章介紹。
func divideByZero() {
	defer func() {
		fmt.Println("panic:", recover())
	}()
	fmt.Println(divide(10, 0))
}

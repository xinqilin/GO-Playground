package main

import (
	"fmt"
	"strconv"
)

// binary 把非負整數寫成二進位字串。遞迴的兩個部分：
// 終止條件（n 小於 2 時直接回答），以及把問題縮小後呼叫自己。
func binary(n int) string {
	if n < 2 {
		return strconv.Itoa(n)
	}
	return binary(n/2) + strconv.Itoa(n%2)
}

// recursion 示範遞迴。
func recursion() {
	fmt.Println(binary(13), binary(255), binary(0))
}

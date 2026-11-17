package main

import (
	"fmt"
	"math"
)

// seat 是座位的位置。欄位都可以比較，所以 seat 可以當 map 的 key。
type seat struct {
	row, col int
}

// structKey 示範 struct 和陣列都能當 key，只要它們可以用 == 比較。
func structKey() {
	booked := map[seat]string{
		{1, 2}: "小明",
		{3, 4}: "小華",
	}
	fmt.Println(booked[seat{1, 2}], len(booked))
	_, taken := booked[seat{2, 2}]
	fmt.Println(taken)

	trips := make(map[[2]string]int) // 起點和終點組成的路段
	trips[[2]string{"甲站", "乙站"}]++
	trips[[2]string{"甲站", "乙站"}]++
	trips[[2]string{"乙站", "甲站"}]++
	fmt.Println(trips[[2]string{"甲站", "乙站"}], len(trips))
}

// nanKey 示範 NaN 不等於自己，所以用 NaN 當 key 放進去的值永遠查不到。
func nanKey() {
	m := make(map[float64]string)
	m[math.NaN()] = "a"
	m[math.NaN()] = "b"
	_, ok := m[math.NaN()]
	fmt.Println(len(m), ok)
}

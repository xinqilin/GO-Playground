package main

import (
	"fmt"
	"slices"
)

// slicesPackage 示範 slices 套件常用的函式。
func slicesPackage() {
	prices := []int{90, 85, 120, 85, 70}
	fmt.Println(slices.Contains(prices, 120), slices.Index(prices, 85), slices.Index(prices, 999))
	fmt.Println(slices.Min(prices), slices.Max(prices))

	sorted := slices.Clone(prices)
	slices.Sort(sorted)
	fmt.Println(sorted)
	i, found := slices.BinarySearch(sorted, 90)
	fmt.Println(i, found)

	sorted = slices.Compact(sorted) // 移除相鄰的重複值
	fmt.Println(sorted)
	slices.Reverse(sorted)
	fmt.Println(sorted)
	fmt.Println(slices.Equal(sorted, []int{120, 90, 85, 70}))
}

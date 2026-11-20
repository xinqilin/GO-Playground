package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

// stdlibDemo：slices、maps、strings 都提供回傳疊代器的函式。
func stdlibDemo() {
	stock := map[string]int{"雞腿飯": 0, "排骨飯": 10, "滷蛋": 25}
	fmt.Println(slices.Sorted(maps.Keys(stock)))

	queue := []int{17, 18, 19}
	for i, n := range slices.Backward(queue) {
		fmt.Println(i, n)
	}

	receipt := "排骨飯 x2\n紅茶 x1\n合計 195\n"
	for line := range strings.Lines(receipt) {
		fmt.Printf("%q\n", line)
	}

	for batch := range slices.Chunk([]int{1, 2, 3, 4, 5}, 2) {
		fmt.Println(batch)
	}
}

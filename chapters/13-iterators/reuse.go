package main

import (
	"bufio"
	"fmt"
	"slices"
	"strings"
)

// reuseDemo：有些疊代器每次 range 都從頭來過，有些只能用一次。
func reuseDemo() {
	seq := callNumbers(1, 3)
	fmt.Println(slices.Collect(seq), slices.Collect(seq))

	sc := bufio.NewScanner(strings.NewReader("排骨飯\n紅茶\n"))
	once := func(yield func(string) bool) {
		for sc.Scan() {
			if !yield(sc.Text()) {
				return
			}
		}
	}
	fmt.Println(slices.Collect(once), slices.Collect(once))
}

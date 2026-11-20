package main

import (
	"fmt"
	"io"
	"maps"
	"os"
	"slices"

	"github.com/xinqilin/GO-Playground/chapters/33-idiomatic-go/ticket"
)

// waiter 是 busiest 和 printBoard 唯一需要的能力。interface 定義在使用它的地方，
// ticket package 完全不知道有這個 interface。
type waiter interface {
	Waiting() int
}

// busiest 回傳等最多人的櫃台；人數相同時，取名稱排序在前的那一個。
func busiest(counters map[string]waiter) string {
	best, most := "", -1
	for _, name := range slices.Sorted(maps.Keys(counters)) {
		if n := counters[name].Waiting(); n > most {
			best, most = name, n
		}
	}
	return best
}

// printBoard 把每個櫃台的等候人數寫到 w，依櫃台名稱排序。
func printBoard(w io.Writer, counters map[string]waiter) error {
	for _, name := range slices.Sorted(maps.Keys(counters)) {
		if _, err := fmt.Fprintf(w, "%s櫃台：%d 人在等\n", name, counters[name].Waiting()); err != nil {
			return err
		}
	}
	return nil
}

func interfaceDemo() {
	var pay, pickup ticket.Machine
	for range 3 {
		pay.Take()
	}
	pickup.Take()
	counters := map[string]waiter{"結帳": &pay, "取餐": &pickup}
	printBoard(os.Stdout, counters)
	fmt.Printf("最忙的是：%s\n", busiest(counters))
}

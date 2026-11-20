package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

// sink 讓編譯器不能把示範裡的結果當成沒用到而整個省略。
var (
	sink      string
	sinkBytes []byte
)

// allocDemo 用 testing.AllocsPerRun 數出「排骨飯 x2 180」的四種組法各配置幾次記憶體。
func allocDemo() {
	dish, qty, price := "排骨飯", 2, 180
	buf := make([]byte, 0, 64) // 重複使用的緩衝區

	ways := []struct {
		name string
		f    func()
	}{
		{"fmt.Sprintf", func() {
			sink = fmt.Sprintf("%s x%d %d", dish, qty, price)
		}},
		{"字串相加", func() {
			sink = dish + " x" + strconv.Itoa(qty) + " " + strconv.Itoa(price)
		}},
		{"strings.Builder", func() {
			var b strings.Builder
			b.Grow(32)
			b.WriteString(dish)
			b.WriteString(" x")
			b.WriteString(strconv.Itoa(qty))
			b.WriteString(" ")
			b.WriteString(strconv.Itoa(price))
			sink = b.String()
		}},
		{"append 到舊的 buf", func() {
			b := append(buf[:0], dish...)
			b = append(b, " x"...)
			b = strconv.AppendInt(b, int64(qty), 10)
			b = append(b, ' ')
			b = strconv.AppendInt(b, int64(price), 10)
			sinkBytes = b
		}},
	}
	for _, w := range ways {
		fmt.Printf("%s：每次配置 %.0f 次\n", w.name, testing.AllocsPerRun(100, w.f))
	}
}

// collectNames 一邊讀一邊 append，slice 不夠大時就換一塊更大的記憶體。
func collectNames(orders []Order) []string {
	var names []string
	for _, o := range orders {
		names = append(names, o.Dish)
	}
	return names
}

// collectNamesPrealloc 事先知道要放幾個，一開始就配置好。
func collectNamesPrealloc(orders []Order) []string {
	names := make([]string, 0, len(orders))
	for _, o := range orders {
		names = append(names, o.Dish)
	}
	return names
}

// preallocDemo 比較兩種寫法最後的 len 和 cap。配置了幾次要看 BenchmarkCollectNames。
func preallocDemo() {
	orders := make([]Order, 1000)
	a, b := collectNames(orders), collectNamesPrealloc(orders)
	fmt.Printf("沒有預先配置：len %d cap %d\n", len(a), cap(a))
	fmt.Printf("預先配置：len %d cap %d\n", len(b), cap(b))
}

package main

import (
	"fmt"
	"strconv"
	"strings"
)

// receipt 用 strings.Builder 組出一張收據。Builder 在同一塊空間往後寫，最後才變成字串。
func receipt(items []string, prices []int) string {
	var b strings.Builder
	total := 0
	for i, item := range items {
		fmt.Fprintf(&b, "%s %d\n", item, prices[i])
		total += prices[i]
	}
	b.WriteString("合計 ")
	b.WriteString(strconv.Itoa(total))
	return b.String()
}

// builderDemo 示範使用 receipt。
func builderDemo() {
	fmt.Println(receipt([]string{"排骨", "雞腿"}, []int{85, 90}))
}

// concatPlus 用 + 串接：每一次都配置一個新字串，再把目前為止的內容整個複製過去。
func concatPlus(parts []string) string {
	s := ""
	for _, p := range parts {
		s += p
	}
	return s
}

// concatBuilder 用 Builder 串接，並先用 Grow 保留需要的空間，整個過程只配置一次。
func concatBuilder(parts []string) string {
	n := 0
	for _, p := range parts {
		n += len(p)
	}
	var b strings.Builder
	b.Grow(n)
	for _, p := range parts {
		b.WriteString(p)
	}
	return b.String()
}

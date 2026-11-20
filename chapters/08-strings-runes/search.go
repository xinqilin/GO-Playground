package main

import (
	"fmt"
	"strings"
)

// searchFuncs 示範在字串裡找東西。Index 回傳的是 byte 位置。
func searchFuncs() {
	addr := "台北市大安區便當路 100 號"
	fmt.Println(strings.Contains(addr, "大安區"), strings.HasPrefix(addr, "台北市"), strings.HasSuffix(addr, "號"))
	fmt.Println(strings.Index(addr, "大安"), strings.Index(addr, "信義"))
	fmt.Println(strings.Count("排骨,雞腿,排骨", "排骨"))
	fmt.Println(strings.EqualFold("Go", "GO"))
}

// splitJoin 示範切開與接起來。用 %q 印出，才看得出空字串和空白。
func splitJoin() {
	fmt.Printf("%q\n", strings.Split("排骨,雞腿,,素食", ","))
	fmt.Printf("%q\n", strings.Fields("  排骨  雞腿\t素食\n"))
	fmt.Println(strings.Join([]string{"排骨", "雞腿", "素食"}, "、"))
}

// trimCut 示範去掉頭尾，以及用 Cut 切成兩半。
func trimCut() {
	fmt.Printf("%q\n", strings.TrimSpace("  排骨飯 \n"))
	fmt.Println(strings.TrimPrefix("INV-AB12345678", "INV-"), strings.TrimSuffix("menu.txt", ".txt"))
	fmt.Println(strings.Trim("**特價**", "*"))
	name, qty, ok := strings.Cut("排骨x2", "x")
	fmt.Println(name, qty, ok)
	_, _, ok = strings.Cut("排骨", "x")
	fmt.Println(ok)
}

// replaceFuncs 示範取代、重複和大小寫。
func replaceFuncs() {
	fmt.Println(strings.ReplaceAll("0912-345-678", "-", ""))
	fmt.Println(strings.Replace("a-b-c", "-", "+", 1))
	fmt.Println(strings.Repeat("=", 10))
	fmt.Println(strings.ToUpper("go 語言"), strings.ToLower("HELLO"))
	r := strings.NewReplacer("台", "臺", "-", "")
	fmt.Println(r.Replace("台北-台中"))
}

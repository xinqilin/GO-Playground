package main

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"unicode/utf8"
)

type dish struct {
	name  string
	stock int
}

// rangeCopy 示範 range 給的元素是複本：改複本不會改到 slice 本身。
func rangeCopy() {
	menu := []dish{{"雞腿飯", 3}, {"排骨飯", 5}}
	for _, d := range menu {
		d.stock = 0 // 改到的是複本 d
	}
	fmt.Println(menu)

	for i := range menu {
		menu[i].stock = 0 // 透過索引才改得到 menu 裡的元素
	}
	fmt.Println(menu)
}

// rangeArrayCopy 比較 range 陣列和 range slice：陣列會在迴圈開始前整個複製一份。
func rangeArrayCopy() {
	arr := [3]string{"空", "空", "空"}
	var seen []string
	for i, s := range arr {
		if i == 0 {
			arr[2] = "有人"
		}
		seen = append(seen, s)
	}
	fmt.Println("陣列：", seen)

	sl := []string{"空", "空", "空"}
	seen = nil
	for i, s := range sl {
		if i == 0 {
			sl[2] = "有人"
		}
		seen = append(seen, s)
	}
	fmt.Println("slice：", seen)
}

// rangeString 走過字串時，每一圈拿到的是一個 rune（字元）和它的起始 byte 位置。
func rangeString() {
	s := "Go語言"
	for i, r := range s {
		fmt.Printf("位置 %d：%c（%U）\n", i, r, r)
	}
	fmt.Printf("len = %d，字元數 = %d\n", len(s), utf8.RuneCountInString(s))
}

// rangeMap 先把 key 排序再走，輸出順序才會固定；直接 range map 的順序是隨機的。
func rangeMap() {
	stock := map[string]int{"排骨": 5, "雞腿": 3, "魚排": 0}
	for _, name := range slices.Sorted(maps.Keys(stock)) {
		fmt.Println(name, stock[name])
	}
}

// rangeChannel 一直從 channel 收值，直到 channel 被 close 而且收完為止。
func rangeChannel() {
	tickets := make(chan string, 3)
	tickets <- "A01"
	tickets <- "A02"
	tickets <- "A03"
	close(tickets) // 沒有 close，range 會一直等下去
	for t := range tickets {
		fmt.Println("叫號", t)
	}
}

// countdown 回傳一個可以被 range 的函式，第 13 章會詳細介紹。
func countdown(from int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for n := from; n > 0; n-- {
			if !yield(n) {
				return
			}
		}
	}
}

func rangeFunc() {
	for n := range countdown(3) {
		fmt.Println(n)
	}
	fmt.Println("發射")
}

// rangeOnce 示範 range 後面的運算式只在開始時求值一次，迴圈裡 append 不會讓圈數變多。
func rangeOnce() {
	nums := []int{1, 2, 3}
	count := 0
	for _, n := range nums {
		nums = append(nums, n*10)
		count++
	}
	fmt.Println("跑了", count, "圈，nums =", nums)
}

// loopVarPerIteration 比較「每一圈都是新變數」和「大家共用一個變數」。
func loopVarPerIteration() {
	var perIter, shared []func() int
	for i := range 3 {
		perIter = append(perIter, func() int { return i }) // 每一圈的 i 都是新的
	}

	var j int
	for j = 0; j < 3; j++ {
		shared = append(shared, func() int { return j }) // j 宣告在迴圈外，三個 closure 共用它
	}

	for k := range 3 {
		fmt.Println(perIter[k](), shared[k]())
	}
}

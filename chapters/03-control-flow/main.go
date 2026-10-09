// 第 03 章：流程控制。
//
// Go 的流程控制只有三樣：if、for、switch。沒有 while，也沒有三元運算子。
// 執行：go run ./chapters/03-control-flow
package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("== if 可以先做一件事 ==")
	ifWithInit("42")
	ifWithInit("abc")

	fmt.Println("== for 的四種寫法 ==")
	loops()

	fmt.Println("== 九九乘法表（前三列） ==")
	table(3)

	fmt.Println("== 有標籤的 break ==")
	findPair([]int{3, 8, 1, 6}, 9)

	fmt.Println("== switch ==")
	for _, score := range []int{95, 72, 40} {
		fmt.Println(score, grade(score))
	}
	fmt.Println(dayType("六"), dayType("二"))

	fmt.Println("== fallthrough ==")
	fallthroughDemo("a")
}

// ifWithInit 示範 if 前面可以先放一個簡短敘述，宣告的變數只活在 if/else 裡。
func ifWithInit(s string) {
	if n, err := strconv.Atoi(s); err != nil {
		fmt.Println(s, "不是數字")
	} else {
		fmt.Println(s, "的兩倍是", n*2)
	}
}

// loops 示範 for 的四種寫法。
func loops() {
	var got []int
	for i := 0; i < 3; i++ { // 1. 三段式
		got = append(got, i)
	}
	fmt.Println(got)

	n := 1
	for n < 100 { // 2. 只寫條件，就是其他語言的 while
		n *= 3
	}
	fmt.Println(n)

	got = got[:0]
	for i := range 3 { // 3. Go 1.22 起可以直接 range 一個整數
		got = append(got, i)
	}
	fmt.Println(got)

	count := 0
	for { // 4. 什麼都不寫是無窮迴圈，用 break 離開
		count++
		if count == 5 {
			break
		}
	}
	fmt.Println(count)
}

// table 用兩層 for 印出 n×n 的乘法表。
func table(n int) {
	for i := 1; i <= n; i++ {
		row := make([]string, 0, n)
		for j := 1; j <= n; j++ {
			row = append(row, fmt.Sprintf("%d×%d=%d", j, i, i*j))
		}
		fmt.Println(strings.Join(row, "  "))
	}
}

// findPair 找第一組加起來等於 target 的數字，找到就用標籤一次跳出兩層迴圈。
func findPair(nums []int, target int) {
outer:
	for i := range nums {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				fmt.Println("找到", nums[i], nums[j])
				break outer
			}
		}
	}
}

// grade 用沒有條件的 switch 取代一長串 if-else，由上往下找第一個成立的 case。
func grade(score int) string {
	switch {
	case score >= 90:
		return "優"
	case score >= 60:
		return "及格"
	default:
		return "不及格"
	}
}

// dayType 示範一個 case 可以列多個值。
func dayType(day string) string {
	switch day {
	case "六", "日":
		return "週末"
	default:
		return "平日"
	}
}

// fallthroughDemo 示範 fallthrough：不檢查下一個 case 的條件，直接執行它。實務上很少用。
func fallthroughDemo(ans string) {
	switch ans {
	case "a":
		fmt.Println("ans 是 a")
		fallthrough
	case "b":
		fmt.Println("進到 b 的區塊")
	default:
		fmt.Println("沒有答案")
	}
}

package main

import "fmt"

// labeledBreak 在內層迴圈找到答案後，用 break 加標籤一次跳出兩層。
func labeledBreak() {
	prices := []int{45, 55, 35, 65}
	budget := 100
search:
	for i := range prices {
		for j := i + 1; j < len(prices); j++ {
			if prices[i]+prices[j] == budget {
				fmt.Println("第", i, "和第", j, "個加起來剛好", budget, "元")
				break search
			}
		}
	}
}

// labeledContinue 一發現不及格就跳到下一班，不再看這一班剩下的分數。
func labeledContinue() {
	classes := [][]int{{88, 92, 75}, {90, 55, 80}, {60, 70, 65}}
nextClass:
	for c, scores := range classes {
		for _, s := range scores {
			if s < 60 {
				fmt.Println("第", c+1, "班有人不及格：", s)
				continue nextClass
			}
		}
		fmt.Println("第", c+1, "班全部及格")
	}
}

// countUntilStopBuggy 想在遇到 "stop" 時停止，但 switch 裡的 break 只會跳出 switch。
func countUntilStopBuggy(cmds []string) int {
	handled := 0
	for _, c := range cmds {
		switch c {
		case "stop":
			break
		default:
			handled++
		}
	}
	return handled
}

// countUntilStop 幫迴圈加上標籤，break loop 才會離開 for。
func countUntilStop(cmds []string) int {
	handled := 0
loop:
	for _, c := range cmds {
		switch c {
		case "stop":
			break loop
		default:
			handled++
		}
	}
	return handled
}

func breakInSwitch() {
	cmds := []string{"開燈", "開冷氣", "stop", "關門"}
	fmt.Println("沒有標籤：", countUntilStopBuggy(cmds))
	fmt.Println("有標籤：", countUntilStop(cmds))
}

// reverseBuggy 用無號整數倒著數：i >= 0 永遠成立，減到 0 之後會繞回最大值，接著索引越界。
func reverseBuggy(s []string) []string {
	var out []string
	for i := uint(len(s)) - 1; i >= 0; i-- {
		out = append(out, s[i])
	}
	return out
}

// reverse 用 int 倒著數，i 變成 -1 時條件不成立，迴圈正常結束。
func reverse(s []string) []string {
	var out []string
	for i := len(s) - 1; i >= 0; i-- {
		out = append(out, s[i])
	}
	return out
}

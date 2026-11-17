package main

import (
	"fmt"
	"os"
	"strings"
)

// printSpacing 示範 Print 和 Println 補空白的規則不一樣。
func printSpacing() {
	fmt.Print("便當", 3, 85, "元", "\n")
	fmt.Println("便當", 3, 85, "元")
}

// sprintDemo 示範 S 開頭的函式不印出來，而是回傳組好的字串。
func sprintDemo() {
	label := fmt.Sprintf("%s x%d", "排骨飯", 2)
	total := fmt.Sprint(170)
	fmt.Println(label + "：" + total + " 元")
}

// fprintDemo 示範 F 開頭的函式寫到指定的 io.Writer：可以是 os.Stdout，也可以是 strings.Builder。
func fprintDemo() {
	fmt.Fprintln(os.Stdout, "寫到標準輸出")

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d 號請到 %d 號櫃台", 42, 3)
	fmt.Println("Builder 裡有：" + sb.String())
}

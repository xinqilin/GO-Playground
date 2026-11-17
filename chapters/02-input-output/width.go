package main

import "fmt"

// widthBasics 示範寬度、對齊、補零和精度。| 標出欄位的邊界。
func widthBasics() {
	fmt.Printf("|%6d|%-6d|%06d|\n", 85, 85, 85)
	fmt.Printf("|%8.2f|%-8.2f|%.1f|\n", 3.14159, 3.14159, 3.14159)
	fmt.Printf("|%6s|%-6s|%.2s|\n", "Go", "Go", "Golang")
	fmt.Printf("|%*d|\n", 6, 85)
}

// fareTable 印出票價表。數字欄放前面並靠右對齊，中文放在最後一欄，就不會被字寬影響。
// 標題的「票價」兩個中文字在終端機上剛好佔 4 格，和 %4d 一樣寬。
func fareTable() {
	fmt.Println("票價  目的地")
	fmt.Printf("%4d  %s\n", 20, "台北車站")
	fmt.Printf("%4d  %s\n", 35, "南港展覽館")
	fmt.Printf("%4d  %s\n", 65, "淡水")
}

// widthCJK 示範 %-6s 是用「字元數」補空白，但中文字在終端機上佔兩格寬，所以右邊的 | 對不齊。
func widthCJK() {
	fmt.Printf("|%-6s|\n", "Go")
	fmt.Printf("|%-6s|\n", "排骨飯")
}

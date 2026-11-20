// forever 示範忘了更新變數造成的無窮迴圈：cups 永遠是 0，條件永遠成立。
// 它會一直印同一行，要按 Ctrl-C 才會停，所以沒有測試。
//
// 執行：go run ./chapters/pre-3-decisions/forever
package main

import "fmt"

func main() {
	cups := 0
	for cups < 3 {
		fmt.Println("倒一杯，目前有", cups, "杯")
	}
	fmt.Println("倒好了")
}

// deadlock 示範最常見的死結：在同一個 goroutine 裡，對無緩衝的 channel 送值，卻沒有人在收。
// 這支程式一定會當掉，所以放在獨立的目錄，沒有測試。
//
// 執行：go run ./chapters/15-channels/deadlock
package main

import "fmt"

func main() {
	window := make(chan string)
	window <- "排骨飯" // 沒有其他 goroutine 會來收：main 永遠卡在這一行
	fmt.Println(<-window)
}

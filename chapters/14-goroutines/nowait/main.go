// nowait 示範 main 不會等 goroutine：main 一結束，整個程式就結束了，還沒跑完的 goroutine 也一起消失。
// 輸出取決於時間，所以這支程式沒有測試。
//
// 執行：go run ./chapters/14-goroutines/nowait
package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("排骨飯做好了")
	}()
	fmt.Println("點完餐，走人")
}

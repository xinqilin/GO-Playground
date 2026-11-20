// crash 示範 goroutine 裡的 panic 會讓整個程式結束，main 裡的 recover 接不到。
// 這支程式一定會當掉，所以放在獨立的目錄，沒有測試。
//
// 執行：go run ./chapters/14-goroutines/crash
package main

import (
	"fmt"
	"sync"
)

func main() {
	defer func() {
		fmt.Println("main 接到了：", recover())
	}()

	var wg sync.WaitGroup
	wg.Add(1) // wg.Go 的文件規定傳進去的函式不可以 panic，所以這裡用 Add 和 Done
	go func() {
		var stock map[string]int
		stock["排骨飯"]-- // 寫入 nil map：panic
		wg.Done()
	}()
	wg.Wait()
	fmt.Println("不會印出來")
}

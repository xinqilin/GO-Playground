// procs 印出這台機器的 CPU 數、GOMAXPROCS 和目前的 goroutine 數。
// 輸出隨機器而不同，所以這支程式沒有測試。
//
// 執行：go run ./chapters/14-goroutines/procs
package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("CPU 數", runtime.NumCPU())
	fmt.Println("GOMAXPROCS", runtime.GOMAXPROCS(0))
	fmt.Println("goroutine 數", runtime.NumGoroutine())
}

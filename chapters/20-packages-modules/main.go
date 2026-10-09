// 第 20 章：package 與 module。
//
// module 是一整個專案（一個 go.mod），package 是裡面的一個資料夾。
// import 路徑 = module 路徑 + 資料夾路徑。
// 執行：go run ./chapters/20-packages-modules
package main

import (
	"fmt"

	"github.com/xinqilin/GO-Playground/chapters/20-packages-modules/greet"
	"github.com/xinqilin/GO-Playground/chapters/20-packages-modules/internal/counter"
)

// init 在 main 之前自動執行，每個檔案可以有多個。適合做簡單的初始化，不要放太多邏輯。
func init() {
	fmt.Println("init：在 main 之前執行")
}

func main() {
	fmt.Println(greet.Hello("Gopher"))
	fmt.Println(greet.Hello("  "))
	// greet.normalize("x") 會編譯錯誤：normalize 沒有匯出。

	greet.Prefix = "嗨" // 匯出的變數，其他套件也改得到
	fmt.Println(greet.Hello("Go"))

	fmt.Println(counter.Next(), counter.Next())
}

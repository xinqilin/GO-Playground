// 第 24 章：工具鏈。
//
// go 指令是一整套工具：格式化、檢查、文件、測試、交叉編譯、效能分析都在裡面。
// 執行：go run ./chapters/24-toolchain
//
//	go run -tags pro ./chapters/24-toolchain
//	go build -ldflags "-X main.version=1.2.3" -o bin/app ./chapters/24-toolchain
package main

import (
	"fmt"
	"runtime"
)

// version 預設是 dev；發版時用 -ldflags "-X main.version=1.2.3" 在編譯時改掉。
var version = "dev"

// describe 回傳版本與版別。edition 來自 edition_free.go 或 edition_pro.go，由 build tag 決定編哪一個。
func describe() string {
	return fmt.Sprintf("版本 %s，%s", version, edition)
}

func main() {
	fmt.Println(describe())
	// 交叉編譯時，GOOS 和 GOARCH 決定產生哪個平台的執行檔。
	fmt.Println("平台", runtime.GOOS+"/"+runtime.GOARCH)
}

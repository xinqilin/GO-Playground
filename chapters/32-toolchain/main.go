// 第 32 章：工具鏈。
//
// 這一章的程式本身是一個小小的「便當店點餐機」，重點在怎麼用 go 指令建置它：
// build tag 決定版別、-ldflags 注入版本、GOOS/GOARCH 交叉編譯、go generate 產生菜單。
// main_test.go 用 Example 鎖住不隨建置方式改變的輸出；toolchain_test.go 實際執行
// go build、go list、go fix，確認頁面上的指令和輸出是真的。
//
// 執行：go run ./chapters/32-toolchain
//
//	go run -tags pro ./chapters/32-toolchain
//	go build -ldflags "-X main.version=1.2.3" -o app ./chapters/32-toolchain
package main

import (
	"fmt"
	"runtime/debug"
)

// version 預設是 dev；發版時用 -ldflags "-X main.version=1.2.3" 在編譯時改掉。
var version = "dev"

// describe 回傳版本與版別。edition 來自 edition_free.go 或 edition_pro.go，由 build tag 決定編哪一個。
func describe() string {
	return fmt.Sprintf("版本 %s，%s", version, edition)
}

// buildInfoDemo 印出 go build 自動寫進執行檔的建置資訊。內容取決於怎麼建置，所以沒有 Example。
func buildInfoDemo() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("沒有建置資訊")
		return
	}
	fmt.Printf("Go 版本：%s\n", info.GoVersion)
	fmt.Printf("module 版本：%s\n", info.Main.Version)
	for _, s := range info.Settings {
		switch s.Key {
		case "-tags", "-ldflags", "GOOS", "GOARCH", "vcs.revision", "vcs.modified":
			fmt.Printf("%s = %s\n", s.Key, s.Value)
		}
	}
}

func main() {
	section("1. 版本與版別")
	fmt.Println(describe())
	fmt.Printf("打開網頁用的指令：%v\n", openCommand())

	section("2. 建置資訊")
	buildInfoDemo()

	section("3. 菜單（go generate 產生）")
	menuDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}

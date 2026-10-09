// 第 21 章：測試。
//
// 測試是自動檢查員：每次改完程式，跑一次 go test，它就幫你把該檢查的地方全部再看一遍。
// 這一章的重點在 textutil/textutil_test.go。
// 執行：go test -v ./chapters/21-testing/...
package main

import (
	"fmt"

	"github.com/xinqilin/GO-Playground/chapters/21-testing/textutil"
)

func main() {
	fmt.Println(textutil.Reverse("Go語言"))
	fmt.Println(textutil.IsPalindrome("Never odd or even"))
}

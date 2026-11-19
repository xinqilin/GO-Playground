// args 把 os.Args 的每一個元素印出來，看 shell 把命令列切成了哪幾段。
//
// 執行（先編譯，os.Args[0] 才會是簡短的 ./args）：
//
//	go build -o args ./chapters/21-cli/args
//	./args total -tax 10 "週一 營業額.txt"
package main

import (
	"fmt"
	"os"
)

func main() {
	for i, a := range os.Args {
		fmt.Printf("os.Args[%d] = %q\n", i, a)
	}
}

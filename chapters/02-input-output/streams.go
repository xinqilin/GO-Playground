package main

import (
	"fmt"
	"os"
)

// streams 示範 stdout 和 stderr 是兩條不同的管道。
// Example 只比對 stdout，所以 stderr 那一行不會出現在測試的輸出裡。
func streams() {
	fmt.Fprintln(os.Stderr, "提示：這一行寫到 stderr")
	fmt.Fprintln(os.Stdout, "結果：這一行寫到 stdout")
}

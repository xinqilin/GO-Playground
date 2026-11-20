// 和 workspace/shop 一樣的程式，這次靠 go.mod 裡的 replace 找到 example.com/price。
package main

import (
	"fmt"

	"example.com/price"
)

func main() {
	fmt.Printf("排骨飯 %d 元\n", price.Of("排骨飯"))
}

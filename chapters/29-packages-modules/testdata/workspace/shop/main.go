// shop 用到 example.com/price。這個 module 沒有發佈到任何地方，是 go.work 讓 go 找得到它。
package main

import (
	"fmt"

	"example.com/price"
)

func main() {
	for _, dish := range []string{"排骨飯", "滷肉飯"} {
		fmt.Printf("%s %d 元\n", dish, price.Of(dish))
	}
}

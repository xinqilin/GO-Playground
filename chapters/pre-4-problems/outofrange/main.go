// outofrange 示範索引超出範圍：迴圈條件寫成 <=，多跑了一圈。
// 它故意會當掉，所以沒有測試。
//
// 執行：go run ./chapters/pre-4-problems/outofrange
package main

import "fmt"

func main() {
	menu := []string{"排骨飯", "雞腿飯", "素食餐"}
	for i := 0; i <= len(menu); i++ {
		fmt.Println(i, menu[i])
	}
}

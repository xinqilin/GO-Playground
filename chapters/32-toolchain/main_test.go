package main

import "fmt"

// 版本和版別取決於建置方式；這裡是不加任何參數時的結果。
func Example_describe() {
	fmt.Println(describe())
	// Output: 版本 dev，免費版
}

func Example_menuDemo() {
	menuDemo()
	// Output:
	// 排骨飯 90 元
	// 雞腿飯 100 元
	// 滷肉飯 50 元
	// 魚排飯 95 元
	// 素食便當 80 元
}

package main

import (
	"fmt"
	"strings"

	"github.com/xinqilin/GO-Playground/chapters/29-packages-modules/internal/steps"
)

// 套件層級的變數依照「相依關係」初始化：total 寫在最前面，
// 但它用到 price 和 qty，所以要等它們兩個都有值之後才輪到它。
var total = steps.Record("total", price*qty)
var price = steps.Record("price", 90)
var qty = steps.Record("qty", 2)

// init 在所有套件層級的變數都初始化完之後執行，一個檔案可以有好幾個。
func init() {
	steps.Record("main 的 init", 0)
}

// initOrderDemo 印出 main 開始之前發生的事。
func initOrderDemo() {
	fmt.Println(strings.Join(steps.Log, " → "))
	fmt.Println("total =", total)
}

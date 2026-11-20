package main

import (
	"log"
	"os"
)

// logDemo 用 log.New 建立自己的 Logger：寫到哪裡、前綴、旗標都自己決定。
// 這裡寫到 stdout、旗標用 0（不印日期時間），輸出才每次一樣；
// 套件層級的 log.Println 寫到 stderr，旗標是 log.LstdFlags（日期和時間）。
func logDemo() {
	kitchen := log.New(os.Stdout, "[廚房] ", 0)
	kitchen.Println("開火")
	kitchen.Printf("第 %d 桌 %s 出餐", 7, "排骨飯") // 結尾沒有換行，log 會自動補上

	counter := log.New(kitchen.Writer(), "[外場] ", kitchen.Flags())
	counter.Println("叫號", 42)
}

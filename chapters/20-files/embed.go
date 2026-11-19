package main

import (
	"embed"
	"fmt"
	"io/fs"
)

// 下面兩個變數的內容在編譯時就放進執行檔：部署時只要帶一個執行檔，不用另外帶資料夾。
// //go:embed 寫在變數的正上方，路徑相對於這個 .go 檔所在的目錄。

//go:embed welcome.txt
var welcome string

//go:embed receipts
var receipts embed.FS

// embedDemo 讀 embed 進來的檔案。receipts 是一個 fs.FS，可以直接交給 receiptTotal。
func embedDemo() {
	fmt.Print(welcome)

	total, err := receiptTotal(receipts)
	fmt.Println("embed 進來的發票合計：", total, err)

	names, err := fs.Glob(receipts, "receipts/*/*.txt")
	fmt.Println(names, err)
}

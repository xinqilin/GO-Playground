package main

import "fmt"

// deferOrder 示範 defer 在函式結束前才執行，後註冊的先執行。
func deferOrder() {
	defer fmt.Println("關店")
	defer fmt.Println("關燈")
	defer fmt.Println("收銀機結帳")
	fmt.Println("營業中")
}

// deferArgs 示範 defer 的參數在註冊當下就算好；closure 則在真正執行時才讀變數。
func deferArgs() {
	x := 10
	defer fmt.Println("參數：", x)
	defer func() { fmt.Println("closure：", x) }()
	x = 20
}

// deferResult 示範 defer 可以修改具名回傳值：return 先把結果放進 total，defer 再改它。
func deferResult() (total int) {
	defer func() {
		total += 10 // 外送費
	}()
	return 85
}

// open 和 closeRes 模擬需要關閉的資源，例如檔案。
func open(name string) string {
	fmt.Println("開啟", name)
	return name
}

func closeRes(name string) {
	fmt.Println("關閉", name)
}

// deferInLoop 的 defer 要等整個函式結束才執行，所以前面的資源一直沒有關閉。
func deferInLoop() {
	for _, name := range []string{"a.txt", "b.txt"} {
		f := open(name)
		defer closeRes(f)
		fmt.Println("處理", f)
	}
	fmt.Println("迴圈結束")
}

// deferInLoopFixed 把每一圈的工作包成一個函式，defer 在那個函式結束時就執行。
func deferInLoopFixed() {
	for _, name := range []string{"a.txt", "b.txt"} {
		func() {
			f := open(name)
			defer closeRes(f)
			fmt.Println("處理", f)
		}()
	}
	fmt.Println("迴圈結束")
}

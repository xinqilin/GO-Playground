package main

import "fmt"

// take 從庫存拿貨。拿太多是呼叫端的程式寫錯了，所以直接 panic。
func take(stock, n int) int {
	if n > stock {
		panic(fmt.Sprintf("庫存不足：剩 %d 個，要拿 %d 個", stock, n))
	}
	return stock - n
}

// safeTake 在邊界把 panic 轉成 error：recover 只有在 defer 的函式裡直接呼叫才有效。
func safeTake(stock, n int) (left int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("拿貨失敗：%v", r)
		}
	}()
	return take(stock, n), nil
}

// panicRecover 示範 recover 接住 panic 之後，程式可以繼續往下走。
func panicRecover() {
	fmt.Println(safeTake(5, 2))
	fmt.Println(safeTake(3, 10))
	fmt.Println("程式繼續執行")
}

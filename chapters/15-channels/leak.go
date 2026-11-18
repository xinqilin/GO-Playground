package main

import "fmt"

// findLeaky 在出餐流裡找 want，回傳它是第幾份（從 0 開始）。
//
// 錯誤示範：找到就 return，不再收了。送的那個 goroutine 卡在下一次送出，永遠不會結束。
func findLeaky(items []string, want string) int {
	i := 0
	for o := range orders(items...) {
		if o == want {
			return i
		}
		i++
	}
	return -1
}

// find 找到之後，由 defer 把剩下的收完，讓送的一方能送完、close、結束。
func find(items []string, want string) int {
	ch := orders(items...)
	defer func() {
		for range ch {
		}
	}()
	i := 0
	for o := range ch {
		if o == want {
			return i
		}
		i++
	}
	return -1
}

func findDemo() {
	fmt.Println(find([]string{"排骨飯", "雞腿飯", "魚排飯"}, "雞腿飯"))
}

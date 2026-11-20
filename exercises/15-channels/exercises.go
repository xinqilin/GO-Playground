//go:build !solution

package relay

import "fmt"

// closed 回傳一個已經 close 的 channel，題目檔用它當暫時的回傳值。
func closed[T any]() <-chan T {
	ch := make(chan T)
	close(ch)
	return ch
}

// Tickets 是發號碼牌的機器：回傳一個 channel，在另一個 goroutine 裡依序送出 from、from+1……to，送完就 close。
// from 大於 to 時什麼都不送，直接 close。
func Tickets(from, to int) <-chan int {
	// TODO: 實作
	return closed[int]()
}

// Collect 收 ch 裡的每一個值，直到 ch 被 close，依照收到的順序回傳。
func Collect[T any](ch <-chan T) []T {
	// TODO: 實作
	return nil
}

// Stage 是管線的一站：從 in 收值，把 f 的結果送到回傳的 channel；in 被 close 之後，也 close 回傳的 channel。
func Stage[T, U any](in <-chan T, f func(T) U) <-chan U {
	// TODO: 實作
	return closed[U]()
}

// FanIn 把好幾個 channel 收到的值全部送到同一個 channel。
// 所有的輸入都被 close、而且值都送完之後，才 close 回傳的 channel。順序不限。
func FanIn[T any](chs ...<-chan T) <-chan T {
	// TODO: 實作
	return closed[T]()
}

// FirstOK 同時執行每一個 try，回傳第一個成功（err 是 nil）的結果，不必等其他的做完。
// 全部失敗時，把所有錯誤依照 tries 的順序用 errors.Join 合起來回傳。
// 不可以留下卡住的 goroutine：FirstOK 回傳之後，還在跑的 try 做完時也要能順利結束。
func FirstOK(tries ...func() (string, error)) (string, error) {
	// TODO: 實作
	return "", nil
}

// Receipts 回傳一個 channel，依序送出每一筆金額的收據文字，例如 85 是 "NT$85"。
//
// 找錯：值都收得到，但 range 這個 channel 的人永遠等不到結束。
func Receipts(amounts []int) <-chan string {
	out := make(chan string)
	go func() {
		for _, a := range amounts {
			out <- fmt.Sprintf("NT$%d", a)
		}
	}()
	return out
}

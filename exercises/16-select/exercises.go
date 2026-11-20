//go:build !solution

package desk

import "time"

// TryTake 不等待地從 ch 收一個值：有值就回傳它和 true；
// 沒有值、或 ch 已經 close 而且收完了，就立刻回傳零值和 false。
func TryTake[T any](ch <-chan T) (T, bool) {
	// TODO: 實作
	var zero T
	return zero, false
}

// Within 最多等 d，看 ch 有沒有送來一個值：等到了回傳它和 true；
// 時間到了、或 ch 被 close，回傳零值和 false。
func Within[T any](ch <-chan T, d time.Duration) (T, bool) {
	// TODO: 實作
	var zero T
	return zero, false
}

// NewBell 建立一個服務鈴。
func NewBell() *Bell {
	// TODO: 實作
	return &Bell{}
}

// Ring 按一下鈴。不管有沒有人在聽、之前按過幾次，都不可以卡住。
// 店員還沒來看之前，按很多次也只算一次。
func (b *Bell) Ring() {
	// TODO: 實作
}

// C 回傳店員要聽的 channel：有人按過鈴，就收得到一個值。
func (b *Bell) C() <-chan struct{} {
	return b.ch
}

// Tally 收完 sales 和 refunds 兩個 channel，回傳營業額總和減掉退款總和。
// 兩個 channel 都被 close 才回傳；傳進來的 channel 是 nil，就當作它沒有任何值。
func Tally(sales, refunds <-chan int) int {
	// TODO: 實作
	return 0
}

// Counter 回傳一個 channel，依序送出 1、2、3……沒有上限。
// done 被 close 之後就停止，並 close 回傳的 channel；不可以留下卡住的 goroutine。
func Counter(done <-chan struct{}) <-chan int {
	// TODO: 實作
	out := make(chan int)
	close(out)
	return out
}

// Debounce 等輸入安靜下來才送出：每收到一個值就重新計時，
// 連續 wait 這麼久都沒有新的值，才把最後收到的那個值送出去。
// in 被 close 時，如果還有等著送的值，立刻送出，然後 close 回傳的 channel。
func Debounce[T any](in <-chan T, wait time.Duration) <-chan T {
	// TODO: 實作
	out := make(chan T)
	close(out)
	return out
}

// CollectFor 從 ch 收值，直到 ch 被 close，或從呼叫開始算起已經過了 d，回傳收到的值。
//
// 找錯：來源每隔一陣子就有新的值時，它永遠不會因為時間到而停下來。
func CollectFor[T any](ch <-chan T, d time.Duration) []T {
	var got []T
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return got
			}
			got = append(got, v)
		case <-time.After(d):
			return got
		}
	}
}

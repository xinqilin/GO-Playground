//go:build solution

package desk

import "time"

// TryTake 用 default：沒有任何 case 能進行時，select 走 default，不等待。
func TryTake[T any](ch <-chan T) (T, bool) {
	select {
	case v, ok := <-ch:
		return v, ok
	default:
		var zero T
		return zero, false
	}
}

func Within[T any](ch <-chan T, d time.Duration) (T, bool) {
	select {
	case v, ok := <-ch:
		return v, ok
	case <-time.After(d):
		var zero T
		return zero, false
	}
}

// NewBell 用容量 1 的 channel：第一次按鈴放進一個值，之後再按，channel 已滿，走 default 直接離開。
func NewBell() *Bell {
	return &Bell{ch: make(chan struct{}, 1)}
}

func (b *Bell) Ring() {
	select {
	case b.ch <- struct{}{}:
	default:
	}
}

func (b *Bell) C() <-chan struct{} {
	return b.ch
}

// Tally 在某一邊 close 之後把它設成 nil，select 就不會再選它；兩邊都是 nil 時結束。
// 一開始就是 nil 的 channel，也自然被當成已經收完。
func Tally(sales, refunds <-chan int) int {
	total := 0
	for sales != nil || refunds != nil {
		select {
		case v, ok := <-sales:
			if !ok {
				sales = nil
				continue
			}
			total += v
		case v, ok := <-refunds:
			if !ok {
				refunds = nil
				continue
			}
			total -= v
		}
	}
	return total
}

// Counter 每次送出時同時等 done，done 被 close 就結束；defer 負責 close out。
func Counter(done <-chan struct{}) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := 1; ; n++ {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}()
	return out
}

// Debounce 用一個 nil 的計時 channel 代表「目前沒有等著送的值」：
// nil channel 永遠收不到東西，那個 case 就像不存在一樣。
func Debounce[T any](in <-chan T, wait time.Duration) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		var pending T
		var fire <-chan time.Time
		for {
			select {
			case v, ok := <-in:
				if !ok {
					if fire != nil {
						out <- pending
					}
					return
				}
				pending = v
				fire = time.After(wait)
			case <-fire:
				out <- pending
				fire = nil
			}
		}
	}()
	return out
}

// CollectFor 在迴圈外面只建立一次計時器，d 才是從呼叫開始算起的整段期限。
// 寫在 select 裡的 time.After(d) 每一圈都會建立新的計時器，等於每收到一個值就重新計時。
func CollectFor[T any](ch <-chan T, d time.Duration) []T {
	var got []T
	deadline := time.After(d)
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return got
			}
			got = append(got, v)
		case <-deadline:
			return got
		}
	}
}

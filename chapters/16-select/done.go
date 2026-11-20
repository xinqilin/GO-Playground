package main

import "fmt"

// fromUntil 和 from 一樣依序送出 items，但每次送的時候也同時等 done：
// done 被 close，就不再送了，直接結束。
func fromUntil(done <-chan struct{}, items ...string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for _, it := range items {
			select {
			case out <- it:
			case <-done:
				return
			}
		}
	}()
	return out
}

// find 找到 want 之後 return，defer close(done) 通知送的一方停下來，不必把剩下的收完。
func find(items []string, want string) int {
	done := make(chan struct{})
	defer close(done)
	i := 0
	for o := range fromUntil(done, items...) {
		if o == want {
			return i
		}
		i++
	}
	return -1
}

func findDemo() {
	fmt.Println(find([]string{"排骨飯", "雞腿飯", "魚排飯", "滷肉飯"}, "雞腿飯"))
}

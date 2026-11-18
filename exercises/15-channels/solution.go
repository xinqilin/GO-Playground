//go:build solution

package relay

import (
	"errors"
	"fmt"
	"sync"
)

func Tickets(from, to int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := from; n <= to; n++ {
			out <- n
		}
	}()
	return out
}

func Collect[T any](ch <-chan T) []T {
	var out []T
	for v := range ch {
		out = append(out, v)
	}
	return out
}

func Stage[T, U any](in <-chan T, f func(T) U) <-chan U {
	out := make(chan U)
	go func() {
		defer close(out)
		for v := range in {
			out <- f(v)
		}
	}()
	return out
}

// FanIn 每個輸入一個 goroutine 負責轉送；另外一個 goroutine 等全部轉送完才 close。
func FanIn[T any](chs ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup
	for _, ch := range chs {
		wg.Go(func() {
			for v := range ch {
				out <- v
			}
		})
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

type outcome struct {
	i   int
	val string
	err error
}

// FirstOK 的 channel 容量等於 try 的數量：就算 FirstOK 已經回傳、沒有人再收，
// 每個 goroutine 都還有位子放結果，送完就結束，不會卡住。
func FirstOK(tries ...func() (string, error)) (string, error) {
	results := make(chan outcome, len(tries))
	for i, try := range tries {
		go func() {
			v, err := try()
			results <- outcome{i, v, err}
		}()
	}
	errs := make([]error, len(tries))
	for range tries {
		r := <-results
		if r.err == nil {
			return r.val, nil
		}
		errs[r.i] = r.err
	}
	return "", errors.Join(errs...)
}

// Receipts 送完之後 close，range 的一方才知道沒有了。
func Receipts(amounts []int) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for _, a := range amounts {
			out <- fmt.Sprintf("NT$%d", a)
		}
	}()
	return out
}

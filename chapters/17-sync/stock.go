package main

import (
	"fmt"
	"sync"
)

// Stock 是某一道菜還剩幾份，可以給多個 goroutine 同時使用。
// mu 保護 left：讀或寫 left 之前都要先拿到 mu。零值的 Mutex 就能直接用。
type Stock struct {
	mu   sync.Mutex
	left int
}

// Take 賣出一份：還有就減一並回傳 true，賣完了回傳 false。
// 「檢查還有沒有」和「減一」在同一段上鎖的區間裡，中間不會有別人插進來。
func (s *Stock) Take() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.left == 0 {
		return false
	}
	s.left--
	return true
}

// Left 回傳還剩幾份。只是讀，也要上鎖。
func (s *Stock) Left() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.left
}

// set 把份數改成 n。
func (s *Stock) set(n int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.left = n
}

// takeSplit 是錯誤的寫法：先用 Left 查、再用 set 改，兩次各自上鎖。
// 每一次讀寫都有鎖，race detector 不會抱怨，但兩次之間別人可以插進來。
func (s *Stock) takeSplit() bool {
	n := s.Left()
	if n == 0 {
		return false
	}
	s.set(n - 1)
	return true
}

func stockDemo() {
	s := &Stock{left: 30}
	sold := make([]bool, 100)
	var wg sync.WaitGroup
	for i := range sold {
		wg.Go(func() { sold[i] = s.Take() })
	}
	wg.Wait()

	n := 0
	for _, ok := range sold {
		if ok {
			n++
		}
	}
	fmt.Printf("100 個人搶 30 份：賣出 %d 份，剩 %d 份\n", n, s.Left())
}

// oversellDemo 把兩個客人同時用 takeSplit 的步驟排出來：兩個人都先查到還有 1 份，再各自扣掉。
func oversellDemo() {
	s := &Stock{left: 1}
	a := s.Left() // 客人 A 查：還有 1 份
	b := s.Left() // 客人 B 也查：還有 1 份
	if a > 0 {
		s.set(a - 1)
		fmt.Println("A 買到了")
	}
	if b > 0 {
		s.set(b - 1)
		fmt.Println("B 買到了")
	}
	fmt.Println("只有 1 份，卻賣出 2 份，剩", s.Left())
}

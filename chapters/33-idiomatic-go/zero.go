package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/xinqilin/GO-Playground/chapters/33-idiomatic-go/ticket"
)

// zeroDemo：宣告完就能用，不需要 New。
func zeroDemo() {
	var m ticket.Machine
	m.Take()
	m.Take()
	fmt.Printf("叫號機：%d 人在等\n", m.Waiting())

	var b strings.Builder
	b.WriteString("排骨飯")
	b.WriteString(" x2")
	fmt.Printf("strings.Builder：%s\n", b.String())

	var mu sync.Mutex
	mu.Lock()
	mu.Unlock()
	fmt.Println("sync.Mutex：零值是沒有上鎖的鎖")

	var orders []string
	orders = append(orders, "滷肉飯")
	fmt.Printf("nil slice 也能 append：%v\n", orders)

	var prices map[string]int
	fmt.Printf("nil map 可以讀：%d，寫入會 panic\n", prices["排骨飯"])
}

// calledDemo：Called 回傳複本，呼叫端怎麼改都不影響叫號機。
func calledDemo() {
	var m ticket.Machine
	for range 3 {
		m.Take()
	}
	m.Call()
	m.Call()
	got := m.Called()
	got[0] = 99
	fmt.Printf("改過的複本：%v\n", got)
	fmt.Printf("叫號機裡的：%v\n", m.Called())
}

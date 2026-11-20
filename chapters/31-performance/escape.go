package main

import "fmt"

// Ticket 是一張號碼牌。
type Ticket struct {
	No   int
	Note string
}

// newTicket 回傳指標：t 在函式結束後還要被用到，只能放在 heap。
func newTicket(no int) *Ticket {
	t := Ticket{No: no}
	return &t
}

// ticketValue 回傳值：複製一份給呼叫的人，t 本身留在 stack。
func ticketValue(no int) Ticket {
	t := Ticket{No: no}
	return t
}

// announce 把 no 交給 fmt.Println。參數型別是 any，編譯器看不出它會被留多久。
func announce(no int) {
	fmt.Println("請", no, "號到櫃台")
}

// lastFive 的陣列大小固定，而且沒有離開函式，留在 stack。
func lastFive(nos []int) int {
	var recent [5]int
	copy(recent[:], nos[max(len(nos)-5, 0):])
	sum := 0
	for _, n := range recent {
		sum += n
	}
	return sum
}

// queueOf 的大小要到執行時才知道，只能放在 heap。
func queueOf(n int) []int {
	return make([]int, n)
}

// escapeDemo 呼叫上面的函式。哪些變數放在 heap，要用 go build -gcflags=-m 看。
func escapeDemo() {
	p := newTicket(7)
	v := ticketValue(8)
	announce(p.No + v.No)
	fmt.Println(lastFive([]int{1, 2, 3, 4, 5, 6}), len(queueOf(3)))
}

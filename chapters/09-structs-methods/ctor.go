package main

import (
	"errors"
	"fmt"
)

// Queue 是排隊叫號機。零值就能直接使用：號碼從 0 開始。
type Queue struct {
	next    int // 最後發出的號碼
	serving int // 目前叫到的號碼
}

// Take 發一張號碼牌。
func (q *Queue) Take() int {
	q.next++
	return q.next
}

// Call 叫下一號；沒有人在等時回傳 0, false。
func (q *Queue) Call() (int, bool) {
	if q.serving >= q.next {
		return 0, false
	}
	q.serving++
	return q.serving, true
}

// Waiting 回傳還有幾個人在等。
func (q *Queue) Waiting() int {
	return q.next - q.serving
}

func zeroValueDemo() {
	var q Queue
	a := q.Take()
	b := q.Take()
	fmt.Println(a, b, q.Waiting())

	n, ok := q.Call()
	fmt.Println(n, ok, q.Waiting())

	q.Call()
	n, ok = q.Call()
	fmt.Println(n, ok)
}

// Reservation 是一筆訂位。
type Reservation struct {
	Name   string
	People int
}

// NewReservation 先檢查輸入，不合理的訂位不會被建立出來。
func NewReservation(name string, people int) (*Reservation, error) {
	if name == "" {
		return nil, errors.New("訂位要留名字")
	}
	if people < 1 || people > 12 {
		return nil, fmt.Errorf("人數 %d 不在 1 到 12 之間", people)
	}
	return &Reservation{Name: name, People: people}, nil
}

func constructorDemo() {
	r, err := NewReservation("王先生", 4)
	fmt.Println(r.Name, r.People, err)

	_, err = NewReservation("林小姐", 20)
	fmt.Println(err)

	_, err = NewReservation("", 2)
	fmt.Println(err)
}

// Tally 統計每道菜賣了幾份。零值就能用：第一次 Add 時才建立 map。
type Tally struct {
	counts map[string]int
}

// Add 記一份。
func (t *Tally) Add(name string) {
	if t.counts == nil {
		t.counts = map[string]int{}
	}
	t.counts[name]++
}

// Count 回傳賣了幾份。讀取 nil map 不會出錯，會得到零值。
func (t *Tally) Count(name string) int {
	return t.counts[name]
}

func tallyDemo() {
	var t Tally
	fmt.Println(t.Count("排骨飯"))
	t.Add("排骨飯")
	t.Add("排骨飯")
	t.Add("雞腿飯")
	fmt.Println(t.Count("排骨飯"), t.Count("雞腿飯"))
}

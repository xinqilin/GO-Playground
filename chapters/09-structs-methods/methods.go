package main

import "fmt"

// StampCard 是集點卡，集滿 10 點換一個便當。
type StampCard struct {
	Stamps int
}

// StampCopy 是值 receiver：s 是呼叫者那張卡的複本，蓋章蓋在複本上。
func (s StampCard) StampCopy() {
	s.Stamps++
}

// Stamp 是指標 receiver：s 指向呼叫者那張卡，蓋章蓋在原本的卡上。
func (s *StampCard) Stamp() {
	s.Stamps++
}

// Free 只讀不改。
func (s *StampCard) Free() bool {
	return s.Stamps >= 10
}

func receiverDemo() {
	card := StampCard{}
	card.StampCopy()
	fmt.Println(card.Stamps)
	card.Stamp()
	card.Stamp()
	fmt.Println(card.Stamps)
}

// autoAddress 變數可以取地址時，Go 自動補上 & 或 *。
func autoAddress() {
	card := StampCard{Stamps: 9}
	p := &card
	p.Stamp()
	card.Stamp()
	fmt.Println(p.Free(), card.Stamps)
}

// Ticket 的 String 是指標 receiver。
type Ticket struct {
	No int
}

// String 讓 *Ticket 被 fmt 印出來時顯示號碼。
func (t *Ticket) String() string {
	return fmt.Sprintf("第 %d 號", t.No)
}

// methodSetDemo：String 掛在 *Ticket 上，所以只有傳指標給 Println 時才會被用到。
func methodSetDemo() {
	t := Ticket{No: 7}
	fmt.Println(t)
	fmt.Println(&t)
}

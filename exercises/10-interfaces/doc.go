// Package kiosk 是第 10 章「interface」的練習題。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/10-interfaces/
//
// 只想跑某一題時加上 -run：
//
//	go test ./exercises/10-interfaces/ -run TestAsFloat
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/10-interfaces/
package kiosk

import (
	"errors"
	"fmt"
	"io"
)

// Shipper 是一種寄送方式。Fee 回傳寄 grams 公克要多少運費；寄不了時回傳 -1。
type Shipper interface {
	Fee(grams int) int
}

// StorePickup 是超商取貨：5000 公克以內 60 元，超過就寄不了。
type StorePickup struct{}

func (StorePickup) Fee(grams int) int {
	if grams > 5000 {
		return -1
	}
	return 60
}

// HomeDelivery 是宅配。
type HomeDelivery struct{}

// Masker 把寫進來的內容轉寫到 W，但每個 ASCII 數字都換成 '*'。
type Masker struct {
	W io.Writer
}

// Item 是購物車裡的一樣東西。
type Item interface {
	Price() int
}

// Discounter 是可以折價的東西，Discount 回傳折多少元。
type Discounter interface {
	Discount() int
}

// Bento 是便當，沒有折扣。
type Bento struct {
	Name string
	Cost int
}

func (b Bento) Price() int { return b.Cost }

// Combo 是套餐，可以折 Off 元。
type Combo struct {
	Name string
	Cost int
	Off  int
}

func (c Combo) Price() int    { return c.Cost }
func (c Combo) Discount() int { return c.Off }

// Sender 是一種送出訊息的方式，送不出去時回傳 error。
type Sender interface {
	Send(msg string) error
}

// Fallback 依序嘗試每一種 Sender，直到有一種成功。Fallback 自己也是一種 Sender。
type Fallback []Sender

// ErrNoSender 表示 Fallback 裡一種 Sender 都沒有。
var ErrNoSender = errors.New("沒有可用的傳送方式")

// Order 是一張訂單。
type Order struct {
	Items int
	Paid  bool
}

// OrderError 說明訂單哪裡有問題。
type OrderError struct {
	Reason string
}

func (e *OrderError) Error() string { return "訂單錯誤：" + e.Reason }

// Price 是金額，印出來時帶著「元」。
type Price int

func (p Price) String() string { return fmt.Sprintf("%d 元", int(p)) }

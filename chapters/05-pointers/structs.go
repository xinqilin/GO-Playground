package main

import "fmt"

// structPointer 示範指向 struct 的指標可以直接用 . 存取欄位，Go 會自動解參考。
func structPointer() {
	b := bento{"雞腿", 90}
	p := &b
	p.price = 95 // 等於 (*p).price = 95
	fmt.Println(b.price, (*p).price)
}

// newVsLiteral 示範兩種取得「指向新變數的指標」的方法。
func newVsLiteral() {
	p := new(bento) // 配置一個零值的 bento，回傳它的位址
	p.name = "素食"
	q := &bento{name: "素食"} // 複合字面值加上 &，可以順便填欄位
	fmt.Println(*p == *q, p == q)
}

// newWithValue 示範 new 直接接一個值（Go 1.26 起），常用在需要「指向某個值的指標」的欄位。
func newWithValue() {
	discount := new(15)
	fmt.Println(*discount)
	o := order{id: 3, note: new("少飯")}
	fmt.Println(*o.note)
}

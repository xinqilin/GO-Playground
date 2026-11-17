package main

import "fmt"

// addOneValue 收到的是複本，改了也不影響呼叫端。
func addOneValue(n int) { n++ }

// addOnePointer 收到的是位址，改的是原本那一格。
func addOnePointer(n *int) { *n++ }

// passValueVsPointer 示範 Go 一律傳值：想讓函式改到呼叫端的變數，就傳指標。
func passValueVsPointer() {
	a, b := 1, 1
	addOneValue(a)
	addOnePointer(&b)
	fmt.Println("a =", a, "b =", b)
}

type bento struct {
	name  string
	price int
}

// raiseCopy 收到整個 struct 的複本。
func raiseCopy(b bento) { b.price += 5 }

// raisePointer 收到 struct 的位址。
func raisePointer(b *bento) { b.price += 5 }

// structCopy 示範 struct 也是值：指定或傳參數都會整個複製一份。
func structCopy() {
	pork := bento{"排骨", 85}
	copied := pork
	copied.price = 999
	fmt.Println(pork.price, copied.price)

	raiseCopy(pork)
	fmt.Println(pork.price)
	raisePointer(&pork)
	fmt.Println(pork.price)
}

type order struct {
	id   int
	note *string
}

// shallowCopy 示範複製 struct 時，指標欄位複製的是位址：兩份 struct 指向同一個字串。
func shallowCopy() {
	note := "不要香菜"
	first := order{id: 1, note: &note}
	second := first
	second.id = 2
	*second.note = "加辣"
	fmt.Println(first.id, *first.note)
	fmt.Println(second.id, *second.note)
}

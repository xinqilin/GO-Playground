package main

import "fmt"

// interfaceValue：interface 值裡裝著兩樣東西，實際的型別和實際的值。
func interfaceValue() {
	var n Notifier
	fmt.Println(n == nil)
	fmt.Printf("%T\n", n)

	n = SMS{Phone: "0912-345-678"}
	fmt.Printf("%T %v\n", n, n)

	n = Email{Addr: "amei@example.com"}
	fmt.Printf("%T %v\n", n, n)
	fmt.Println(n == nil)
}

// Buzzer 是取餐呼叫器，記得自己響了幾次。Notify 要修改 count，所以用指標 receiver。
type Buzzer struct {
	No    int
	count int
}

func (b *Buzzer) Notify(msg string) {
	b.count++
	fmt.Printf("呼叫器 %d 號響了（第 %d 次）：%s\n", b.No, b.count, msg)
}

// pointerInInterface：放進 interface 的是 &b，透過 interface 呼叫改到的就是 b。
func pointerInInterface() {
	b := Buzzer{No: 5}
	var n Notifier = &b
	n.Notify("餐點好了")
	n.Notify("請到櫃台取餐")
	fmt.Println(b.count)
}

// compareInterfaces：interface 可以用 == 比較，但裝著不能比較的值時會在執行期 panic。
func compareInterfaces() {
	var a, b Notifier = SMS{Phone: "0912"}, SMS{Phone: "0912"}
	var c Notifier = Email{Addr: "0912"}
	fmt.Println(a == b, a == c)

	defer func() {
		fmt.Println("panic:", recover())
	}()
	var x, y any = []string{"排骨"}, []string{"排骨"}
	fmt.Println(x == y)
}

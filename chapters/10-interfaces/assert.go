package main

import (
	"errors"
	"fmt"
)

// assertDemo：type assertion 從 interface 取出具體的值；帶 ok 的寫法不會 panic。
func assertDemo() {
	var n Notifier = SMS{Phone: "0912-345-678"}

	s := n.(SMS)
	fmt.Println(s.Phone)

	e, ok := n.(Email)
	fmt.Println(ok, e)
}

// assertPanic：不帶 ok 而型別又不對，就會 panic。
func assertPanic() {
	defer func() {
		fmt.Println("panic:", recover())
	}()
	var n Notifier = SMS{Phone: "0912-345-678"}
	e := n.(Email)
	fmt.Println(e.Addr)
}

// Urgent 是額外的能力：有些通知方式可以標成緊急。
type Urgent interface {
	NotifyUrgent(msg string)
}

func (s SMS) NotifyUrgent(msg string) {
	fmt.Printf("簡訊 %s：【緊急】%s\n", s.Phone, msg)
}

// alert 先問 n 有沒有緊急通知的能力，有就用，沒有就退回一般通知。
func alert(n Notifier, msg string) {
	if u, ok := n.(Urgent); ok {
		u.NotifyUrgent(msg)
		return
	}
	n.Notify(msg)
}

func optionalDemo() {
	alert(SMS{Phone: "0912-345-678"}, "訂單被取消了")
	alert(Email{Addr: "amei@example.com"}, "訂單被取消了")
}

// label 用 type switch 依實際型別回傳不同的說明。
func label(v any) string {
	switch x := v.(type) {
	case nil:
		return "什麼都沒有"
	case int, int64:
		return fmt.Sprintf("整數 %v", x)
	case string:
		return "字串 " + x
	case error:
		return "錯誤：" + x.Error()
	case Notifier:
		return fmt.Sprintf("通知方式 %T", x)
	default:
		return fmt.Sprintf("不認得的 %T", x)
	}
}

func typeSwitchDemo() {
	values := []any{nil, 42, int64(7), "排骨飯", errors.New("賣完了"), SMS{Phone: "0912"}, 3.5}
	for _, v := range values {
		fmt.Println(label(v))
	}
}

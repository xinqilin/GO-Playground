package main

import "fmt"

// Notifier 是「能通知客人」的東西：只要有 Notify 方法就算。
type Notifier interface {
	Notify(msg string)
}

// SMS 用簡訊通知。
type SMS struct {
	Phone string
}

func (s SMS) Notify(msg string) {
	fmt.Printf("簡訊 %s：%s\n", s.Phone, msg)
}

// Email 用電子郵件通知。
type Email struct {
	Addr string
}

func (e Email) Notify(msg string) {
	fmt.Printf("Email %s：%s\n", e.Addr, msg)
}

// 編譯期檢查：SMS 或 Email 少了 Notify，這兩行就編譯不過。
var (
	_ Notifier = SMS{}
	_ Notifier = Email{}
)

// callNumber 只認得 Notifier，不管背後是簡訊還是 Email。
func callNumber(n Notifier, no int) {
	n.Notify(fmt.Sprintf("%d 號的餐點好了", no))
}

// implicitDemo：SMS 和 Email 都沒有寫「我實作了 Notifier」，有 Notify 方法就自動符合。
func implicitDemo() {
	callNumber(SMS{Phone: "0912-345-678"}, 17)
	callNumber(Email{Addr: "amei@example.com"}, 18)

	list := []Notifier{SMS{Phone: "0987-654-321"}, Email{Addr: "kuo@example.com"}}
	for _, n := range list {
		n.Notify("明天公休")
	}
}

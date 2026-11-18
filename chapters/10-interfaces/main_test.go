package main

import (
	"slices"
	"testing"
)

func Example_implicitDemo() {
	implicitDemo()
	// Output:
	// 簡訊 0912-345-678：17 號的餐點好了
	// Email amei@example.com：18 號的餐點好了
	// 簡訊 0987-654-321：明天公休
	// Email kuo@example.com：明天公休
}

func Example_interfaceValue() {
	interfaceValue()
	// Output:
	// true
	// <nil>
	// main.SMS {0912-345-678}
	// main.Email {amei@example.com}
	// false
}

func Example_pointerInInterface() {
	pointerInInterface()
	// Output:
	// 呼叫器 5 號響了（第 1 次）：餐點好了
	// 呼叫器 5 號響了（第 2 次）：請到櫃台取餐
	// 2
}

func Example_compareInterfaces() {
	compareInterfaces()
	// Output:
	// true false
	// panic: runtime error: comparing uncomparable type []string
}

func Example_assertDemo() {
	assertDemo()
	// Output:
	// 0912-345-678
	// false {}
}

func Example_assertPanic() {
	assertPanic()
	// Output:
	// panic: interface conversion: main.Notifier is main.SMS, not main.Email
}

func Example_optionalDemo() {
	optionalDemo()
	// Output:
	// 簡訊 0912-345-678：【緊急】訂單被取消了
	// Email amei@example.com：訂單被取消了
}

func Example_typeSwitchDemo() {
	typeSwitchDemo()
	// Output:
	// 什麼都沒有
	// 整數 42
	// 整數 7
	// 字串 排骨飯
	// 錯誤：賣完了
	// 通知方式 main.SMS
	// 不認得的 float64
}

func Example_nilTrapDemo() {
	nilTrapDemo()
	// Output:
	// false
	// *main.FormError
	// true
	// 姓名沒有填
}

func Example_writerDemo() {
	writerDemo()
	// Output:
	// 13 <nil>
	// 6 紙用完了
	// "排骨飯 x2\n滷蛋"
}

func Example_sortDemo() {
	sortDemo()
	// Output:
	// [{1 true} {4 true} {2 false} {3 false}]
	// [{1 true} {2 false} {3 false} {4 true}]
}

func Example_embedInterfaceDemo() {
	embedInterfaceDemo()
	// Output:
	// 排骨飯
	// 第一次來
}

func Example_remindDemo() {
	remindDemo()
	// Output:
	// 簡訊 0912-345-678：17 號請取餐
	// 簡訊 0912-345-678：21 號請取餐
}

// recorder 是測試用的假 Notifier：不真的發通知，只記下收到的訊息。
type recorder struct {
	msgs []string
}

func (r *recorder) Notify(msg string) {
	r.msgs = append(r.msgs, msg)
}

func TestRemindAll(t *testing.T) {
	var r recorder
	remindAll(&r, []int{17, 21})
	want := []string{"17 號請取餐", "21 號請取餐"}
	if !slices.Equal(r.msgs, want) {
		t.Errorf("收到 %q, want %q", r.msgs, want)
	}
}

func TestCheckName(t *testing.T) {
	if err := checkName("王小明"); err != nil {
		t.Errorf("checkName(王小明) = %v, want nil", err)
	}
	if err := checkName(""); err == nil {
		t.Error("checkName(\"\") = nil, want error")
	}
	if err := checkNameBuggy("王小明"); err == nil {
		t.Error("checkNameBuggy 應該示範出「不是 nil」的錯誤")
	}
}

func TestLastOrder(t *testing.T) {
	s := &memStore{data: map[string]string{"小明": "排骨飯"}}
	tests := []struct {
		member, want string
	}{
		{"小明", "排骨飯"},
		{"阿美", "第一次來"},
	}
	for _, tt := range tests {
		if got := lastOrder(s, tt.member); got != tt.want {
			t.Errorf("lastOrder(%q) = %q, want %q", tt.member, got, tt.want)
		}
	}
}

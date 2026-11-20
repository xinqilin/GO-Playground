package kiosk

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestAsFloat(t *testing.T) {
	tests := []struct {
		v      any
		want   float64
		wantOK bool
	}{
		{85, 85, true},
		{int64(-3), -3, true},
		{float32(0.5), 0.5, true},
		{1.25, 1.25, true},
		{"85.5", 85.5, true},
		{"abc", 0, false},
		{" 12", 0, false},
		{true, 0, false},
		{nil, 0, false},
		{[]int{1}, 0, false},
	}
	for _, tt := range tests {
		got, ok := AsFloat(tt.v)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("AsFloat(%#v) = %v, %v, want %v, %v", tt.v, got, ok, tt.want, tt.wantOK)
		}
	}
}

// flat 是測試用的寄送方式：不論重量都是同一個價錢。
type flat int

func (f flat) Fee(int) int { return int(f) }

func TestShipping(t *testing.T) {
	fees := []struct {
		grams, want int
	}{
		{0, 100}, {5000, 100}, {5001, 130}, {6000, 130}, {6001, 160}, {8000, 190},
	}
	for _, tt := range fees {
		if got := (HomeDelivery{}).Fee(tt.grams); got != tt.want {
			t.Errorf("HomeDelivery.Fee(%d) = %d, want %d", tt.grams, got, tt.want)
		}
	}

	both := []Shipper{StorePickup{}, HomeDelivery{}}
	tests := []struct {
		name    string
		options []Shipper
		grams   int
		want    int
	}{
		{"超商比較便宜", both, 3000, 0},
		{"太重只能宅配", both, 8000, 1},
		{"沒有選項", nil, 1000, -1},
		{"全部寄不了", []Shipper{StorePickup{}}, 6000, -1},
		{"一樣便宜取前面的", []Shipper{flat(60), StorePickup{}}, 1000, 0},
		{"假的比較便宜", []Shipper{HomeDelivery{}, flat(80)}, 1000, 1},
	}
	for _, tt := range tests {
		if got := Cheapest(tt.options, tt.grams); got != tt.want {
			t.Errorf("%s: Cheapest(%d 公克) = %d, want %d", tt.name, tt.grams, got, tt.want)
		}
	}
}

// failWriter 是每次寫入都失敗的 io.Writer。
type failWriter struct{}

func (failWriter) Write([]byte) (int, error) { return 0, errors.New("磁碟滿了") }

func TestMasker(t *testing.T) {
	var b strings.Builder
	fmt.Fprintf(Masker{W: &b}, "電話 %s 金額 %d", "0912-345-678", 85)
	if got, want := b.String(), "電話 ****-***-*** 金額 **"; got != want {
		t.Errorf("寫出 %q, want %q", got, want)
	}

	b.Reset()
	p := []byte("A1")
	n, err := Masker{W: &b}.Write(p)
	if n != 2 || err != nil || b.String() != "A*" {
		t.Errorf("Write(%q) = %d, %v，寫出 %q, want 2, nil，寫出 \"A*\"", "A1", n, err, b.String())
	}
	if string(p) != "A1" {
		t.Errorf("呼叫者的 p 被改成 %q 了", p)
	}

	if _, err := (Masker{W: failWriter{}}).Write([]byte("1")); err == nil {
		t.Error("底下的 Writer 失敗時，Write 應該回傳錯誤")
	}
}

// voucher 是測試用的品項：本身 0 元，卻宣稱可以折 30 元。
type voucher struct{}

func (voucher) Price() int    { return 0 }
func (voucher) Discount() int { return 30 }

func TestTotal(t *testing.T) {
	tests := []struct {
		name  string
		items []Item
		want  int
	}{
		{"空的", nil, 0},
		{"沒有折扣", []Item{Bento{"排骨飯", 85}, Bento{"雞腿飯", 95}}, 180},
		{"套餐有折扣", []Item{Combo{"A 餐", 120, 20}, Bento{"排骨飯", 85}}, 185},
		{"折扣比價錢多", []Item{Combo{"試吃", 50, 80}}, 0},
		{"其他型別的折扣", []Item{voucher{}, Bento{"排骨飯", 85}}, 85},
	}
	for _, tt := range tests {
		if got := Total(tt.items); got != tt.want {
			t.Errorf("%s: Total = %d, want %d", tt.name, got, tt.want)
		}
	}
}

// fakeSender 是測試用的 Sender：err 不是 nil 就失敗，否則記下訊息。
type fakeSender struct {
	err error
	got []string
}

func (f *fakeSender) Send(msg string) error {
	if f.err != nil {
		return f.err
	}
	f.got = append(f.got, msg)
	return nil
}

func TestFallback(t *testing.T) {
	var _ Sender = Fallback{}

	if err := (Fallback{}).Send("hi"); !errors.Is(err, ErrNoSender) {
		t.Errorf("空的 Fallback.Send = %v, want ErrNoSender", err)
	}

	errA, errB := errors.New("簡訊失敗"), errors.New("Email 失敗")
	bad, ok, spare := &fakeSender{err: errA}, &fakeSender{}, &fakeSender{}
	if err := (Fallback{bad, ok, spare}).Send("17 號"); err != nil {
		t.Errorf("第二個成功時 Send = %v, want nil", err)
	}
	if len(ok.got) != 1 || len(spare.got) != 0 {
		t.Errorf("ok 收到 %q、spare 收到 %q；成功之後不應該再嘗試後面的", ok.got, spare.got)
	}

	err := (Fallback{&fakeSender{err: errA}, &fakeSender{err: errB}}).Send("18 號")
	if err != errB {
		t.Errorf("全部失敗時 Send = %v, want 最後一個錯誤 %v", err, errB)
	}

	inner := &fakeSender{}
	nested := Fallback{Fallback{&fakeSender{err: errA}}, inner}
	if err := nested.Send("19 號"); err != nil || len(inner.got) != 1 {
		t.Errorf("Fallback 裡放 Fallback：Send = %v，inner 收到 %q", err, inner.got)
	}
}

func TestCheckOrder(t *testing.T) {
	if err := CheckOrder(Order{Items: 2, Paid: true}); err != nil {
		t.Errorf("正常的訂單 CheckOrder = %v (%T), want nil", err, err)
	}
	tests := []struct {
		o    Order
		want string
	}{
		{Order{Items: 0, Paid: true}, "沒有點東西"},
		{Order{Items: 1, Paid: false}, "還沒付款"},
	}
	for _, tt := range tests {
		err := CheckOrder(tt.o)
		oe, ok := err.(*OrderError)
		if !ok || oe == nil || oe.Reason != tt.want {
			t.Errorf("CheckOrder(%+v) = %v, want *OrderError 的 Reason 是 %q", tt.o, err, tt.want)
		}
	}
}

// tag 是測試用的型別，有 String 方法。
type tag string

func (t tag) String() string { return "#" + string(t) }

func TestKind(t *testing.T) {
	tests := []struct {
		v    any
		want string
	}{
		{Price(85), "價格"},
		{tag("便當"), "可以印的"},
		{errors.New("賣完了"), "錯誤"},
		{42, "其他"},
		{nil, "其他"},
	}
	for _, tt := range tests {
		if got := Kind(tt.v); got != tt.want {
			t.Errorf("Kind(%#v) = %q, want %q", tt.v, got, tt.want)
		}
	}
}

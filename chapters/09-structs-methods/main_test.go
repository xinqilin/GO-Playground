package main

import "testing"

func Example_definedTypes() {
	definedTypes()
	// Output:
	// 1.5
	// 900 公克
	// 350 公克
	// 300
}

func Example_aliasDemo() {
	aliasDemo()
	// Output:
	// int main.Points 90
}

func Example_structLiterals() {
	structLiterals()
	// Output:
	// { 0 false}
	// {Name:排骨飯 Price:85 Spicy:false}
	// main.Bento{Name:"雞腿飯", Price:95, Spicy:true}
	// 排骨飯 90 false
}

func Example_structPointer() {
	structPointer()
	// Output:
	// 85
	// &{Name:素食餐 Price:85 Spicy:false}
}

func Example_anonymousStruct() {
	anonymousStruct()
	// Output:
	// 一 3200
	// 二 4100
}

func Example_copyStruct() {
	copyStruct()
	// Output:
	// 85 100
}

func Example_shallowCopy() {
	shallowCopy()
	// Output:
	// A 餐 [豆干 青菜]
	// B 餐 [豆干 青菜]
}

func Example_compareStructs() {
	compareStructs()
	// Output:
	// true false
	// 2 2
}

func Example_receiverDemo() {
	receiverDemo()
	// Output:
	// 0
	// 2
}

func Example_autoAddress() {
	autoAddress()
	// Output:
	// true 11
}

func Example_methodSetDemo() {
	methodSetDemo()
	// Output:
	// {7}
	// 第 7 號
}

func Example_zeroValueDemo() {
	zeroValueDemo()
	// Output:
	// 1 2 2
	// 1 true 1
	// 0 false
}

func Example_constructorDemo() {
	constructorDemo()
	// Output:
	// 王先生 4 <nil>
	// 人數 20 不在 1 到 12 之間
	// 訂位要留名字
}

func Example_tallyDemo() {
	tallyDemo()
	// Output:
	// 0
	// 2 1
}

func Example_embedDemo() {
	embedDemo()
	// Output:
	// 阿姨便當
	// 阿姨便當 02-2345-6789
	// 阿姨便當
}

func Example_promotedKeys() {
	promotedKeys()
	// Output:
	// 巷口便當 02-8765-4321
}

func Example_shadowDemo() {
	shadowDemo()
	// Output:
	// 手搖飲 02-1111-2222 LINE:@drink
	// 手搖飲 02-1111-2222
}

func Example_notInheritance() {
	notInheritance()
	// Output:
	// 貴賓
	// 歡迎光臨，客人
}

func Example_tagsDemo() {
	tagsDemo()
	// Output:
	// {"name":"小明","points":120} <nil>
	// phone,omitempty
}

func TestNewReservation(t *testing.T) {
	tests := []struct {
		name    string
		people  int
		wantErr bool
	}{
		{"王先生", 4, false},
		{"王先生", 1, false},
		{"王先生", 12, false},
		{"王先生", 0, true},
		{"王先生", 13, true},
		{"", 2, true},
	}
	for _, tt := range tests {
		r, err := NewReservation(tt.name, tt.people)
		if (err != nil) != tt.wantErr {
			t.Errorf("NewReservation(%q, %d) error = %v, wantErr %v", tt.name, tt.people, err, tt.wantErr)
			continue
		}
		if err == nil && (r.Name != tt.name || r.People != tt.people) {
			t.Errorf("NewReservation(%q, %d) = %+v", tt.name, tt.people, *r)
		}
		if err != nil && r != nil {
			t.Errorf("NewReservation(%q, %d) 失敗時應該回傳 nil，結果是 %+v", tt.name, tt.people, *r)
		}
	}
}

func TestQueueEmpty(t *testing.T) {
	var q Queue
	if n, ok := q.Call(); ok || n != 0 {
		t.Errorf("空的 Queue 叫號 = %d, %v, want 0, false", n, ok)
	}
	q.Take()
	q.Call()
	if n, ok := q.Call(); ok || n != 0 {
		t.Errorf("叫完之後再叫號 = %d, %v, want 0, false", n, ok)
	}
	if q.Waiting() != 0 {
		t.Errorf("Waiting() = %d, want 0", q.Waiting())
	}
}

func TestJinGrams(t *testing.T) {
	tests := []struct {
		j    Jin
		want Gram
	}{
		{0, 0},
		{0.5, 300},
		{1, 600},
		{2.5, 1500},
	}
	for _, tt := range tests {
		if got := tt.j.Grams(); got != tt.want {
			t.Errorf("Jin(%v).Grams() = %v, want %v", float64(tt.j), got, tt.want)
		}
	}
}

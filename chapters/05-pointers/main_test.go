package main

import "testing"

func Example_addressAndDeref() {
	addressAndDeref()
	// Output:
	// stock = 5
	// *p = 5
	// p == &stock: true
	// *int
}

func Example_modifyThroughPointer() {
	modifyThroughPointer()
	// Output:
	// stock = 13
}

func Example_aliasing() {
	aliasing()
	// Output:
	// 8 true
	// true false
}

func Example_passValueVsPointer() {
	passValueVsPointer()
	// Output:
	// a = 1 b = 2
}

func Example_structCopy() {
	structCopy()
	// Output:
	// 85 999
	// 85
	// 90
}

func Example_shallowCopy() {
	shallowCopy()
	// Output:
	// 1 加辣
	// 2 加辣
}

func Example_structPointer() {
	structPointer()
	// Output:
	// 95 95
}

func Example_newVsLiteral() {
	newVsLiteral()
	// Output:
	// true false
}

func Example_newWithValue() {
	newWithValue()
	// Output:
	// 15
	// 少飯
}

func Example_nilPointer() {
	nilPointer()
	// Output:
	// true
	// panic: runtime error: invalid memory address or nil pointer dereference
}

func Example_optionalField() {
	optionalField()
	// Output:
	// 100 元（沒有設定折扣）
	// 100 元（折扣 0%）
	// 80 元（折扣 20%）
}

func Example_returnLocalPointer() {
	returnLocalPointer()
	// Output:
	// 2 1 false
}

func TestPriceAfter(t *testing.T) {
	ten, hundred := 10, 100
	tests := []struct {
		price    int
		discount *int
		want     string
	}{
		{85, nil, "85 元（沒有設定折扣）"},
		{85, &ten, "76 元（折扣 10%）"},
		{85, &hundred, "0 元（折扣 100%）"},
	}
	for _, tt := range tests {
		if got := priceAfter(tt.price, tt.discount); got != tt.want {
			t.Errorf("priceAfter(%d, …) = %q, want %q", tt.price, got, tt.want)
		}
	}
}

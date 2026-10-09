package main

import "testing"

func Example() {
	main()
	// Output:
	// == 取地址與解參考 ==
	// x = 5
	// *p = 5
	// p == &x: true
	// == 透過指標修改 ==
	// x = 10
	// == 傳值 vs 傳指標 ==
	// a = 1 b = 2
	// == nil 與 new ==
	// p == nil: true
	// *q = 42
	// *r = 7
	// == 回傳區域變數的指標 ==
	// c.n = 2
}

func TestAddOne(t *testing.T) {
	a, b := 1, 1
	addOneValue(a)
	addOnePointer(&b)
	if a != 1 {
		t.Errorf("addOneValue 不該改到呼叫端：a = %d, want 1", a)
	}
	if b != 2 {
		t.Errorf("addOnePointer 應該改到呼叫端：b = %d, want 2", b)
	}
}

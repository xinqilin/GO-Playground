package main

import "testing"

func Example() {
	main()
	// Output:
	// == 字串是 byte 序列 ==
	// 8
	// 4
	// 71 232
	// 語
	// == range 一次拿一個 rune ==
	// 0 G 71
	// 1 o 111
	// 2 語 35486
	// 5 言 35328
	// == strings 套件 ==
	// HELLO, GO true true
	// [a b c] [a b]
	// lang go true
	// == 組字串用 Builder ==
	// Go-Go-Go
	// == 數字與字串互轉 ==
	// 124 <nil>
	// strconv.Atoi: parsing "12a": invalid syntax
	// 7 <nil>
	// 42! "你好\n"
	// C:\path\no\escape
}

func TestRepeatJoin(t *testing.T) {
	tests := []struct {
		word string
		n    int
		want string
	}{
		{"Go", 0, ""},
		{"Go", 1, "Go"},
		{"中", 3, "中-中-中"},
	}
	for _, tt := range tests {
		if got := repeatJoin(tt.word, tt.n); got != tt.want {
			t.Errorf("repeatJoin(%q, %d) = %q, want %q", tt.word, tt.n, got, tt.want)
		}
	}
}

package main

import (
	"maps"
	"testing"
)

func Example() {
	main()
	// Output:
	// == 基本操作 ==
	// map[小明:19 小美:22] 2
	// == 找不到的 key ==
	// 0
	// true false
	// == 走訪順序不固定 ==
	// a 1
	// b 2
	// c 3
	// == 數字數 ==
	// map[and:1 fast:1 fun:1 go:2 is:2]
	// == nil map ==
	// 0 0 true
	// 寫入 nil map： assignment to entry in nil map
}

func TestWordCount(t *testing.T) {
	tests := []struct {
		in   string
		want map[string]int
	}{
		{"", map[string]int{}},
		{"a a b", map[string]int{"a": 2, "b": 1}},
		{"  多個   空白  ", map[string]int{"多個": 1, "空白": 1}},
	}
	for _, tt := range tests {
		if got := wordCount(tt.in); !maps.Equal(got, tt.want) {
			t.Errorf("wordCount(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

package main

import "testing"

func Example() {
	main()
	// Output:
	// == 泛型函式 ==
	// 6
	// 42
	// [1 2 3]
	// zoo true
	// 0 false
	// == 泛型型別 ==
	// 2 [1]
	// == generic method ==
	// [#1 #3]
}

func TestStack(t *testing.T) {
	var s Stack[string]
	if _, ok := s.Pop(); ok {
		t.Fatal("空的 Stack Pop 應該回傳 false")
	}
	s.Push("a")
	s.Push("b")
	for _, want := range []string{"b", "a"} {
		got, ok := s.Pop()
		if !ok || got != want {
			t.Fatalf("Pop() = %q, %v; want %q, true", got, ok, want)
		}
	}
}

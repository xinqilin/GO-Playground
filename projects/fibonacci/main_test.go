package main

import (
	"os"
	"strings"
	"testing"
)

// 舊版在 n = 0 會無限遞迴、n = 2 會算錯，這裡把它們都列進來。
func TestImplementations(t *testing.T) {
	tests := []struct{ n, want int }{
		{-3, 0},
		{0, 0},
		{1, 1},
		{2, 1},
		{10, 55},
		{20, 6765},
	}
	memo := Memoized()
	for _, tt := range tests {
		if got := Recursive(tt.n); got != tt.want {
			t.Errorf("Recursive(%d) = %d, want %d", tt.n, got, tt.want)
		}
		if got := Iterative(tt.n); got != tt.want {
			t.Errorf("Iterative(%d) = %d, want %d", tt.n, got, tt.want)
		}
		if got := memo(tt.n); got != tt.want {
			t.Errorf("Memoized(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

// 大的 n 只比較兩個快的版本，並確認最大項沒有溢位。
func TestLargeN(t *testing.T) {
	memo := Memoized()
	for n := range maxN + 1 {
		if a, b := Iterative(n), memo(n); a != b {
			t.Fatalf("n=%d: Iterative %d != Memoized %d", n, a, b)
		}
	}
	if got := Iterative(maxN); got != 7540113804746346429 {
		t.Errorf("F(%d) = %d", maxN, got)
	}
}

func TestRunRejectsBadN(t *testing.T) {
	for _, arg := range []string{"-1", "93", "abc"} {
		var out strings.Builder
		if err := run([]string{"-n", arg}, &out); err == nil {
			t.Errorf("-n %s: 應該回傳錯誤", arg)
		}
	}
}

func Example() {
	if err := run([]string{"-n", "10"}, os.Stdout); err != nil {
		panic(err)
	}
	// Output:
	// F(10)
	// 迭代： 55
	// 記憶化： 55
	// 遞迴： 55
}

func BenchmarkRecursive25(b *testing.B) {
	for b.Loop() {
		Recursive(25)
	}
}

func BenchmarkIterative25(b *testing.B) {
	for b.Loop() {
		Iterative(25)
	}
}

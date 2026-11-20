package main

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"strings"
	"testing"
)

// 舊版在 n = 0 會無限遞迴、n = 2 會算錯，這裡把它們都列進來。
func TestSmallN(t *testing.T) {
	tests := []struct{ n, want int }{
		{-3, 0},
		{0, 0},
		{1, 1},
		{2, 1},
		{10, 55},
		{20, 6765},
	}
	for _, tt := range tests {
		memo := Memoized()
		checked, err := Checked(tt.n)
		got := map[string]int{
			"Recursive": Recursive(tt.n),
			"Memoized":  memo(tt.n),
			"Iterative": Iterative(tt.n),
			"Checked":   checked,
			"Big":       int(Big(tt.n).Int64()),
		}
		if err != nil {
			t.Errorf("Checked(%d) 的錯誤：%v", tt.n, err)
		}
		for name, v := range got {
			if v != tt.want {
				t.Errorf("%s(%d) = %d，想要 %d", name, tt.n, v, tt.want)
			}
		}
	}
}

// TestAgreeWithBig 把每一種 int 的寫法都和 math/big 的結果比對，一直到 F(maxN)。
func TestAgreeWithBig(t *testing.T) {
	memo := Memoized()
	var seq []int
	for v := range Seq() {
		seq = append(seq, v)
	}
	if len(seq) != maxN+1 {
		t.Fatalf("Seq 產生了 %d 項，想要 %d 項", len(seq), maxN+1)
	}
	for n := range maxN + 1 {
		want := Big(n)
		checked, err := Checked(n)
		if err != nil {
			t.Fatalf("Checked(%d) 的錯誤：%v", n, err)
		}
		for name, v := range map[string]int{"Memoized": memo(n), "Iterative": Iterative(n), "Checked": checked, "Seq": seq[n]} {
			if big.NewInt(int64(v)).Cmp(want) != 0 {
				t.Fatalf("%s(%d) = %d，math/big 是 %s", name, n, v, want)
			}
		}
	}
	if got := Iterative(maxN); got != 7540113804746346429 {
		t.Errorf("F(%d) = %d", maxN, got)
	}
}

func TestOverflow(t *testing.T) {
	if _, err := Checked(maxN + 1); !errors.Is(err, ErrOverflow) {
		t.Errorf("Checked(%d) 的錯誤 = %v，想要 ErrOverflow", maxN+1, err)
	}
	if Big(maxN + 1).IsInt64() {
		t.Errorf("F(%d) 應該超過 int64", maxN+1)
	}
	// 不檢查的版本會默默算錯：結果變成負數。
	if got := Iterative(maxN + 1); got >= 0 {
		t.Errorf("Iterative(%d) = %d，想要溢位後的負數", maxN+1, got)
	}
}

// 呼叫次數剛好是 2·F(n+1) - 1。
func TestCountCalls(t *testing.T) {
	for n := range 26 {
		if got, want := countCalls(n), 2*Iterative(n+1)-1; got != want {
			t.Errorf("countCalls(%d) = %d，想要 %d", n, got, want)
		}
	}
}

func TestRunErrors(t *testing.T) {
	for _, args := range [][]string{{"-n", "-1"}, {"-n", "abc"}, {"-list", "-2"}, {"-list", "94"}} {
		var out strings.Builder
		if err := run(args, &out); err == nil {
			t.Errorf("%q：想要錯誤", args)
		}
	}
}

// TestRunHelp：-h 只印用法，不算錯誤（main 不會再印「錯誤：」，exit code 是 0）。
func TestRunHelp(t *testing.T) {
	var out strings.Builder
	if err := run([]string{"-h"}, &out); err != nil || !strings.Contains(out.String(), "-list") {
		t.Errorf("run(-h) = %v，輸出\n%s想要 nil 和用法", err, out.String())
	}
}

func Example_run() {
	if err := run([]string{"-n", "30"}, os.Stdout); err != nil {
		fmt.Println(err)
	}
	// Output:
	// F(30) = 832040
	// 遞迴：832040（呼叫了 2692537 次）
	// 記憶化：832040
	// 迴圈：832040
	// math/big：832040
}

func Example_runBig() {
	if err := run([]string{"-n", "100"}, os.Stdout); err != nil {
		fmt.Println(err)
	}
	// Output:
	// F(100) = 354224848179261915075
	// 超過 int 的範圍（最大到 F(92)），只用 math/big 計算
}

func Example_runList() {
	if err := run([]string{"-list", "15"}, os.Stdout); err != nil {
		fmt.Println(err)
	}
	// Output:
	// 0 1 1 2 3 5 8 13 21 34 55 89 144 233 377
}

func Example_countCalls() {
	for _, n := range []int{5, 10, 20, 30} {
		fmt.Printf("Recursive(%d) 呼叫了 %d 次\n", n, countCalls(n))
	}
	// Output:
	// Recursive(5) 呼叫了 15 次
	// Recursive(10) 呼叫了 177 次
	// Recursive(20) 呼叫了 21891 次
	// Recursive(30) 呼叫了 2692537 次
}

func Example_overflow() {
	fmt.Printf("math.MaxInt   = %d\n", math.MaxInt)
	fmt.Printf("F(92)         = %d\n", Iterative(92))
	fmt.Printf("F(93)         = %s\n", Big(93))
	fmt.Printf("Iterative(93) = %d\n", Iterative(93))
	_, err := Checked(93)
	fmt.Println(err)
	// Output:
	// math.MaxInt   = 9223372036854775807
	// F(92)         = 7540113804746346429
	// F(93)         = 12200160415121876738
	// Iterative(93) = -6246583658587674878
	// F(93)：超過 int 的範圍
}

func Example_big() {
	f := Big(1000)
	s := f.String()
	fmt.Printf("F(1000) 有 %d 位數，開頭是 %s…，結尾是 …%s\n", len(s), s[:10], s[len(s)-10:])
	// Output:
	// F(1000) 有 209 位數，開頭是 4346655768…，結尾是 …6849228875
}

// BenchmarkFib 在同一個 n 比較五種寫法。
func BenchmarkFib(b *testing.B) {
	const n = 30
	b.Run("recursive", func(b *testing.B) {
		for b.Loop() {
			Recursive(n)
		}
	})
	b.Run("memoized", func(b *testing.B) {
		for b.Loop() {
			Memoized()(n)
		}
	})
	b.Run("iterative", func(b *testing.B) {
		for b.Loop() {
			Iterative(n)
		}
	})
	b.Run("seq", func(b *testing.B) {
		for b.Loop() {
			i := 0
			for range Seq() {
				if i == n {
					break
				}
				i++
			}
		}
	})
	b.Run("big", func(b *testing.B) {
		for b.Loop() {
			Big(n)
		}
	})
}

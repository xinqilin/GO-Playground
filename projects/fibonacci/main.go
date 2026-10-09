// 專題 1：費氏數列。
//
// 改寫自這個 repo 最早的 Exercise/fibonacci.go。舊版有三個問題：
//  1. 隨機抽到 n = 0 時，fibonacci(0) 會呼叫 fibonacci(-1)、fibonacci(-2)……永遠停不下來，最後 stack overflow。
//  2. fibonacci(2) 回傳 2，整個數列錯開一格（正確是 0, 1, 1, 2, 3, 5…）。
//  3. rand.Seed 自 Go 1.20 起已經 deprecated，現在不用自己設種子。
//
// 執行：go run ./projects/fibonacci -n 30
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// maxN 是 int64 放得下的最大項：F(92) = 7540113804746346429，F(93) 就會溢位。
const maxN = 92

// Recursive 直接照定義遞迴：F(n) = F(n-1) + F(n-2)。很好懂，但同樣的值會被重複算很多次。
// n ≤ 0 時回傳 0。
func Recursive(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	return Recursive(n-1) + Recursive(n-2)
}

// Iterative 只記住最後兩項，一路往前推，時間 O(n)、空間 O(1)。
func Iterative(n int) int {
	a, b := 0, 1
	for range max(n, 0) {
		a, b = b, a+b
	}
	return a
}

// Memoized 回傳一個會記住算過結果的函式：每個 n 只算一次。
func Memoized() func(int) int {
	memo := map[int]int{}
	var fib func(int) int
	fib = func(n int) int {
		if n <= 0 {
			return 0
		}
		if n == 1 {
			return 1
		}
		if v, ok := memo[n]; ok {
			return v
		}
		v := fib(n-1) + fib(n-2)
		memo[n] = v
		return v
	}
	return fib
}

// run 解析參數並輸出結果；和 main 分開，測試時就能傳入自己的參數和輸出。
func run(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("fibonacci", flag.ContinueOnError)
	fs.SetOutput(out)
	n := fs.Int("n", 10, fmt.Sprintf("要算第幾項（0 到 %d）", maxN))
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *n < 0 || *n > maxN {
		return fmt.Errorf("n 必須介於 0 到 %d，收到 %d", maxN, *n)
	}

	fmt.Fprintf(out, "F(%d)\n", *n)
	fmt.Fprintln(out, "迭代：", Iterative(*n))
	fmt.Fprintln(out, "記憶化：", Memoized()(*n))
	if *n <= 35 { // 直接遞迴超過 35 項就慢到不值得等
		fmt.Fprintln(out, "遞迴：", Recursive(*n))
	} else {
		fmt.Fprintln(out, "遞迴：略過（n 太大，會跑很久）")
	}
	return nil
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "錯誤：", err)
		os.Exit(2)
	}
}

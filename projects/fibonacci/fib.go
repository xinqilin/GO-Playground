package main

import (
	"errors"
	"fmt"
	"iter"
	"math"
	"math/big"
)

// maxN 是 int（64 位元）放得下的最大項：F(92) = 7540113804746346429，F(93) 就超過上限。
const maxN = 92

// ErrOverflow 表示結果超過 int 能表示的範圍。
var ErrOverflow = errors.New("超過 int 的範圍")

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

// countCalls 回傳 Recursive(n) 總共會被呼叫幾次，包含最外面那一次。
func countCalls(n int) int {
	calls := 0
	var f func(int) int
	f = func(n int) int {
		calls++
		if n <= 0 {
			return 0
		}
		if n == 1 {
			return 1
		}
		return f(n-1) + f(n-2)
	}
	f(n)
	return calls
}

// Memoized 回傳一個會記住算過結果的函式：每個 n 只算一次。n 不能超過 maxN。
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

// Iterative 只記住最後兩項，一路往前推：時間 O(n)、空間 O(1)。n 不能超過 maxN。
func Iterative(n int) int {
	a, b := 0, 1
	for range max(n, 0) {
		a, b = b, a+b
	}
	return a
}

// Checked 和 Iterative 一樣用迴圈，但在加法溢位之前停下來，回傳 ErrOverflow。
func Checked(n int) (int, error) {
	if n <= 0 {
		return 0, nil
	}
	a, b := 0, 1 // F(0)、F(1)
	for range n - 1 {
		if a > math.MaxInt-b {
			return 0, fmt.Errorf("F(%d)：%w", n, ErrOverflow)
		}
		a, b = b, a+b
	}
	return b, nil
}

// Seq 依序產生 F(0)、F(1)、F(2)…，到 F(maxN) 為止，下一項會溢位時就結束。
func Seq() iter.Seq[int] {
	return func(yield func(int) bool) {
		a, b := 0, 1
		nextFits := true // b 是正確的值嗎（之前的加法沒有溢位）
		for {
			if !yield(a) || !nextFits {
				return
			}
			nextFits = a <= math.MaxInt-b
			a, b = b, a+b
		}
	}
}

// Big 用 math/big 計算 F(n)，沒有上限。a、b 兩個數字重複使用，不會每一步都配置新的。
func Big(n int) *big.Int {
	a, b := big.NewInt(0), big.NewInt(1)
	for range max(n, 0) {
		a.Add(a, b)
		a, b = b, a
	}
	return a
}

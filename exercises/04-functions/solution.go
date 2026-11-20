//go:build solution

package functions

import "fmt"

func MinMax(nums ...int) (lo, hi int, ok bool) {
	if len(nums) == 0 {
		return 0, 0, false
	}
	lo, hi = nums[0], nums[0]
	for _, n := range nums[1:] {
		lo = min(lo, n)
		hi = max(hi, n)
	}
	return lo, hi, true
}

func CountIf(nums []int, keep func(int) bool) int {
	count := 0
	for _, n := range nums {
		if keep(n) {
			count++
		}
	}
	return count
}

// NewDispenser 的兩個 closure 抓住同一個 n，所以 reset 改了它，next 也看得到。
func NewDispenser(prefix string) (next func() string, reset func()) {
	n := 0
	next = func() string {
		n++
		return fmt.Sprintf("%s%03d", prefix, n)
	}
	reset = func() { n = 0 }
	return next, reset
}

func PercentOff(pct int) Rule {
	return func(price int) int {
		return price * (100 - pct) / 100
	}
}

func MinusOver(threshold, off int) Rule {
	return func(price int) int {
		if price >= threshold {
			return price - off
		}
		return price
	}
}

func ApplyRules(price int, rules ...Rule) int {
	for _, rule := range rules {
		price = rule(price)
	}
	return max(price, 0)
}

// ProcessAll 把每一圈包成一個函式，defer 在那個函式結束時就執行。
func ProcessAll(t *Tracker, names []string) {
	for _, name := range names {
		func() {
			t.Open(name)
			defer t.Close(name)
			t.Use(name)
		}()
	}
}

// SafeDiv 的回傳值要有名字，defer 裡的 closure 才能在 panic 之後改到真正要回傳的 err。
func SafeDiv(a, b int) (q int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("除法失敗：%v", r)
		}
	}()
	return a / b, nil
}

// PayWays 把問題分成兩種情況：至少用一個 coins[0]，或者完全不用 coins[0]。
func PayWays(amount int, coins []int) int {
	if amount == 0 {
		return 1
	}
	if amount < 0 || len(coins) == 0 {
		return 0
	}
	return PayWays(amount-coins[0], coins) + PayWays(amount, coins[1:])
}

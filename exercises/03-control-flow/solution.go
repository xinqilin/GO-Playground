//go:build solution

package controlflow

import (
	"fmt"
	"strings"
	"unicode"
)

func TicketPrice(km int) int {
	switch {
	case km <= 0:
		return 0
	case km <= 5:
		return 20
	case km <= 10:
		return 25
	case km <= 20:
		return 35
	case km <= 40:
		return 50
	default:
		return 65
	}
}

// Restock 要透過索引修改 slice 裡的元素；range 給的 it 只是複本。
func Restock(items []Item, n int) {
	for i := range items {
		items[i].Stock += n
	}
}

func Table(n int) string {
	var b strings.Builder
	for i := 1; i <= n; i++ {
		if i > 1 {
			b.WriteString("\n")
		}
		for j := 1; j <= n; j++ {
			if j > 1 {
				b.WriteString("\t")
			}
			fmt.Fprintf(&b, "%d×%d=%d", i, j, i*j)
		}
	}
	return b.String()
}

func InvoicePrize(winning, ticket string) int {
	prizes := []int{200000, 40000, 10000, 4000, 1000, 200} // 末 8 碼到末 3 碼
	for i, prize := range prizes {
		digits := 8 - i
		if winning[8-digits:] == ticket[8-digits:] {
			return prize
		}
	}
	return 0
}

func CountKinds(s string) (han, letter, digit, other int) {
	for _, r := range s {
		switch {
		case unicode.Is(unicode.Han, r):
			han++
		case 'a' <= r && r <= 'z', 'A' <= r && r <= 'Z':
			letter++
		case '0' <= r && r <= '9':
			digit++
		default:
			other++
		}
	}
	return han, letter, digit, other
}

// TotalUntilNegative 在 switch 裡要離開 for，break 必須加上迴圈的標籤。
func TotalUntilNegative(nums []int) int {
	total := 0
loop:
	for _, n := range nums {
		switch {
		case n < 0:
			break loop
		default:
			total += n
		}
	}
	return total
}

func Spiral(n int) [][]int {
	if n < 1 {
		return [][]int{}
	}
	m := make([][]int, n)
	for i := range m {
		m[i] = make([]int, n)
	}
	top, bottom, left, right := 0, n-1, 0, n-1
	v := 1
	for top <= bottom && left <= right {
		for c := left; c <= right; c++ {
			m[top][c] = v
			v++
		}
		top++
		for r := top; r <= bottom; r++ {
			m[r][right] = v
			v++
		}
		right--
		for c := right; c >= left && top <= bottom; c-- {
			m[bottom][c] = v
			v++
		}
		bottom--
		for r := bottom; r >= top && left <= right; r-- {
			m[r][left] = v
			v++
		}
		left++
	}
	return m
}

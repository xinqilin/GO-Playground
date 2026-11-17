//go:build solution

package text

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func MaskName(name string) string {
	r := []rune(name)
	switch len(r) {
	case 0, 1:
		return name
	case 2:
		return string(r[0]) + "*"
	}
	return string(r[0]) + strings.Repeat("*", len(r)-2) + string(r[len(r)-1])
}

func NormalizePhone(s string) (string, bool) {
	s = strings.NewReplacer(" ", "", "-", "").Replace(s)
	if rest, ok := strings.CutPrefix(s, "+886"); ok {
		s = "0" + rest
	}
	if len(s) != 10 || !strings.HasPrefix(s, "09") {
		return "", false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return "", false
		}
	}
	return s, true
}

func ParseOrder(s string) (map[string]int, error) {
	order := make(map[string]int)
	if s == "" {
		return order, nil
	}
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		name, qtyText, ok := strings.Cut(part, "x")
		if !ok || name == "" {
			return nil, fmt.Errorf("格式錯誤：%q", part)
		}
		qty, err := strconv.Atoi(qtyText)
		if err != nil {
			return nil, fmt.Errorf("數量錯誤：%q：%w", part, err)
		}
		if qty <= 0 {
			return nil, errors.New("數量必須是正整數：" + part)
		}
		order[name] += qty
	}
	return order, nil
}

// Truncate 要以「字」為單位計算，所以先轉成 []rune。
func Truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}

// IsPalindrome 比較的是 rune，不是 byte：中文字的 byte 倒過來就不是同一個字了。
func IsPalindrome(s string) bool {
	r := []rune(s)
	for i := 0; i < len(r)/2; i++ {
		if r[i] != r[len(r)-1-i] {
			return false
		}
	}
	return true
}

func FormatMenu(items []MenuItem) string {
	width := func(s string) int {
		w := 0
		for _, r := range s {
			w += RuneWidth(r)
		}
		return w
	}
	nameW, priceW := 0, 0
	for _, it := range items {
		nameW = max(nameW, width(it.Name))
		priceW = max(priceW, len(strconv.Itoa(it.Price)))
	}
	lines := make([]string, 0, len(items))
	for _, it := range items {
		pad := strings.Repeat(" ", nameW-width(it.Name))
		lines = append(lines, fmt.Sprintf("%s%s  %*d", it.Name, pad, priceW, it.Price))
	}
	return strings.Join(lines, "\n")
}

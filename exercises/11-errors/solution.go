//go:build solution

package orders

import (
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"strconv"
	"strings"
)

// ParsePickup 先檢查格式，再檢查營業時間；兩種錯誤包住不同的 sentinel。
func ParsePickup(s string) (int, error) {
	hs, ms, ok := strings.Cut(s, ":")
	if !ok {
		return 0, fmt.Errorf("取餐時間 %q：%w", s, ErrTimeFormat)
	}
	h, err1 := strconv.Atoi(hs)
	m, err2 := strconv.Atoi(ms)
	if err1 != nil || err2 != nil || h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, fmt.Errorf("取餐時間 %q：%w", s, ErrTimeFormat)
	}
	minutes := h*60 + m
	if minutes < 11*60 || minutes > 20*60 {
		return 0, fmt.Errorf("取餐時間 %s：%w", s, ErrClosed)
	}
	return minutes, nil
}

func (e *LimitError) Error() string {
	return fmt.Sprintf("%s 每人限購 %d 份", e.Item, e.Max)
}

// Buy 沒有問題時直接 return nil，不經過 *LimitError 型別的變數（第 10 章的 nil 陷阱）。
func Buy(item string, qty int, limits map[string]int) error {
	limit, ok := limits[item]
	if ok && qty > limit {
		return &LimitError{Item: item, Max: limit}
	}
	return nil
}

// CheckCart 收集所有錯誤；errors.Join 會略過 nil，全部是 nil 時回傳 nil。
func CheckCart(cart map[string]int, limits map[string]int) error {
	var errs []error
	for _, item := range slices.Sorted(maps.Keys(cart)) {
		errs = append(errs, Buy(item, cart[item], limits))
	}
	return errors.Join(errs...)
}

func Retry(times int, op func() error) error {
	times = max(times, 1)
	var err error
	for range times {
		err = op()
		if err == nil {
			return nil
		}
		if errors.Is(err, ErrPermanent) {
			return err
		}
	}
	return fmt.Errorf("試了 %d 次：%w", times, err)
}

// Causes 用 type switch 分辨兩種 Unwrap 方法，再遞迴走進每一個被包著的錯誤。
func Causes(err error) []string {
	if err == nil {
		return nil
	}
	out := []string{err.Error()}
	switch x := err.(type) {
	case interface{ Unwrap() error }:
		out = append(out, Causes(x.Unwrap())...)
	case interface{ Unwrap() []error }:
		for _, e := range x.Unwrap() {
			out = append(out, Causes(e)...)
		}
	}
	return out
}

// IsSoldOut 用 errors.Is：比的是錯誤值本身，而且會一層一層拆開包裝去找。
func IsSoldOut(err error) bool {
	return errors.Is(err, ErrSoldOut)
}

// SaveLines 讓回傳值有名字，defer 裡才能把 Close 的錯誤併進回傳值。
func SaveLines(w io.WriteCloser, lines []string) (err error) {
	defer func() {
		err = errors.Join(err, w.Close())
	}()
	for _, l := range lines {
		if _, err = fmt.Fprintln(w, l); err != nil {
			return err
		}
	}
	return nil
}

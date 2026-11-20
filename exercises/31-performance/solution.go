//go:build solution

package rush

import (
	"bytes"
	"math"
	"slices"
	"strconv"
	"strings"
)

func Summarize(samples []float64) (median, spread float64) {
	if len(samples) == 0 {
		return 0, 0
	}
	s := slices.Sorted(slices.Values(samples))
	n := len(s)
	median = s[n/2]
	if n%2 == 0 {
		median = (s[n/2-1] + s[n/2]) / 2
	}
	for _, x := range s {
		spread = max(spread, math.Abs(x-median)/median*100)
	}
	return median, spread
}

func AppendLine(dst []byte, dish string, qty, price int) []byte {
	dst = append(dst, dish...)
	dst = append(dst, " x"...)
	dst = strconv.AppendInt(dst, int64(qty), 10)
	dst = append(dst, ' ')
	dst = strconv.AppendInt(dst, int64(qty*price), 10)
	return append(dst, '\n')
}

func ParseOrder(line string) (Order, error) {
	t, rest, ok1 := strings.Cut(line, ",")
	dish, rest, ok2 := strings.Cut(rest, ",")
	q, p, ok3 := strings.Cut(rest, ",")
	if !ok1 || !ok2 || !ok3 || !validTime(t) || dish == "" || strings.Contains(p, ",") {
		return Order{}, ErrBadLine
	}
	qty, err1 := strconv.Atoi(q)
	price, err2 := strconv.Atoi(p)
	if err1 != nil || err2 != nil || qty < 1 || qty > 99 || price <= 0 {
		return Order{}, ErrBadLine
	}
	return Order{Time: t, Dish: dish, Qty: qty, Price: price}, nil
}

func validTime(t string) bool {
	if len(t) != 5 || t[2] != ':' {
		return false
	}
	for _, i := range []int{0, 1, 3, 4} {
		if t[i] < '0' || t[i] > '9' {
			return false
		}
	}
	return true
}

func (p *BufferPool) Get() *bytes.Buffer {
	b, _ := p.pool.Get().(*bytes.Buffer)
	if b == nil {
		return new(bytes.Buffer)
	}
	b.Reset()
	return b
}

func (p *BufferPool) Put(b *bytes.Buffer) {
	if p.MaxCap > 0 && b.Cap() > p.MaxCap {
		return
	}
	p.pool.Put(b)
}

// Receipt 先用 AppendLine 在 stack 上的小陣列裡量出每一行的長度，再一次配置好。
func Receipt(items []Item) string {
	var scratch [64]byte
	size, total := 0, 0
	for _, it := range items {
		size += len(AppendLine(scratch[:0], "", it.Qty, it.Price)) + len(it.Dish)
		total += it.Qty * it.Price
	}
	totalLine := strconv.AppendInt(append(scratch[:0], "合計 "...), int64(total), 10)
	size += len(totalLine) + 1

	var b strings.Builder
	b.Grow(size)
	for _, it := range items {
		b.Write(AppendLine(scratch[:0], it.Dish, it.Qty, it.Price))
	}
	totalLine = strconv.AppendInt(append(scratch[:0], "合計 "...), int64(total), 10)
	b.Write(totalLine)
	b.WriteByte('\n')
	return b.String()
}

// FirstLine 複製出第一行：結果有自己的記憶體，data 可以被回收或重複使用。
func FirstLine(data []byte) []byte {
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		data = data[:i]
	}
	return bytes.Clone(data)
}

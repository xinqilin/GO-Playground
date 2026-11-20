//go:build solution

package inout

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func Percent(part, total int) string {
	if total < 1 {
		return "0.0%"
	}
	return fmt.Sprintf("%.1f%%", float64(part)*100/float64(total))
}

func HexColor(r, g, b uint8) string {
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

// String 從個位數往左，每三位插一個逗號。
// 要用 m.String() 以外的方式取得數字（這裡是 strconv），不然 fmt 會再呼叫 String，無限遞迴。
func (m Money) String() string {
	n := int(m)
	sign := ""
	if n < 0 {
		sign, n = "-", -n
	}
	digits := strconv.Itoa(n)
	var sb strings.Builder
	for i, d := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			sb.WriteByte(',')
		}
		sb.WriteRune(d)
	}
	return sign + "NT$" + sb.String()
}

func WriteReceipt(w io.Writer, items []Item) error {
	total := 0
	for _, it := range items {
		sub := it.Qty * it.Price
		total += sub
		if _, err := fmt.Fprintf(w, "%5d  %s x%d\n", sub, it.Name, it.Qty); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(w, "-----"); err != nil {
		return err
	}
	_, err := fmt.Fprintf(w, "%5d  合計\n", total)
	return err
}

func ParseOrders(r io.Reader) (int, error) {
	sc := bufio.NewScanner(r)
	total, lineNo := 0, 0
	for sc.Scan() {
		lineNo++
		line := sc.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		fields := strings.Fields(trimmed)
		if len(fields) != 3 {
			return 0, fmt.Errorf("第 %d 行格式錯誤：%q", lineNo, line)
		}
		qty, err1 := strconv.Atoi(fields[1])
		price, err2 := strconv.Atoi(fields[2])
		if err1 != nil || err2 != nil {
			return 0, fmt.Errorf("第 %d 行格式錯誤：%q", lineNo, line)
		}
		total += qty * price
	}
	if err := sc.Err(); err != nil {
		return 0, err
	}
	return total, nil
}

// CallNumbers 的結尾也要寫進 bw：直接寫到 w 的東西會比還在 bw 裡排隊的先出去。
func CallNumbers(w io.Writer, names []string) error {
	bw := bufio.NewWriter(w)
	for _, name := range names {
		fmt.Fprintf(bw, "請 %s 到櫃台\n", name)
	}
	fmt.Fprintln(bw, "== 以上 ==")
	return bw.Flush()
}

func LongestLine(r io.Reader) (int, error) {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024+1)
	longest := 0
	for sc.Scan() {
		longest = max(longest, len(sc.Bytes()))
	}
	if err := sc.Err(); err != nil {
		return 0, err
	}
	return longest, nil
}

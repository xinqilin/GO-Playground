package main

import (
	"errors"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// Order 是便當店訂單紀錄的一行：「08:05,排骨飯,2,90」。
type Order struct {
	Time  string
	Dish  string
	Qty   int
	Price int
}

var errBadLine = errors.New("訂單格式不對")

var dishes = []struct {
	name  string
	price int
}{{"排骨飯", 90}, {"雞腿飯", 100}, {"滷肉飯", 50}, {"魚排飯", 95}, {"素食便當", 80}}

// sampleLog 產生 n 行固定的訂單紀錄，不用亂數，每次結果都一樣。
func sampleLog(n int) string {
	var b strings.Builder
	for i := range n {
		d := dishes[i*7%len(dishes)]
		minute := 10*60 + i%(4*60) // 10:00 到 13:59
		fmt.Fprintf(&b, "%02d:%02d,%s,%d,%d\n", minute/60, minute%60, d.name, i%3+1, d.price)
	}
	return b.String()
}

var lineRE = regexp.MustCompile(`^(\d{2}:\d{2}),([^,]+),(\d+),(\d+)$`)

// parseRegexp 用 regexp 解析一行：好寫，但每一行都要跑一次比對。
func parseRegexp(line string) (Order, error) {
	m := lineRE.FindStringSubmatch(line)
	if m == nil {
		return Order{}, fmt.Errorf("%w：%q", errBadLine, line)
	}
	qty, err1 := strconv.Atoi(m[3]) // \d+ 擋掉了正負號，但擋不住太長、放不進 int 的數字
	price, err2 := strconv.Atoi(m[4])
	if err1 != nil || err2 != nil {
		return Order{}, fmt.Errorf("%w：%q", errBadLine, line)
	}
	return Order{Time: m[1], Dish: m[2], Qty: qty, Price: price}, nil
}

// parseCut 用 strings.Cut 一段一段切：子字串和原本的字串共用記憶體，不用另外配置。
// regexp 檢查的每一條規則都要自己再寫一次，少寫一條，兩個版本的結果就不一樣了。
func parseCut(line string) (Order, error) {
	t, rest, ok1 := strings.Cut(line, ",")
	dish, rest, ok2 := strings.Cut(rest, ",")
	q, p, ok3 := strings.Cut(rest, ",")
	qty, err1 := strconv.Atoi(q)
	price, err2 := strconv.Atoi(p)
	clock := len(t) == 5 && t[2] == ':' && digits(t[:2]) && digits(t[3:])
	if !ok1 || !ok2 || !ok3 || !clock || dish == "" || !digits(q) || !digits(p) || err1 != nil || err2 != nil {
		return Order{}, fmt.Errorf("%w：%q", errBadLine, line)
	}
	return Order{Time: t, Dish: dish, Qty: qty, Price: price}, nil
}

// digits 回傳 s 是不是由一個以上的 0–9 組成。Atoi 還接受 + 和 - 號，regexp 的 \d 不接受。
func digits(s string) bool {
	for i := range len(s) {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return s != ""
}

// revenue 算出每道菜的營業額，parse 決定用哪一種方法解析每一行。
func revenue(log string, parse func(string) (Order, error)) (map[string]int, error) {
	totals := make(map[string]int)
	for line := range strings.Lines(log) {
		o, err := parse(strings.TrimSuffix(line, "\n"))
		if err != nil {
			return nil, err
		}
		totals[o.Dish] += o.Qty * o.Price
	}
	return totals, nil
}

// revenueDemo：兩種解析方法算出來的結果要一模一樣，比速度才有意義。
func revenueDemo() {
	log := sampleLog(1000)
	slow, err1 := revenue(log, parseRegexp)
	fast, err2 := revenue(log, parseCut)
	for _, dish := range slices.Sorted(maps.Keys(fast)) {
		fmt.Printf("%s %d 元\n", dish, fast[dish])
	}
	fmt.Println("兩種結果相同：", maps.Equal(slow, fast), err1, err2)
	_, err := parseCut("10:05,排骨飯,兩個,90")
	fmt.Println(err)
}

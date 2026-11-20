package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// OrderLine 是點餐單上的一行。
type OrderLine struct {
	Dish string
	Qty  int
	Note string
}

// parseOrder 解析櫃台手打的點餐單。每一行是「品項 [x數量] [備註…]」，
// 空行和 # 開頭的註解略過；數量沒寫就是 1。
func parseOrder(text string) ([]OrderLine, error) {
	var order []OrderLine
	n := 0
	for line := range strings.Lines(text) {
		n++
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		ol := OrderLine{Dish: fields[0], Qty: 1}
		rest := fields[1:]
		if len(rest) > 0 {
			if q, ok := strings.CutPrefix(rest[0], "x"); ok {
				qty, err := strconv.Atoi(q)
				if err != nil {
					return nil, fmt.Errorf("第 %d 行：數量 %q：%w", n, q, err)
				}
				if qty < 1 {
					return nil, fmt.Errorf("第 %d 行：數量至少是 1", n)
				}
				ol.Qty = qty
				rest = rest[1:]
			}
		}
		ol.Note = strings.Join(rest, " ")
		order = append(order, ol)
	}
	return order, nil
}

func parseDemo() {
	const ticket = `# 第 7 桌
排骨飯 x2 少飯
滷肉飯

珍珠奶茶 x1 微糖, 去冰
`
	order, err := parseOrder(ticket)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, ol := range order {
		fmt.Printf("%s %d 份 %q\n", ol.Dish, ol.Qty, ol.Note)
	}

	_, err = parseOrder("排骨飯\n雞腿飯 x二\n")
	fmt.Println(err)
}

// numErrorDemo strconv 的錯誤是 *strconv.NumError，分得出「不是數字」和「超出範圍」。
func numErrorDemo() {
	for _, s := range []string{"二", "99999999999999999999"} {
		_, err := strconv.Atoi(s)
		ne, _ := errors.AsType[*strconv.NumError](err)
		fmt.Println(ne.Func, ne.Num, errors.Is(err, strconv.ErrSyntax), errors.Is(err, strconv.ErrRange))
	}

	_, err := strconv.ParseInt("300", 10, 8) // 第三個參數是位元數：int8 放不下 300
	fmt.Println(err)

	for _, s := range []string{"true", "1", "T", "yes"} {
		b, err := strconv.ParseBool(s)
		fmt.Println(s, b, err == nil)
	}
}

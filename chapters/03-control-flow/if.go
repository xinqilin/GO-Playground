package main

import (
	"errors"
	"fmt"
	"strconv"
)

// ifElse 由上往下檢查條件，執行第一個成立的區塊，其餘跳過。
func ifElse() {
	for _, temp := range []int{33, 24, 12} {
		if temp >= 30 {
			fmt.Println(temp, "度：開冷氣")
		} else if temp >= 20 {
			fmt.Println(temp, "度：開窗就好")
		} else {
			fmt.Println(temp, "度：穿外套")
		}
	}
}

// ifInit 在 if 前面先做一件事；宣告的 n、err 只活在這組 if / else 裡。
func ifInit() {
	for _, input := range []string{"3", "三", "0"} {
		if n, err := strconv.Atoi(input); err != nil {
			fmt.Printf("%q 不是數字\n", input)
		} else if n <= 0 {
			fmt.Println("至少要訂 1 個")
		} else {
			fmt.Println("訂", n, "個便當，共", n*85, "元")
		}
		// 這裡已經看不到 n 和 err
	}
}

type order struct {
	paid    bool
	address string
	items   int
}

// canShipNested 把條件一層一層包起來：每多一個檢查，主要邏輯就往右縮一層。
func canShipNested(o order) string {
	if o.paid {
		if o.address != "" {
			if o.items > 0 {
				return "可以出貨"
			} else {
				return "購物車是空的"
			}
		} else {
			return "缺少地址"
		}
	} else {
		return "還沒付款"
	}
}

// canShip 用 guard clause：先把不合格的情況處理掉並 return，正常流程留在最外層。
func canShip(o order) string {
	if !o.paid {
		return "還沒付款"
	}
	if o.address == "" {
		return "缺少地址"
	}
	if o.items == 0 {
		return "購物車是空的"
	}
	return "可以出貨"
}

// guardClause 證明兩種寫法的結果一樣，差別只在好不好讀。
func guardClause() {
	orders := []order{
		{paid: false},
		{paid: true},
		{paid: true, address: "台北市信義路 1 號"},
		{paid: true, address: "台北市信義路 1 號", items: 2},
	}
	for _, o := range orders {
		fmt.Println(canShipNested(o), "|", canShip(o))
	}
}

// shadowing 示範在內層區塊用 := 會宣告一個同名的新變數，外層的變數不受影響。
func shadowing() {
	limit := 100
	if limit > 50 {
		limit := 50 // 新的 limit，只活在這個區塊
		fmt.Println("區塊裡：", limit)
	}
	fmt.Println("區塊外：", limit)
}

var settings = map[string]string{"theme": "dark"}

// readSetting 找不到設定時回傳錯誤。
func readSetting(key string) (string, error) {
	v, ok := settings[key]
	if !ok {
		return "", errors.New("找不到設定 " + key)
	}
	return v, nil
}

// loadTheme 想把讀到的值存進外層的 theme，卻在 if 裡用 := 宣告了新的 theme 和 err。
func loadTheme(fromSettings bool) string {
	theme := "light"
	if fromSettings {
		theme, err := readSetting("theme")
		if err != nil {
			return "讀取失敗"
		}
		fmt.Println("if 裡讀到：", theme)
	}
	return theme
}

// loadThemeFixed 先宣告 err，再用 = 賦值，寫進的就是外層的 theme。
func loadThemeFixed(fromSettings bool) string {
	theme := "light"
	if fromSettings {
		var err error
		theme, err = readSetting("theme")
		if err != nil {
			return "讀取失敗"
		}
	}
	return theme
}

func shadowPitfall() {
	fmt.Println("錯誤寫法回傳：", loadTheme(true))
	fmt.Println("正確寫法回傳：", loadThemeFixed(true))
}

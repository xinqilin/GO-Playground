package main

import (
	"fmt"
	"regexp"
)

// 正規表示式在套件層級編譯一次，之後重複使用。MustCompile 在寫錯時 panic，
// 程式一啟動就會發現，不會等到第一次用到。
var (
	invoiceNo = regexp.MustCompile(`^[A-Z]{2}-\d{8}$`)
	looseNo   = regexp.MustCompile(`[A-Z]{2}-\d{8}`) // 沒有 ^ 和 $：只要字串裡「有一段」符合就算
)

// matchDemo 同一組字串，有沒有錨點（^ 和 $）結果不一樣。
func matchDemo() {
	for _, s := range []string{"AB-12345678", "AB-123456789", "發票 AB-12345678", "ab-12345678"} {
		fmt.Printf("%q：完整比對 %v，部分比對 %v\n", s, invoiceNo.MatchString(s), looseNo.MatchString(s))
	}
}

// trip 解析悠遊卡的交易紀錄。(?<name>…) 是命名群組；(?:…) 只分組、不擷取；最後的 ? 表示整段可有可無。
var trip = regexp.MustCompile(`^(?<date>\d{4}/\d{2}/\d{2}) (?<time>\d{2}:\d{2}) (?<kind>進站|出站) (?<station>\S+)(?: 扣款 (?<fare>\d+) 元)?$`)

// submatchDemo 用 FindStringSubmatch 取出每個群組，再用 SubexpIndex 依名稱找位置。
func submatchDemo() {
	for _, line := range []string{
		"2026/11/19 08:15 進站 台北車站",
		"2026/11/19 08:41 出站 市政府 扣款 25 元",
		"2026/11/19 08:41 轉乘 市政府",
	} {
		m := trip.FindStringSubmatch(line)
		if m == nil {
			fmt.Printf("看不懂：%s\n", line)
			continue
		}
		fmt.Printf("%s %s %s 車資=%q\n",
			m[trip.SubexpIndex("time")], m[trip.SubexpIndex("kind")],
			m[trip.SubexpIndex("station")], m[trip.SubexpIndex("fare")])
	}

	amounts := regexp.MustCompile(`NT\$[\d,]+`)
	fmt.Printf("%q\n", amounts.FindAllString("早餐 NT$85，午餐 NT$120，晚餐 NT$1,280", -1))
}

package main

import (
	"fmt"
	"unicode/utf8"
)

// switchValue 拿一個值去比對每個 case；一個 case 可以列好幾個值。
func switchValue() {
	for _, day := range []string{"六", "二", "日"} {
		switch day {
		case "六", "日":
			fmt.Printf("星期%s：假日\n", day)
		default:
			fmt.Printf("星期%s：上班\n", day)
		}
	}
}

// switchInit 和 if 一樣，switch 前面可以先放一個簡短敘述。
func switchInit() {
	for _, name := range []string{"王明", "陳小美", "歐陽靖文", "李"} {
		switch n := utf8.RuneCountInString(name); n {
		case 2, 3:
			fmt.Printf("%s：一般長度\n", name)
		case 4:
			fmt.Printf("%s：可能是複姓\n", name)
		default:
			fmt.Printf("%s：%d 個字，請再確認\n", name, n)
		}
	}
}

// aqiLevel 用沒有條件的 switch 取代一長串 if-else，由上往下找第一個成立的 case。
func aqiLevel(aqi int) string {
	switch {
	case aqi <= 50:
		return "良好"
	case aqi <= 100:
		return "普通"
	case aqi <= 150:
		return "對敏感族群不健康"
	case aqi <= 200:
		return "對所有族群不健康"
	case aqi <= 300:
		return "非常不健康"
	default:
		return "危害"
	}
}

func switchTrue() {
	for _, v := range []int{35, 120, 320} {
		fmt.Printf("AQI %d：%s\n", v, aqiLevel(v))
	}
}

// check 印出自己被呼叫了，用來觀察 case 的求值順序。
func check(label string, ok bool) bool {
	fmt.Println("檢查", label)
	return ok
}

// switchOrder 證明 case 由上往下求值，找到成立的就停，後面的 case 不會被求值。
func switchOrder() {
	n := 7
	switch {
	case check("n > 10", n > 10):
		fmt.Println("大於 10")
	case check("n > 5", n > 5):
		fmt.Println("大於 5")
	case check("n > 0", n > 0):
		fmt.Println("大於 0")
	}
}

// perks 用 fallthrough 讓高等級的會員也拿到低等級的優惠。
func perks(level string) []string {
	var got []string
	switch level {
	case "金卡":
		got = append(got, "機場接送")
		fallthrough
	case "銀卡":
		got = append(got, "免運")
		fallthrough
	case "一般":
		got = append(got, "生日禮")
	}
	return got
}

func switchFallthrough() {
	for _, lv := range []string{"金卡", "銀卡", "一般", "訪客"} {
		fmt.Println(lv, perks(lv))
	}
}

// describe 用 type switch 依照值的實際型別走不同分支，第 10 章會詳細介紹。
func describe(x any) string {
	switch v := x.(type) {
	case nil:
		return "nil"
	case int:
		return fmt.Sprintf("整數 %d，兩倍是 %d", v, v*2)
	case string:
		return fmt.Sprintf("字串 %q，長度 %d", v, len(v))
	default:
		return fmt.Sprintf("其他型別 %T", v)
	}
}

func typeSwitch() {
	for _, x := range []any{42, "便當", 3.5, nil} {
		fmt.Println(describe(x))
	}
}

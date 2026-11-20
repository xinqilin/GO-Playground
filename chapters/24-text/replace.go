package main

import (
	"fmt"
	"regexp"
	"strconv"
)

var slashDate = regexp.MustCompile(`(\d{4})/(\d{2})/(\d{2})`)

// replaceDemo 取代字串裡用 $1 或 ${1} 引用群組。名稱後面緊接著文字時，一定要用大括號。
func replaceDemo() {
	s := "訂位日期 2026/11/19，取消期限 2026/11/17"
	fmt.Println(slashDate.ReplaceAllString(s, "$1-$2-$3"))
	fmt.Printf("%q\n", slashDate.ReplaceAllString(s, "$1年$2月$3日"))
	fmt.Println(slashDate.ReplaceAllString(s, "${1}年${2}月${3}日"))
}

var rocDate = regexp.MustCompile(`民國(\d+)年`)

// replaceFuncDemo 取代的結果需要計算時用 ReplaceAllStringFunc：民國年換成西元年。
func replaceFuncDemo() {
	s := "本店民國68年開幕，民國115年搬到現址"
	out := rocDate.ReplaceAllStringFunc(s, func(m string) string {
		y, err := strconv.Atoi(rocDate.FindStringSubmatch(m)[1])
		if err != nil {
			return m // 數字太大轉不過來時保留原文
		}
		return strconv.Itoa(y+1911) + "年"
	})
	fmt.Println(out)

	price := regexp.QuoteMeta("NT$1.5") // 使用者輸入的文字裡有 $ 和 .，要當成一般字元
	fmt.Println(price, regexp.MustCompile(price).MatchString("特價 NT$1.5 起"))
}

// re2Demo Go 的 regexp 是 RE2 語法：不支援 lookahead 和 backreference，寫了會在編譯時就錯。
func re2Demo() {
	for _, expr := range []string{`\d+(?=元)`, `(\w)\1`, `[A-Z`} {
		_, err := regexp.Compile(expr)
		fmt.Println(err)
	}
}

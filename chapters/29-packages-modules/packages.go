package main

import (
	"fmt"
	htmltemplate "html/template"
	"os"
	"text/template"

	"github.com/xinqilin/GO-Playground/chapters/29-packages-modules/bento"
	"github.com/xinqilin/GO-Playground/chapters/29-packages-modules/internal/receipt"
)

// exportDemo 使用另一個 package：只看得到大寫開頭的名稱。
func exportDemo() {
	m := bento.New()
	fmt.Println(m.Add("排骨飯", 90), m.Add("滷肉飯", 50))
	fmt.Println(m.Add("空氣便當", 0))
	// m.prices["排骨飯"] = -1 // 編譯錯誤：prices 沒有匯出
	// bento.valid(-1)        // 編譯錯誤：valid 沒有匯出

	price, ok := m.Price("排骨飯")
	fmt.Println(price, ok, m.Dishes())

	fmt.Print(receipt.Format([]receipt.Line{
		{Dish: "排骨飯", Qty: 2, Price: price},
		{Dish: "滷肉飯", Qty: 1, Price: 50},
	}))
}

// aliasDemo：兩個 package 的名字都叫 template，其中一個要取別名。
func aliasDemo() {
	const note = "備註：<b>不要辣</b>\n"
	text := template.Must(template.New("t").Parse("{{.}}"))
	html := htmltemplate.Must(htmltemplate.New("h").Parse("{{.}}"))
	text.Execute(os.Stdout, note)
	html.Execute(os.Stdout, note)
}

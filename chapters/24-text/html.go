package main

import (
	"fmt"
	htmltemplate "html/template"
	"os"
	"text/template"
)

// Review 是顧客留下的評論，每個欄位都是使用者輸入的，不能相信。
type Review struct {
	Name    string
	Website string
	Comment string
}

// reviewPage 裡同一種 {{.}}，出現在 HTML 內容、屬性裡的網址、JavaScript 裡，需要的跳脫方式都不一樣。
const reviewPage = `<p><a href="{{.Website}}">{{.Name}}</a>：{{.Comment}}</p>
<button onclick="reply({{.Name}})">回覆</button>
`

// escapeDemo 同一個範本、同一份惡意的輸入，用 text/template 和 html/template 各執行一次。
func escapeDemo() {
	r := Review{
		Name:    `<b>阿明</b>`,
		Website: `javascript:alert(1)`,
		Comment: `好吃<script>steal()</script>`,
	}

	fmt.Println("text/template：")
	t := template.Must(template.New("review").Parse(reviewPage))
	if err := t.Execute(os.Stdout, r); err != nil {
		fmt.Println(err)
	}

	fmt.Println("html/template：")
	h := htmltemplate.Must(htmltemplate.New("review").Parse(reviewPage))
	if err := h.Execute(os.Stdout, r); err != nil {
		fmt.Println(err)
	}
}

// trustedDemo template.HTML 告訴 html/template「這段已經是安全的 HTML，不要跳脫」。
// 只能用在程式自己產生的內容，絕對不能包使用者的輸入。
func trustedDemo() {
	h := htmltemplate.Must(htmltemplate.New("badge").Parse("<p>{{.}}</p>\n"))
	for _, v := range []any{"<em>招牌</em>", htmltemplate.HTML("<em>招牌</em>")} {
		if err := h.Execute(os.Stdout, v); err != nil {
			fmt.Println(err)
		}
	}
}

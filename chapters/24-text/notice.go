package main

import (
	"embed"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// money 把 1280 寫成 1,280：從個位數往前，每三位插一個逗號。
func money(n int) string {
	s := strconv.Itoa(n)
	var b strings.Builder
	for i, d := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(d)
	}
	return b.String()
}

// funcs 是範本裡可以呼叫的函式，要在 Parse 之前用 Funcs 加進去。
var funcs = template.FuncMap{
	"money": money,
	"join":  strings.Join,
}

// base 是一般的公告；holiday 複製一份 base，再把 body 換成公休的內容。
var (
	base    = template.Must(template.New("").Funcs(funcs).ParseFS(templateFS, "templates/notice.tmpl"))
	holiday = template.Must(template.Must(base.Clone()).ParseFS(templateFS, "templates/holiday.tmpl"))
)

// Notice 是店門口的公告。
type Notice struct {
	Shop     string
	Hours    []string
	MinOrder int
	Date     string
}

func noticeDemo() {
	n := Notice{Shop: "阿嬤便當", Hours: []string{"10:30-14:00", "17:00-20:30"}, MinOrder: 1500, Date: "11/24（一）"}
	for _, t := range []*template.Template{base, holiday} {
		if err := t.ExecuteTemplate(os.Stdout, "notice", n); err != nil {
			fmt.Println(err)
		}
	}
}

//go:build solution

package memo

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"path"
	"regexp"
	"strconv"
	"strings"
)

// ParseKV 用 SplitSeq 切段落、Cut 切 key 和 value。
func ParseKV(s string) (map[string]string, error) {
	kv := map[string]string{}
	for part := range strings.SplitSeq(s, ";") {
		if strings.TrimSpace(part) == "" {
			continue
		}
		k, v, ok := strings.Cut(part, "=")
		if !ok {
			return nil, fmt.Errorf("%q 沒有 =", strings.TrimSpace(part))
		}
		k = strings.TrimSpace(k)
		if k == "" {
			return nil, fmt.Errorf("%q 的 key 是空的", strings.TrimSpace(part))
		}
		if _, dup := kv[k]; dup {
			return nil, fmt.Errorf("%s 重複了", k)
		}
		kv[k] = strings.TrimSpace(v)
	}
	return kv, nil
}

// reviewsTmpl 用 html/template：名字和評論會依照所在的位置自動跳脫。
var reviewsTmpl = template.Must(template.New("reviews").Parse(`{{if .}}<ul>
{{range .}}<li><b>{{.Name}}</b> {{.Stars}} 顆星：{{.Comment}}</li>
{{end}}</ul>
{{else}}<p>還沒有評論</p>
{{end}}`))

func RenderReviews(w io.Writer, reviews []Review) error {
	return reviewsTmpl.Execute(w, reviews)
}

var (
	// \b 是英數字和其他字元的交界：前後連著數字或英文字母就不算。
	phoneInText = regexp.MustCompile(`\b(09\d{2})[- ]?\d{3}[- ]?(\d{3})\b`)
	email       = regexp.MustCompile(`\b([A-Za-z0-9])[A-Za-z0-9._%+-]*@([A-Za-z0-9.-]+\.[A-Za-z]{2,})\b`)
)

// Mask 手機號碼要統一格式，所以用 Func 版本；email 直接用 ${1} 引用群組。
func Mask(s string) string {
	s = phoneInText.ReplaceAllStringFunc(s, func(m string) string {
		g := phoneInText.FindStringSubmatch(m)
		return g[1] + "-***-" + g[2]
	})
	return email.ReplaceAllString(s, "${1}***@${2}")
}

var invoiceLine = regexp.MustCompile(`(?<no>[A-Z]{2}-\d{8})\s+(?<y>\d{4})[/-](?<m>\d{2})[/-](?<d>\d{2})\s+NT\$(?<amount>\d[\d,]*)`)

// FindInvoices 用 FindAllStringSubmatch 一次找出全部，再依群組名稱取值。
func FindInvoices(text string) []Invoice {
	var out []Invoice
	for _, m := range invoiceLine.FindAllStringSubmatch(text, -1) {
		get := func(name string) string { return m[invoiceLine.SubexpIndex(name)] }
		amount, err := strconv.Atoi(strings.ReplaceAll(get("amount"), ",", ""))
		if err != nil {
			continue
		}
		out = append(out, Invoice{
			Number: get("no"),
			Date:   get("y") + "-" + get("m") + "-" + get("d"),
			Amount: amount,
		})
	}
	return out
}

// Site 每個頁面一份獨立的範本：版型 Clone 一份，再加上那個頁面的 content。
type Site struct {
	pages map[string]*template.Template
}

func NewSite(fsys fs.FS) (*Site, error) {
	layout, err := template.ParseFS(fsys, "layout.html")
	if err != nil {
		return nil, err
	}
	files, err := fs.Glob(fsys, "pages/*.html")
	if err != nil {
		return nil, err
	}
	s := &Site{pages: map[string]*template.Template{}}
	for _, f := range files {
		t, err := layout.Clone()
		if err != nil {
			return nil, err
		}
		if _, err := t.ParseFS(fsys, f); err != nil {
			return nil, err
		}
		s.pages[strings.TrimSuffix(path.Base(f), ".html")] = t
	}
	return s, nil
}

// Render 先寫到緩衝區，成功了才寫到 w。
func (s *Site) Render(w io.Writer, name string, data any) error {
	t, ok := s.pages[name]
	if !ok {
		return errors.New("沒有這個頁面：" + name)
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "layout.html", data); err != nil {
		return err
	}
	_, err := buf.WriteTo(w)
	return err
}

// phonePattern 加上 ^ 和 $：整個字串都要符合，不是「裡面有一段符合」。
var phonePattern = regexp.MustCompile(`^09\d{2}-?\d{3}-?\d{3}$`)

func ValidPhone(s string) bool {
	return phonePattern.MatchString(s)
}

package memo

import (
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
)

func TestParseKV(t *testing.T) {
	got, err := ParseKV(" name = 阿嬤便當 ; phone=02-2345-6789;note=a=b ;")
	want := map[string]string{"name": "阿嬤便當", "phone": "02-2345-6789", "note": "a=b"}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("ParseKV = %v（%v），want %v", got, err, want)
	}
	if got, err := ParseKV(" ; "); err != nil || len(got) != 0 {
		t.Errorf("只有空的段落：%v（%v），want 空的 map", got, err)
	}
	if got, err := ParseKV("open="); err != nil || got["open"] != "" || len(got) != 1 {
		t.Errorf("value 是空的也可以：%v（%v）", got, err)
	}
	for _, in := range []string{"name", "=阿嬤便當", "a=1; a=2", "a=1; b"} {
		if _, err := ParseKV(in); err == nil {
			t.Errorf("ParseKV(%q) 應該回傳錯誤", in)
		}
	}
}

func TestRenderReviews(t *testing.T) {
	var b strings.Builder
	err := RenderReviews(&b, []Review{
		{"阿明", 5, "好吃"},
		{"<i>小美</i>", 3, `太鹹了 & "太油"`},
	})
	want := "<ul>\n" +
		"<li><b>阿明</b> 5 顆星：好吃</li>\n" +
		"<li><b>&lt;i&gt;小美&lt;/i&gt;</b> 3 顆星：太鹹了 &amp; &#34;太油&#34;</li>\n" +
		"</ul>\n"
	if err != nil || b.String() != want {
		t.Errorf("RenderReviews（%v）寫出\n%s\nwant\n%s", err, b.String(), want)
	}

	b.Reset()
	if err := RenderReviews(&b, nil); err != nil || b.String() != "<p>還沒有評論</p>\n" {
		t.Errorf("沒有評論（%v）寫出 %q，want %q", err, b.String(), "<p>還沒有評論</p>\n")
	}
}

func TestMask(t *testing.T) {
	tests := map[string]string{
		"請回電 0912-345-678，謝謝":        "請回電 0912-***-678，謝謝",
		"手機0912345678或0987 654 321":  "手機0912-***-678或0987-***-321",
		"信箱 bill@example.com.tw 收到了": "信箱 b***@example.com.tw 收到了",
		"訂單 A091234567890 已出貨":       "訂單 A091234567890 已出貨",
		"市話 02-2345-6789 不用遮":        "市話 02-2345-6789 不用遮",
		"聯絡 0912-345-678 或 a@b.co":   "聯絡 0912-***-678 或 a***@b.co",
		"沒有個資的句子":                    "沒有個資的句子",
	}
	for in, want := range tests {
		if got := Mask(in); got != want {
			t.Errorf("Mask(%q) = %q，want %q", in, got, want)
		}
	}
}

func TestFindInvoices(t *testing.T) {
	text := `電子發票證明聯
AB-12345678  2026/11/19   NT$1,280
隨機碼 9527
CD-87654321 2026-11-20 NT$85
EF-1234567 2026/11/21 NT$300
總計 NT$1,365
GH-00000001 2026/12/01 NT$12,345,678`
	got := FindInvoices(text)
	want := []Invoice{
		{"AB-12345678", "2026-11-19", 1280},
		{"CD-87654321", "2026-11-20", 85},
		{"GH-00000001", "2026-12-01", 12345678},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FindInvoices = %+v\nwant %+v", got, want)
	}
	if got := FindInvoices("今天沒有發票"); got != nil {
		t.Errorf("找不到時 = %+v，want nil", got)
	}
}

func TestSite(t *testing.T) {
	fsys := fstest.MapFS{
		"layout.html":      {Data: []byte("<title>{{.Title}}</title>\n<main>{{block \"content\" .}}預設內容{{end}}</main>\n")},
		"pages/menu.html":  {Data: []byte(`{{define "content"}}{{range .Items}}<p>{{.}}</p>{{end}}{{end}}`)},
		"pages/about.html": {Data: []byte(`{{define "content"}}關於 {{.Shop}}{{end}}`)},
		"pages/empty.html": {Data: []byte(`這個頁面沒有定義 content`)},
		"notes.txt":        {Data: []byte("不是範本")},
	}
	s, err := NewSite(fsys)
	if err != nil {
		t.Fatalf("NewSite：%v", err)
	}

	type menu struct {
		Title string
		Items []string
	}
	type about struct{ Title, Shop string }
	tests := []struct {
		page string
		data any
		want string
	}{
		{"menu", menu{"菜單", []string{"排骨飯", "<滷肉飯>"}}, "<title>菜單</title>\n<main><p>排骨飯</p><p>&lt;滷肉飯&gt;</p></main>\n"},
		{"about", about{"關於我們", "阿嬤 & 孫"}, "<title>關於我們</title>\n<main>關於 阿嬤 &amp; 孫</main>\n"},
		{"empty", about{"空白", ""}, "<title>空白</title>\n<main>預設內容</main>\n"},
	}
	for _, tt := range tests {
		var b strings.Builder
		if err := s.Render(&b, tt.page, tt.data); err != nil || b.String() != tt.want {
			t.Errorf("Render(%q)（%v）寫出 %q\nwant %q", tt.page, err, b.String(), tt.want)
		}
	}

	var b strings.Builder
	if err := s.Render(&b, "contact", nil); err == nil {
		t.Error("沒有 contact 這個頁面，應該回傳錯誤")
	}
	if err := s.Render(&b, "menu", about{"菜單", "阿嬤"}); err == nil || b.Len() != 0 {
		t.Errorf("about 沒有 Items 欄位，應該回傳錯誤而且什麼都不寫；err = %v，寫出 %q", err, b.String())
	}

	fsys["pages/broken.html"] = &fstest.MapFile{Data: []byte(`{{define "content"}}{{if}}{{end}}`)}
	if _, err := NewSite(fsys); err == nil {
		t.Error("pages/broken.html 有語法錯誤，NewSite 應該回傳錯誤")
	}
}

func TestValidPhone(t *testing.T) {
	for s, want := range map[string]bool{
		"0912345678":     true,
		"0912-345-678":   true,
		"09123456789":    false,
		"電話 0912345678":  false,
		"0912-345-6789":  false,
		"0812345678":     false,
		"0912345678 是我的": false,
		"":               false,
	} {
		if got := ValidPhone(s); got != want {
			t.Errorf("ValidPhone(%q) = %v，want %v", s, got, want)
		}
	}
}

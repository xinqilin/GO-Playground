//go:build !solution

package memo

import (
	"io"
	"io/fs"
	"regexp"
)

// ParseKV 解析用分號隔開的「key=value」設定，例如 "name=阿嬤便當; phone=02-2345-6789"。
// key 和 value 前後的空白要去掉；只在第一個 = 切開，所以 value 裡可以有 =；
// 空的段落（例如結尾多一個分號）略過。
// 段落裡沒有 =、key 是空的、同一個 key 出現兩次時，回傳錯誤。
func ParseKV(s string) (map[string]string, error) {
	// TODO: 實作
	return nil, nil
}

// RenderReviews 把評論寫成 HTML 清單，格式如下（每一行後面都有換行）：
//
//	<ul>
//	<li><b>名字</b> 5 顆星：評論</li>
//	</ul>
//
// 名字和評論是使用者輸入的，一定要正確地跳脫。沒有任何評論時，只寫一行 <p>還沒有評論</p>。
func RenderReviews(w io.Writer, reviews []Review) error {
	// TODO: 實作
	return nil
}

// Mask 把文字裡的手機號碼和 email 遮起來，給客服紀錄用，其他文字不變：
//   - 手機號碼：09 開頭共 10 位數字，數字之間可以有 - 或空白，例如 0912-345-678、0912 345 678。
//     不管原本怎麼寫，一律遮成 0912-***-678 的樣子（保留前 4 碼和後 3 碼）。
//     前後還連著其他數字或英文字母時（例如訂單編號 A091234567890）不算手機號碼。
//   - email：帳號只留第一個字元，其他換成 ***，例如 bill@example.com 變成 b***@example.com。
func Mask(s string) string {
	// TODO: 實作
	return ""
}

// FindInvoices 從一段文字（例如掃描發票之後辨識出來的文字）找出所有的發票。
// 每張發票在同一行，格式是「號碼 日期 金額」，中間隔著一個以上的空白：
//   - 號碼：兩個大寫英文字母、-、八位數字
//   - 日期：2026/11/19 或 2026-11-19
//   - 金額：NT$ 加上數字，數字裡可以有逗號，例如 NT$1,280
//
// 回傳的 Date 一律換成 2006-01-02 的格式。一張都找不到時回傳 nil。
func FindInvoices(text string) []Invoice {
	// TODO: 實作
	return nil
}

// Site 是一組網頁範本：一個共用的版型加上好幾個頁面。
type Site struct {
	// TODO: 加上需要的欄位
}

// NewSite 從 fsys 讀範本：layout.html 是版型，裡面用 {{block "content" .}} 留一個位置；
// pages 目錄下的每個 .html 檔是一個頁面，用 {{define "content"}} 填那個位置。
// 頁面的名稱是檔名去掉 .html，例如 pages/menu.html 叫做 "menu"。
// 每個頁面都用自己的 content，不能互相蓋掉。任何一個範本有語法錯誤時回傳錯誤。
// 名稱、評論這類資料都可能來自使用者，所以要用會自動跳脫的範本套件。
func NewSite(fsys fs.FS) (*Site, error) {
	// TODO: 實作
	return &Site{}, nil
}

// Render 用 name 這個頁面加上版型，把 data 寫到 w。沒有這個頁面時回傳錯誤。
// 執行範本失敗時，不能寫任何東西到 w。
func (s *Site) Render(w io.Writer, name string, data any) error {
	// TODO: 實作
	return nil
}

// phonePattern 是台灣的手機號碼：09 開頭共 10 位數字，數字之間可以有 -。
var phonePattern = regexp.MustCompile(`09\d{2}-?\d{3}-?\d{3}`)

// ValidPhone 回傳 s 整個是不是一支手機號碼，例如 0912345678 或 0912-345-678。
//
// 找錯：「09123456789」（11 位數）和「電話 0912345678」都被當成合法的。
func ValidPhone(s string) bool {
	return phonePattern.MatchString(s)
}

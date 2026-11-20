package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

// ReceiptItem 是收據上的一個品項。
type ReceiptItem struct {
	Name  string
	Qty   int
	Price int
}

// Subtotal 在範本裡寫成 {{.Subtotal}}：沒有參數、回傳一個值的方法可以直接呼叫。
func (it ReceiptItem) Subtotal() int { return it.Qty * it.Price }

// Receipt 是一張收據。
type Receipt struct {
	Shop  string
	Table int
	Items []ReceiptItem
	Paid  bool
	Note  string
}

// Total 是所有品項的小計加起來。
func (r Receipt) Total() int {
	sum := 0
	for _, it := range r.Items {
		sum += it.Subtotal()
	}
	return sum
}

// {{- 和 -}} 會吃掉那一側的空白和換行，讓範本可以分行寫，輸出又不會多出空行。
const receiptText = `{{.Shop}} 第 {{.Table}} 桌
{{range .Items -}}
{{.Name}} x{{.Qty}} {{.Subtotal}} 元
{{else -}}
（沒有品項）
{{end -}}
合計 {{.Total}} 元{{if .Paid}}（已付款）{{end}}
{{with .Note}}備註：{{.}}
{{end -}}
`

var receiptTmpl = template.Must(template.New("receipt").Parse(receiptText))

func receiptDemo() {
	receipts := []Receipt{
		{Shop: "阿嬤便當", Table: 7, Paid: true, Note: "統編 12345678",
			Items: []ReceiptItem{{"排骨飯", 2, 90}, {"滷蛋", 3, 15}}},
		{Shop: "阿嬤便當", Table: 3},
	}
	for _, r := range receipts {
		if err := receiptTmpl.Execute(os.Stdout, r); err != nil {
			fmt.Println(err)
		}
	}
}

// execErrorDemo 欄位名稱寫錯不會在 Parse 時發現，要到 Execute 才出錯；
// map 少了 key 預設印 <no value>，加上 missingkey=error 才會變成錯誤。
// Execute 出錯時，錯誤之前的內容已經寫出去了，所以先寫到緩衝區，成功了再送出。
func execErrorDemo() {
	t := template.Must(template.New("typo").Parse("{{.Shp}}\n"))
	fmt.Println(t.Execute(os.Stdout, Receipt{Shop: "阿嬤便當"}))

	data := map[string]string{"shop": "阿嬤便當"}
	loose := template.Must(template.New("loose").Parse("{{.shop}} {{.phone}}\n"))
	if err := loose.Execute(os.Stdout, data); err != nil {
		fmt.Println(err)
	}
	strict := template.Must(template.New("strict").Option("missingkey=error").Parse("{{.shop}} {{.phone}}\n"))
	var b strings.Builder
	err := strict.Execute(&b, data)
	fmt.Printf("寫出了 %q\n%v\n", b.String(), err)
}

package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// InvoiceLine 是發票明細的一行。
type InvoiceLine struct {
	Number string // 發票號碼
	Item   string
	Amount int
}

// writeInvoices 寫出 CSV。欄位裡有逗號、引號、換行時，csv.Writer 會自動加上引號。
func writeInvoices(w io.Writer, lines []InvoiceLine) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"發票號碼", "品項", "金額"}); err != nil {
		return err
	}
	for _, l := range lines {
		if err := cw.Write([]string{l.Number, l.Item, strconv.Itoa(l.Amount)}); err != nil {
			return err
		}
	}
	cw.Flush() // csv.Writer 有緩衝區，最後一定要 Flush，再用 Error 檢查寫入有沒有出錯
	return cw.Error()
}

// readInvoices 讀 CSV：第一行是標題，之後每一行都要剛好 3 個欄位。
func readInvoices(r io.Reader) ([]InvoiceLine, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = 3
	if _, err := cr.Read(); err != nil { // 標題
		return nil, err
	}
	var lines []InvoiceLine
	for {
		rec, err := cr.Read()
		if errors.Is(err, io.EOF) {
			return lines, nil
		}
		if err != nil {
			return nil, err
		}
		amount, err := strconv.Atoi(rec[2])
		if err != nil {
			line, _ := cr.FieldPos(2)
			return nil, fmt.Errorf("第 %d 行：金額 %q 不是整數", line, rec[2])
		}
		lines = append(lines, InvoiceLine{Number: rec[0], Item: rec[1], Amount: amount})
	}
}

func csvDemo() {
	var buf bytes.Buffer
	err := writeInvoices(&buf, []InvoiceLine{
		{"AB-12345678", "排骨飯", 90},
		{"AB-12345678", "珍珠奶茶, 微糖去冰", 55},
		{"AB-12345679", `"招牌"滷肉飯`, 35},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(buf.String())

	lines, err := readInvoices(&buf)
	fmt.Printf("%d 行 %q %v\n", len(lines), lines[1].Item, err)

	_, err = readInvoices(strings.NewReader("發票號碼,品項,金額\nAB-1,紅茶,30\nAB-2,綠茶,二十\n"))
	fmt.Println(err)
	_, err = readInvoices(strings.NewReader("發票號碼,品項,金額\nAB-1,紅茶\n"))
	fmt.Println(err)
}

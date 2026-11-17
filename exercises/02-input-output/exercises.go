//go:build !solution

package inout

import (
	"bufio"
	"fmt"
	"io"
)

// Percent 回傳 part 佔 total 的百分比，取到小數第一位並加上 %，例如 Percent(1, 3) 是 "33.3%"。
// total 小於 1 時回傳 "0.0%"。
func Percent(part, total int) string {
	// TODO: 實作
	return ""
}

// HexColor 把紅、綠、藍三個 0–255 的值寫成網頁用的色碼，例如 HexColor(255, 136, 0) 是 "#FF8800"。
// 每個顏色固定兩位，英文字母用大寫。
func HexColor(r, g, b uint8) string {
	// TODO: 實作
	return ""
}

// String 讓 Money 印成有千分位的台幣，例如 1234567 印成 "NT$1,234,567"，-1500 印成 "-NT$1,500"。
func (m Money) String() string {
	// TODO: 實作
	return ""
}

// WriteReceipt 把明細寫到 w。每個品項一行，最後是分隔線和合計：
//
//	  170  排骨飯 x2
//	   45  滷蛋 x3
//	-----
//	  215  合計
//
// 金額（數量乘單價）靠右對齊、佔 5 格，後面空兩格接品名、空一格接 x 和數量。
// 寫入失敗時，立刻回傳那個錯誤。
func WriteReceipt(w io.Writer, items []Item) error {
	// TODO: 實作
	return nil
}

// ParseOrders 從 r 讀訂單，回傳總金額。每一行是「品名 數量 單價」，用空白分開，例如 "排骨飯 2 85"。
// 空白行和 # 開頭的註解行要略過。遇到格式不對的行（欄位不是 3 個，或數量、單價不是整數），
// 回傳 0 和錯誤，錯誤訊息是 fmt.Errorf("第 %d 行格式錯誤：%q", 行號, 那一行的原始內容)，行號從 1 開始。
// 讀取本身失敗時，回傳 0 和那個錯誤。
func ParseOrders(r io.Reader) (int, error) {
	// TODO: 實作
	return 0, nil
}

// CallNumbers 把叫號清單寫到 w：每個名字一行「請 <name> 到櫃台」，最後一行「== 以上 ==」。
//
// 找錯：測試發現「== 以上 ==」跑到了最前面。
func CallNumbers(w io.Writer, names []string) error {
	bw := bufio.NewWriter(w)
	for _, name := range names {
		fmt.Fprintf(bw, "請 %s 到櫃台\n", name)
	}
	fmt.Fprintln(w, "== 以上 ==")
	return bw.Flush()
}

// LongestLine 回傳 r 裡最長一行的長度（以 byte 計，不含換行）。一行最長可能到 1 MB。
// 讀取失敗時回傳 0 和那個錯誤。
func LongestLine(r io.Reader) (int, error) {
	// TODO: 實作
	return 0, nil
}

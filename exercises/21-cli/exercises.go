//go:build !solution

package kiosk

import (
	"context"
	"flag"
	"fmt"
	"io"
)

// ParseBooking 用 flag.FlagSet 解析訂位的參數（欄位說明見 Booking）。
// 旗標錯誤、少了 -name、人數不在 1 到 12、時間格式不對時回傳錯誤；
// -h 回傳的錯誤要讓 errors.Is(err, flag.ErrHelp) 成立。
// 這個函式不能印任何東西到 os.Stderr。
func ParseBooking(args []string) (Booking, error) {
	// TODO: 實作
	return Booking{}, nil
}

// LoadConfig 決定設定值（欄位說明見 Config），優先順序是：命令列旗標 > 環境變數 > 預設值。
// 環境變數的值看不懂時，回傳的錯誤訊息要提到那個環境變數的名稱。
// 這個函式不能印任何東西到 os.Stderr。
func LoadConfig(args []string, getenv func(string) string) (Config, error) {
	// TODO: 實作
	return Config{}, nil
}

// Relay 把 r 一行一行抄到 w（每行後面加上 \n），回傳抄了幾行。
// 每抄一行之前先看 ctx：已經取消就停下來，回傳 context.Cause(ctx)。
func Relay(ctx context.Context, r io.Reader, w io.Writer) (int, error) {
	// TODO: 實作
	return 0, nil
}

// Count 是一個小型的 wc：count [-w] [檔案…]。
// 每個輸入印一行「數量 名稱」；沒有 -w 時數行數，有 -w 時數字數（以空白分隔）。
// 沒有給檔案、或檔名是 - 時讀 stdin，名稱印成 -。輸入超過一個時，最後再印「數量 總計」。
// 打不開的檔案：在 stderr 印「count：錯誤」，繼續處理下一個，最後回傳 1。
// 旗標錯誤回傳 2（訊息寫到 stderr），-h 回傳 0，全部成功回傳 0。
func Count(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	// TODO: 實作
	return 0
}

// Report 是 report [-n 個數] 品項…：最多印出前 n 個品項，每行一個。
// 成功回傳 0；-h 把用法寫到 stderr，回傳 0；
// 旗標錯誤、或一個品項都沒有給時，錯誤訊息寫到 stderr，回傳 2。stdout 只放結果。
//
// 找錯：用法錯誤時，結束碼和訊息寫去的地方都不對。
func Report(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	n := fs.Int("n", 3, "最多列出幾個")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stdout, err)
		return 1
	}
	if fs.NArg() == 0 {
		fmt.Fprintln(stdout, "report：至少要給一個品項")
		return 0
	}
	for i, item := range fs.Args() {
		if i == *n {
			break
		}
		fmt.Fprintln(stdout, item)
	}
	return 0
}

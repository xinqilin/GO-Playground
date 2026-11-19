package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
)

// 結束碼。
const (
	exitOK          = 0   // 成功
	exitError       = 1   // 執行時出錯，例如檔案打不開
	exitUsage       = 2   // 用法錯誤，例如旗標的值不對
	exitInterrupted = 130 // 被 Ctrl-C 中斷：128 加上 SIGINT 的編號 2
)

const usage = `用法：tally <子命令> [旗標] [檔案…]

子命令：
  total   加總營業額
  top     列出賣最多的品項

沒有給檔案、或檔名是 - 時，從標準輸入讀。
每一行是「品項 金額」，例如：排骨飯 90
`

// sys 是程式和作業系統打交道的四個地方。main 接上真的，測試接上假的。
type sys struct {
	stdin          io.Reader
	stdout, stderr io.Writer
	getenv         func(string) string
}

// run 是整個程式：看第一個參數決定子命令，回傳結束碼。
func run(ctx context.Context, args []string, s sys) int {
	if len(args) == 0 {
		fmt.Fprint(s.stderr, usage)
		return exitUsage
	}
	switch cmd, rest := args[0], args[1:]; cmd {
	case "total":
		return total(ctx, rest, s)
	case "top":
		return top(ctx, rest, s)
	case "help", "-h", "-help", "--help":
		fmt.Fprint(s.stdout, usage)
		return exitOK
	default:
		fmt.Fprintf(s.stderr, "tally：不認得的子命令 %q\n\n%s", cmd, usage)
		return exitUsage
	}
}

// parse 解析旗標。第二個回傳值是 false 時，呼叫的人直接用第一個回傳值當結束碼。
// 錯誤訊息和用法說明，flag 套件已經寫到 fs.Output() 了，這裡不用再印。
func parse(fs *flag.FlagSet, args []string) (int, bool) {
	err := fs.Parse(args)
	switch {
	case err == nil:
		return exitOK, true
	case errors.Is(err, flag.ErrHelp):
		return exitOK, false // 使用者自己要看說明，不算失敗
	default:
		return exitUsage, false
	}
}

// fail 把錯誤寫到 stderr，回傳對應的結束碼。
func fail(ctx context.Context, stderr io.Writer, prog string, err error) int {
	fmt.Fprintf(stderr, "%s：%v\n", prog, err)
	if ctx.Err() != nil {
		return exitInterrupted
	}
	return exitError
}

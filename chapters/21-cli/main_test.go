package main

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
)

// sales 是 Example 共用的標準輸入：四筆賣出紀錄。
const sales = "排骨飯 90\n雞腿飯 100\n排骨飯 90\n滷肉飯 60\n"

// try 用假的標準輸入和環境變數執行 tally，像終端機一樣印出指令、stdout、stderr 和結束碼。
// stderr 的每一行前面加上「stderr｜」，和 stdout 分得出來。
func try(ctx context.Context, stdin string, env map[string]string, args ...string) {
	var stdout, stderr strings.Builder
	code := run(ctx, args, sys{
		stdin:  strings.NewReader(stdin),
		stdout: &stdout,
		stderr: &stderr,
		getenv: func(key string) string { return env[key] },
	})

	prompt := []string{"$"}
	for _, k := range slices.Sorted(maps.Keys(env)) {
		prompt = append(prompt, k+"="+env[k])
	}
	fmt.Println(strings.Join(slices.Concat(prompt, []string{"tally"}, args), " "))
	fmt.Print(stdout.String())
	for line := range strings.Lines(stderr.String()) {
		fmt.Println(strings.TrimRight("stderr｜"+line, " \n"))
	}
	fmt.Println("結束碼", code)
}

func Example_total() {
	try(context.Background(), sales, nil, "total")
	try(context.Background(), sales, nil, "total", "-tax=10")
	// Output:
	// $ tally total
	// 合計 340 元，含稅 5% 共 357 元
	// 結束碼 0
	// $ tally total -tax=10
	// 合計 340 元，含稅 10% 共 374 元
	// 結束碼 0
}

func Example_badFlag() {
	try(context.Background(), sales, nil, "total", "-tax", "abc")
	// Output:
	// $ tally total -tax abc
	// stderr｜invalid value "abc" for flag -tax: parse error
	// stderr｜用法：tally total [-tax 百分比] [檔案…]
	// stderr｜  -tax int
	// stderr｜    	稅率（百分比），也可以用環境變數 TALLY_TAX 設定 (default 5)
	// 結束碼 2
}

func Example_help() {
	try(context.Background(), "", nil, "top", "-h")
	// Output:
	// $ tally top -h
	// stderr｜用法：tally top [-n 個數] [檔案…]
	// stderr｜  -n int
	// stderr｜    	列出幾個品項 (default 3)
	// 結束碼 0
}

func Example_flagOrder() {
	try(context.Background(), sales, nil, "total", "-", "-tax", "10")
	// Output:
	// $ tally total - -tax 10
	// stderr｜tally total：open -tax: no such file or directory
	// 結束碼 1
}

func Example_subcommands() {
	try(context.Background(), sales, nil, "top", "-n", "2")
	try(context.Background(), "", nil, "sum")
	// Output:
	// $ tally top -n 2
	// 1. 排骨飯 2 份
	// 2. 滷肉飯 1 份
	// 結束碼 0
	// $ tally sum
	// stderr｜tally：不認得的子命令 "sum"
	// stderr｜
	// stderr｜用法：tally <子命令> [旗標] [檔案…]
	// stderr｜
	// stderr｜子命令：
	// stderr｜  total   加總營業額
	// stderr｜  top     列出賣最多的品項
	// stderr｜
	// stderr｜沒有給檔案、或檔名是 - 時，從標準輸入讀。
	// stderr｜每一行是「品項 金額」，例如：排骨飯 90
	// 結束碼 2
}

func Example_files() {
	try(context.Background(), "", nil, "total", "testdata/mon.txt", "testdata/tue.txt")
	try(context.Background(), "滷肉飯 60\n", nil, "top", "testdata/mon.txt", "-")
	try(context.Background(), "", nil, "total", "testdata/bad.txt")
	try(context.Background(), "", nil, "total", "testdata/nope.txt")
	// Output:
	// $ tally total testdata/mon.txt testdata/tue.txt
	// 合計 340 元，含稅 5% 共 357 元
	// 結束碼 0
	// $ tally top testdata/mon.txt -
	// 1. 排骨飯 1 份
	// 2. 滷肉飯 1 份
	// 3. 雞腿飯 1 份
	// 結束碼 0
	// $ tally total testdata/bad.txt
	// stderr｜tally total：testdata/bad.txt 第 2 行：看不懂「雞腿飯 一百」
	// 結束碼 1
	// $ tally total testdata/nope.txt
	// stderr｜tally total：open testdata/nope.txt: no such file or directory
	// 結束碼 1
}

func Example_config() {
	try(context.Background(), sales, map[string]string{"TALLY_TAX": "10"}, "total")
	try(context.Background(), sales, map[string]string{"TALLY_TAX": "10"}, "total", "-tax", "0")
	try(context.Background(), sales, map[string]string{"TALLY_TAX": "十"}, "total")
	// Output:
	// $ TALLY_TAX=10 tally total
	// 合計 340 元，含稅 10% 共 374 元
	// 結束碼 0
	// $ TALLY_TAX=10 tally total -tax 0
	// 合計 340 元，含稅 0% 共 340 元
	// 結束碼 0
	// $ TALLY_TAX=十 tally total
	// stderr｜tally total：環境變數 TALLY_TAX="十" 不是整數
	// 結束碼 2
}

func Example_interrupted() {
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(errors.New("interrupt signal received")) // 假裝使用者已經按了 Ctrl-C
	try(ctx, sales, nil, "total")
	// Output:
	// $ tally total
	// stderr｜tally total：interrupt signal received
	// 結束碼 130
}

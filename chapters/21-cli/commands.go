package main

import (
	"cmp"
	"context"
	"flag"
	"fmt"
	"maps"
	"slices"
	"strings"
)

// total 是 tally total：加總所有輸入的金額，再算出含稅的總額。
func total(ctx context.Context, args []string, s sys) int {
	fs := flag.NewFlagSet("tally total", flag.ContinueOnError)
	fs.SetOutput(s.stderr)

	tax, err := envInt(s.getenv, "TALLY_TAX", 5)
	if err != nil {
		fmt.Fprintf(s.stderr, "tally total：%v\n", err)
		return exitUsage
	}
	fs.IntVar(&tax, "tax", tax, "稅率（百分比），也可以用環境變數 TALLY_TAX 設定")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "用法：tally total [-tax 百分比] [檔案…]")
		fs.PrintDefaults()
	}
	if code, ok := parse(fs, args); !ok {
		return code
	}

	sum := 0
	err = readSales(ctx, fs.Args(), s.stdin, func(_ string, price int) {
		sum += price
	})
	if err != nil {
		return fail(ctx, s.stderr, "tally total", err)
	}

	withTax := (sum*(100+tax) + 50) / 100 // 四捨五入到整數元
	fmt.Fprintf(s.stdout, "合計 %d 元，含稅 %d%% 共 %d 元\n", sum, tax, withTax)
	return exitOK
}

// top 是 tally top：依照賣出的份數排出前 n 名，份數一樣時照品項名稱排。
func top(ctx context.Context, args []string, s sys) int {
	fs := flag.NewFlagSet("tally top", flag.ContinueOnError)
	fs.SetOutput(s.stderr)
	n := fs.Int("n", 3, "列出幾個品項")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "用法：tally top [-n 個數] [檔案…]")
		fs.PrintDefaults()
	}
	if code, ok := parse(fs, args); !ok {
		return code
	}
	if *n < 1 {
		fmt.Fprintln(s.stderr, "tally top：-n 至少要是 1")
		return exitUsage
	}

	counts := map[string]int{}
	err := readSales(ctx, fs.Args(), s.stdin, func(item string, _ int) {
		counts[item]++
	})
	if err != nil {
		return fail(ctx, s.stderr, "tally top", err)
	}

	items := slices.Collect(maps.Keys(counts))
	slices.SortFunc(items, func(a, b string) int {
		return cmp.Or(cmp.Compare(counts[b], counts[a]), strings.Compare(a, b))
	})
	for i, item := range items[:min(*n, len(items))] {
		fmt.Fprintf(s.stdout, "%d. %s %d 份\n", i+1, item, counts[item])
	}
	return exitOK
}

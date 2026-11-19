package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// readSales 依序讀每個輸入，每一筆交給 add。
// 沒有給檔名、或檔名是 - 時讀標準輸入，這是命令列工具共同的慣例。
func readSales(ctx context.Context, names []string, stdin io.Reader, add func(item string, price int)) error {
	if len(names) == 0 {
		names = []string{"-"}
	}
	for _, name := range names {
		if err := readSalesFrom(ctx, name, stdin, add); err != nil {
			return err
		}
	}
	return nil
}

// readSalesFrom 讀一個輸入。每讀一行先看 ctx：使用者按了 Ctrl-C，就回傳取消的原因。
func readSalesFrom(ctx context.Context, name string, stdin io.Reader, add func(string, int)) error {
	r, label := stdin, "標準輸入"
	if name != "-" {
		f, err := os.Open(name)
		if err != nil {
			return err
		}
		defer f.Close()
		r, label = f, name
	}

	sc := bufio.NewScanner(r)
	for line := 1; sc.Scan(); line++ {
		if ctx.Err() != nil {
			return context.Cause(ctx)
		}
		text := strings.TrimSpace(sc.Text())
		if text == "" {
			continue
		}
		i := strings.LastIndexByte(text, ' ')
		price, err := strconv.Atoi(text[i+1:])
		if i < 0 || err != nil {
			return fmt.Errorf("%s 第 %d 行：看不懂「%s」", label, line, text)
		}
		add(text[:i], price)
	}
	return sc.Err()
}

// envInt 讀整數型的環境變數，沒有設定時回傳 def。
// 結果拿來當旗標的預設值，就得到「命令列 > 環境變數 > 內建預設值」的優先順序。
func envInt(getenv func(string) string, key string, def int) (int, error) {
	v := getenv(key)
	if v == "" {
		return def, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("環境變數 %s=%q 不是整數", key, v)
	}
	return n, nil
}

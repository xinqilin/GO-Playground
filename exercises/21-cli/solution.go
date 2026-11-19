//go:build solution

package kiosk

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// ParseBooking 用 io.Discard 當 FlagSet 的輸出，flag 套件就不會自己印錯誤訊息。
func ParseBooking(args []string) (Booking, error) {
	var b Booking
	fs := flag.NewFlagSet("booking", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&b.Name, "name", "", "訂位的人")
	fs.IntVar(&b.People, "people", 2, "人數")
	fs.StringVar(&b.Time, "time", "18:00", "時間")
	fs.BoolVar(&b.Outdoor, "outdoor", false, "坐戶外")
	if err := fs.Parse(args); err != nil {
		return Booking{}, err
	}

	switch {
	case b.Name == "":
		return Booking{}, errors.New("一定要給 -name")
	case b.People < 1 || b.People > 12:
		return Booking{}, fmt.Errorf("人數 %d 不在 1 到 12 之間", b.People)
	}
	if _, err := time.Parse("15:04", b.Time); err != nil {
		return Booking{}, fmt.Errorf("時間 %q 的格式不對：%w", b.Time, err)
	}
	b.Note = strings.Join(fs.Args(), " ")
	return b, nil
}

// LoadConfig 先用環境變數算出每個旗標的預設值，再解析旗標：命令列有寫就蓋過去。
func LoadConfig(args []string, getenv func(string) string) (Config, error) {
	c := Config{Port: 8080, Shop: "總店"}
	if v := getenv("KIOSK_PORT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("KIOSK_PORT=%q 不是整數", v)
		}
		c.Port = n
	}
	if v := getenv("KIOSK_SHOP"); v != "" {
		c.Shop = v
	}
	if v := getenv("KIOSK_DEBUG"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return Config{}, fmt.Errorf("KIOSK_DEBUG=%q 看不懂", v)
		}
		c.Debug = b
	}

	fs := flag.NewFlagSet("kiosk", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.IntVar(&c.Port, "port", c.Port, "埠號")
	fs.StringVar(&c.Shop, "shop", c.Shop, "分店名稱")
	fs.BoolVar(&c.Debug, "debug", c.Debug, "除錯模式")
	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}
	return c, nil
}

// Relay 在讀下一行之前檢查 ctx，取消之後一行都不再寫。
func Relay(ctx context.Context, r io.Reader, w io.Writer) (int, error) {
	n := 0
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		if ctx.Err() != nil {
			return n, context.Cause(ctx)
		}
		if _, err := fmt.Fprintln(w, sc.Text()); err != nil {
			return n, err
		}
		n++
	}
	return n, sc.Err()
}

// Count 把「數一個輸入」抽成 countOne，打不開的檔案記下來，繼續處理下一個。
func Count(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("count", flag.ContinueOnError)
	fs.SetOutput(stderr)
	words := fs.Bool("w", false, "數字數，不數行數")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}

	names := fs.Args()
	if len(names) == 0 {
		names = []string{"-"}
	}
	code, sum := 0, 0
	for _, name := range names {
		n, err := countOne(name, stdin, *words)
		if err != nil {
			fmt.Fprintf(stderr, "count：%v\n", err)
			code = 1
			continue
		}
		sum += n
		fmt.Fprintln(stdout, n, name)
	}
	if len(names) > 1 {
		fmt.Fprintln(stdout, sum, "總計")
	}
	return code
}

// countOne 數一個輸入；在自己的函式裡 defer Close，每個檔案用完馬上關。
func countOne(name string, stdin io.Reader, words bool) (int, error) {
	r := stdin
	if name != "-" {
		f, err := os.Open(name)
		if err != nil {
			return 0, err
		}
		defer f.Close()
		r = f
	}
	sc := bufio.NewScanner(r)
	if words {
		sc.Split(bufio.ScanWords)
	}
	n := 0
	for sc.Scan() {
		n++
	}
	return n, sc.Err()
}

// Report：FlagSet 的輸出接到 stderr，錯誤訊息也寫到 stderr；-h 回傳 0，其他用法錯誤回傳 2。
func Report(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	fs.SetOutput(stderr)
	n := fs.Int("n", 3, "最多列出幾個")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if fs.NArg() == 0 {
		fmt.Fprintln(stderr, "report：至少要給一個品項")
		return 2
	}
	for i, item := range fs.Args() {
		if i == *n {
			break
		}
		fmt.Fprintln(stdout, item)
	}
	return 0
}

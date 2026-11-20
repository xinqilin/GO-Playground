// 專題 3：並行網址檢查器。
//
// 同時檢查一批網址：固定數量的 worker、每次請求的逾時、每秒請求數的上限、失敗後的退避重試，
// 按 Ctrl-C 會取消還沒完成的檢查。有任何網址失敗時以 exit code 1 結束，方便接在 CI 或排程裡。
//
// 執行：go run ./projects/urlchecker https://go.dev https://pkg.go.dev
//
//	go run ./projects/urlchecker -workers 8 -rate 5 < urls.txt
package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

// run 是整支程式：args 不含程式名稱，沒有給網址時從 stdin 每行讀一個，回傳 exit code。
func run(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer, f Fetcher) int {
	fs := flag.NewFlagSet("urlchecker", flag.ContinueOnError)
	fs.SetOutput(stderr)
	var cfg Config
	fs.IntVar(&cfg.Workers, "workers", 4, "同時最多幾個請求")
	fs.DurationVar(&cfg.Timeout, "timeout", 5*time.Second, "每一次請求最多等多久")
	fs.IntVar(&cfg.PerSecond, "rate", 0, "每秒最多送出幾個請求，0 是不限制")
	fs.IntVar(&cfg.Retries, "retries", 2, "失敗後最多再試幾次")
	fs.DurationVar(&cfg.Backoff, "backoff", 200*time.Millisecond, "第一次重試前等多久，之後每次加倍")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if cfg.Workers < 1 || cfg.Timeout <= 0 || cfg.PerSecond < 0 || cfg.Retries < 0 || cfg.Backoff < 0 {
		fmt.Fprintln(stderr, "urlchecker：-workers 至少是 1，-timeout 要大於 0，其他數值不能是負數")
		return 2
	}

	urls := fs.Args()
	if len(urls) == 0 {
		var err error
		if urls, err = readURLs(stdin); err != nil {
			fmt.Fprintf(stderr, "urlchecker：讀取網址：%v\n", err)
			return 2
		}
	}
	if len(urls) == 0 {
		fmt.Fprintln(stderr, "用法：urlchecker [旗標] URL...（或從標準輸入每行讀一個網址）")
		return 2
	}

	start := time.Now()
	results := Check(ctx, urls, cfg, f)
	ok := 0
	for _, r := range results {
		fmt.Fprintln(stdout, formatResult(r))
		if r.OK() {
			ok++
		}
	}
	fmt.Fprintf(stdout, "共 %d 個網址：%d 個正常，%d 個失敗，花了 %v\n",
		len(results), ok, len(results)-ok, time.Since(start).Round(time.Millisecond))
	if ok < len(results) {
		return 1
	}
	return 0
}

// readURLs 每行讀一個網址，略過空白行和 # 開頭的註解。
func readURLs(r io.Reader) ([]string, error) {
	var urls []string
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		s := strings.TrimSpace(sc.Text())
		if s != "" && !strings.HasPrefix(s, "#") {
			urls = append(urls, s)
		}
	}
	return urls, sc.Err()
}

// formatResult 把一個結果排成一行：狀態、狀態碼、花的時間、網址，再加上失敗原因和重試次數。
func formatResult(r Result) string {
	mark, code := "OK", strconv.Itoa(r.Status)
	if !r.OK() {
		mark = "FAIL"
	}
	if r.Err != nil {
		code = "---"
	}
	s := fmt.Sprintf("%-4s %3s %6v  %s", mark, code, r.Took.Round(time.Millisecond), r.URL)

	var notes []string
	if r.Err != nil {
		notes = append(notes, describe(r.Err))
	}
	if r.Attempts > 1 {
		notes = append(notes, fmt.Sprintf("試了 %d 次", r.Attempts))
	}
	if len(notes) > 0 {
		s += "（" + strings.Join(notes, "，") + "）"
	}
	return s
}

// describe 把常見的錯誤換成短的說明。
func describe(err error) string {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return "逾時"
	case errors.Is(err, context.Canceled):
		return "已取消"
	}
	return err.Error()
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	context.AfterFunc(ctx, stop) // 第一次 Ctrl-C 之後恢復預設行為：再按一次就直接結束

	f := httpFetcher{client: &http.Client{}} // 逾時由 context 控制
	code := run(ctx, os.Args[1:], os.Stdin, os.Stdout, os.Stderr, f)
	stop()
	os.Exit(code)
}

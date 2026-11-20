// 專題 5：日誌分析。
//
// 讀 access log 或 slog 的 JSON log（可以是 .gz），算出請求數、狀態碼分布、最忙的一分鐘，
// 以及每一類請求的 p50、p95 和最慢的處理時間。輸出可以是文字、JSON 或 CSV；
// 5xx 的比例超過門檻時以 exit code 1 結束，可以接在排程裡發警報。
//
// 執行：go run ./projects/logstats projects/logstats/testdata/access.log
//
//	go run ./projects/logstats -format json -fail-5xx 5 app.log.gz
package main

import (
	"bufio"
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// maxWarnings 是最多在 stderr 列出幾行無法解析的內容，其他的只算數量。
const maxWarnings = 3

// run 是整支程式，回傳 exit code：0 正常，1 超過 5xx 門檻，2 用法錯誤或讀不到檔案。
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("logstats", flag.ContinueOnError)
	fs.SetOutput(stderr)
	format := fs.String("format", "text", "輸出格式：text、json 或 csv")
	top := fs.Int("top", 10, "最多列出幾類請求")
	var since, until time.Time
	fs.TextVar(&since, "since", time.Time{}, "只算這個時間（含）之後的請求，RFC 3339 格式")
	fs.TextVar(&until, "until", time.Time{}, "只算這個時間之前的請求，RFC 3339 格式")
	fail := fs.Float64("fail-5xx", 100, "5xx 的比例超過這個百分比時，以 exit code 1 結束")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	write, ok := writers[*format]
	if !ok {
		fmt.Fprintf(stderr, "logstats：-format 要是 text、json 或 csv，收到 %q\n", *format)
		return 2
	}

	files := fs.Args()
	if len(files) == 0 {
		files = []string{"-"}
	}
	s := newStats()
	for _, name := range files {
		err := readFile(name, stdin, func(n int, line string) {
			e, err := ParseLine(line)
			if errors.Is(err, errSkip) {
				return
			}
			if err != nil {
				s.Bad++
				if s.Bad <= maxWarnings {
					fmt.Fprintf(stderr, "%s 第 %d 行：%v\n", name, n, err)
				}
				return
			}
			if (since.IsZero() || !e.Time.Before(since)) && (until.IsZero() || e.Time.Before(until)) {
				s.Add(e)
			}
		})
		if err != nil {
			fmt.Fprintf(stderr, "logstats：%v\n", err)
			return 2
		}
	}

	r := s.Report(*top)
	if err := write(stdout, r); err != nil {
		fmt.Fprintf(stderr, "logstats：%v\n", err)
		return 2
	}
	if r.Total > 0 && percent(r.ByClass["5xx"], r.Total) > *fail {
		return 1
	}
	return 0
}

// readFile 一行一行讀 name，把行號（從 1 開始）和內容交給 fn。
// name 是 "-" 時讀 stdin；.gz 結尾的檔案一邊讀一邊解壓縮，不用先解開。
func readFile(name string, stdin io.Reader, fn func(n int, line string)) error {
	var r io.Reader = stdin
	if name != "-" {
		f, err := os.Open(name)
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
		if strings.HasSuffix(name, ".gz") {
			zr, err := gzip.NewReader(f)
			if err != nil {
				return fmt.Errorf("%s：%w", name, err)
			}
			r = zr
		}
	}

	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024) // 一行最長 1MB，預設只有 64KB
	n := 0
	for sc.Scan() {
		n++
		fn(n, sc.Text())
	}
	if err := sc.Err(); err != nil {
		return fmt.Errorf("%s 第 %d 行之後：%w", name, n, err)
	}
	return nil
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

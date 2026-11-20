package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// update 是這個 package 的測試自己加的旗標：go test ./projects/logstats -run TestTextReport -update
var update = flag.Bool("update", false, "用目前的輸出覆寫 testdata/access.golden")

// TestTextReport 把文字報表和 golden file 逐字比對。
func TestTextReport(t *testing.T) {
	var out bytes.Buffer
	if code := run([]string{"testdata/access.log"}, nil, &out, io.Discard); code != 0 {
		t.Errorf("exit code = %d，想要 0", code)
	}
	golden := filepath.Join("testdata", "access.golden")
	if *update {
		if err := os.WriteFile(golden, out.Bytes(), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("%v（第一次請加 -update 產生）", err)
	}
	if got := out.String(); got != string(want) {
		t.Errorf("輸出和 %s 不一樣：\n--- got\n%s--- want\n%s", golden, got, want)
	}
}

// Example_warnings 只看 stderr：無法解析的行列出檔名、行號和原因。
func Example_warnings() {
	run([]string{"testdata/access.log"}, nil, io.Discard, os.Stdout)
	// Output:
	// testdata/access.log 第 15 行：格式不對
	// testdata/access.log 第 18 行：格式不對
	// testdata/access.log 第 21 行：時間不對：parsing time "31/Nov/2026:09:33:40 +0800": day out of range
}

func Example_json() {
	run([]string{"-format", "json", "testdata/app.jsonl"}, nil, os.Stdout, os.Stdout)
	// Output:
	// {
	//   "from": "2026-11-21T09:35:10.12+08:00",
	//   "to": "2026-11-21T09:36:00.01+08:00",
	//   "total": 5,
	//   "bad_lines": 0,
	//   "by_class": {
	//     "2xx": 3,
	//     "4xx": 1,
	//     "5xx": 1
	//   },
	//   "routes": [
	//     {
	//       "route": "POST /todos",
	//       "count": 2,
	//       "errors_5xx": 0,
	//       "p50_ms": 0.15,
	//       "p95_ms": 2.35,
	//       "max_ms": 2.35
	//     },
	//     {
	//       "route": "DELETE /todos/{id}",
	//       "count": 1,
	//       "errors_5xx": 1,
	//       "p50_ms": 5.12,
	//       "p95_ms": 5.12,
	//       "max_ms": 5.12
	//     },
	//     {
	//       "route": "GET /todos",
	//       "count": 1,
	//       "errors_5xx": 0,
	//       "p50_ms": 0.41,
	//       "p95_ms": 0.41,
	//       "max_ms": 0.41
	//     },
	//     {
	//       "route": "PATCH /todos/{id}",
	//       "count": 1,
	//       "errors_5xx": 0,
	//       "p50_ms": 1.98,
	//       "p95_ms": 1.98,
	//       "max_ms": 1.98
	//     }
	//   ]
	// }
}

func Example_csv() {
	run([]string{"-format", "csv", "-top", "3", "testdata/access.log"}, nil, os.Stdout, io.Discard)
	// Output:
	// route,count,errors_5xx,p50_ms,p95_ms,max_ms
	// GET /menu,6,0,11,14,14
	// POST /orders,6,2,52,1204,1204
	// GET /orders/{id},3,0,8,9,9
}

func Example_since() {
	run([]string{"-since", "2026-11-21T09:32:00+08:00", "-top", "2", "testdata/access.log"}, nil, os.Stdout, io.Discard)
	// Output:
	// 期間：2026-11-21 09:32:30 +08:00 到 2026-11-21 09:34:30 +08:00
	// 請求：5 筆，另有 3 行無法解析
	// 狀態：2xx 3、3xx 0、4xx 1、5xx 1（20.0%）
	// 最忙的一分鐘：09:32，2 筆
	//
	//   count  5xx   p50    p95    max  route
	//       2    1  48ms  230ms  230ms  POST /orders
	//       1    0   1ms    1ms    1ms  GET /healthz
}

func TestExitCodes(t *testing.T) {
	tests := []struct {
		args     []string
		wantCode int
		wantErr  string
	}{
		{[]string{"testdata/access.log"}, 0, ""},
		{[]string{"-fail-5xx", "20", "testdata/access.log"}, 0, ""}, // 10.5% 沒超過 20%
		{[]string{"-fail-5xx", "5", "testdata/access.log"}, 1, ""},  // 10.5% 超過 5%
		{[]string{"-format", "xml", "testdata/access.log"}, 2, "-format 要是 text、json 或 csv"},
		{[]string{"-since", "昨天", "testdata/access.log"}, 2, `invalid value "昨天" for flag -since`},
		{[]string{"testdata/nosuch.log"}, 2, "nosuch.log"},
	}
	for _, tt := range tests {
		var stderr strings.Builder
		if code := run(tt.args, nil, io.Discard, &stderr); code != tt.wantCode || !strings.Contains(stderr.String(), tt.wantErr) {
			t.Errorf("run(%q) = %d，stderr %q；想要 %d 和 %q", tt.args, code, stderr.String(), tt.wantCode, tt.wantErr)
		}
	}
}

// TestGzipAndStdin：同一份 log 壓縮成 .gz、或從 stdin 讀，結果都要和直接讀檔案一樣。
func TestGzipAndStdin(t *testing.T) {
	data, err := os.ReadFile("testdata/access.log")
	if err != nil {
		t.Fatal(err)
	}
	gz := filepath.Join(t.TempDir(), "access.log.gz")
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	zw.Write(data)
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(gz, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}

	report := func(args []string, stdin io.Reader) string {
		var out strings.Builder
		run(append([]string{"-format", "csv"}, args...), stdin, &out, io.Discard)
		return out.String()
	}
	want := report([]string{"testdata/access.log"}, nil)
	if got := report([]string{gz}, nil); got != want {
		t.Errorf(".gz 的結果不同：\n%s", got)
	}
	if got := report(nil, bytes.NewReader(data)); got != want {
		t.Errorf("stdin 的結果不同：\n%s", got)
	}
}

func TestLongLine(t *testing.T) {
	long := strings.Repeat("x", 2*1024*1024) + "\n"
	var stderr strings.Builder
	if code := run(nil, strings.NewReader(long), io.Discard, &stderr); code != 2 || !strings.Contains(stderr.String(), "token too long") {
		t.Errorf("超過 1MB 的一行：exit code %d，stderr %q", code, stderr.String())
	}
}

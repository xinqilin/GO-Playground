package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/synctest"
)

// update 是這個 package 的測試自己加的旗標：go test ./projects/urlchecker -run TestReport -update
var update = flag.Bool("update", false, "用目前的輸出覆寫 testdata/report.golden")

// TestReport 用假的網路和假的時鐘跑整支程式，輸出（包括每一個時間）和 golden file 逐字比對。
func TestReport(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		args := []string{"-workers", "2", "-timeout", "500ms",
			"https://ok.test", "https://missing.test", "https://flaky.test", "https://slow.test", "ftp://files.test"}
		var out bytes.Buffer
		if code := run(t.Context(), args, nil, &out, &out, newFake()); code != 1 {
			t.Errorf("exit code = %d，想要 1", code)
		}

		golden := filepath.Join("testdata", "report.golden")
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
	})
}

func TestRunExitCodes(t *testing.T) {
	tests := []struct {
		args     []string
		stdin    string
		wantCode int
		wantOut  string
	}{
		{[]string{"https://ok.test", "https://ok.test"}, "", 0, "共 2 個網址：2 個正常，0 個失敗"},
		{nil, "https://ok.test\n# 註解\n", 0, "共 1 個網址"},
		{nil, "https://missing.test\n", 1, "1 個失敗"},
		{nil, "", 2, "用法"},
		{nil, "# 只有註解\n\n", 2, "用法"},
		{[]string{"-workers", "0", "https://ok.test"}, "", 2, "-workers 至少是 1"},
		{[]string{"-timeout", "soon", "https://ok.test"}, "", 2, `invalid value "soon" for flag -timeout`},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			var out strings.Builder
			code := run(t.Context(), tt.args, strings.NewReader(tt.stdin), &out, &out, newFake())
			if code != tt.wantCode || !strings.Contains(out.String(), tt.wantOut) {
				t.Errorf("run(%q) = %d，輸出\n%s想要 %d 和 %q", tt.args, code, out.String(), tt.wantCode, tt.wantOut)
			}
		})
	}
}

func Example_readURLs() {
	in := strings.NewReader("# 便當店的網站\nhttps://shop.test\n\n  https://shop.test/menu\n# https://old.shop.test 先不檢查\n")
	urls, err := readURLs(in)
	fmt.Println(len(urls), err)
	for _, u := range urls {
		fmt.Println(u)
	}
	// Output:
	// 2 <nil>
	// https://shop.test
	// https://shop.test/menu
}

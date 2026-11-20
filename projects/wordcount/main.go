// 專題 2：字詞統計。
//
// 讀一個或多個檔案（沒有給檔名時讀標準輸入），算出行數、字詞數、字元數和 byte 數，
// 再列出出現最多次的字詞。英文以空白和標點分開，中文依 -han 切成單字或雙字。
//
// 執行：go run ./projects/wordcount projects/wordcount/testdata/*.txt
//
//	go run ./projects/wordcount -han bigram -top 3 projects/wordcount/testdata/notice.txt
//	echo "Go 很好玩，Go 很快" | go run ./projects/wordcount
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

// run 是整支程式：args 不含程式名稱，回傳 exit code。
// 某個檔案讀不到時，在 stderr 回報、繼續處理其他檔案，最後回傳 1。
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("wordcount", flag.ContinueOnError)
	fs.SetOutput(stderr)
	top := fs.Int("top", 5, "列出出現最多次的前幾個字詞")
	mode := HanChar
	fs.TextVar(&mode, "han", HanChar, "中文的切法：char（單字）或 bigram（雙字）")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	files := fs.Args()
	if len(files) == 0 {
		files = []string{"-"}
	}
	var total Stats
	status := 0
	for _, name := range files {
		s, err := countFile(name, stdin, mode)
		if err != nil {
			fmt.Fprintf(stderr, "wordcount：%v\n", err)
			status = 1
			continue
		}
		printStats(stdout, name, s)
		total.Add(s)
	}
	if len(files) > 1 {
		printStats(stdout, "合計", total)
	}
	printTop(stdout, Top(total.Counts, *top))
	return status
}

// countFile 統計一個檔案；name 是 "-" 時讀 stdin。
func countFile(name string, stdin io.Reader, mode HanMode) (Stats, error) {
	if name == "-" {
		return Count(stdin, mode)
	}
	f, err := os.Open(name)
	if err != nil {
		return Stats{}, err
	}
	defer f.Close()
	s, err := Count(f, mode)
	if err != nil {
		return s, fmt.Errorf("%s：%w", name, err)
	}
	return s, nil
}

func printStats(w io.Writer, name string, s Stats) {
	if name == "-" {
		name = "標準輸入"
	}
	fmt.Fprintf(w, "%s：%d 行，%d 個字詞，%d 個字元，%d bytes\n", name, s.Lines, s.Words, s.Runes, s.Bytes)
}

// printTop 印出排行榜，字詞欄依顯示寬度對齊。
func printTop(w io.Writer, top []WordCount) {
	if len(top) == 0 {
		return
	}
	col := 0
	for _, wc := range top {
		col = max(col, displayWidth(wc.Word))
	}
	fmt.Fprintln(w, "出現最多的字詞：")
	for _, wc := range top {
		fmt.Fprintf(w, "  %s  %d\n", padRight(wc.Word, col), wc.Count)
	}
}

// displayWidth 回傳 s 在等寬字型的終端機上大約佔幾格：中日韓文字和全形符號佔 2 格。
func displayWidth(s string) int {
	w := 0
	for _, r := range s {
		w++
		if unicode.In(r, unicode.Han, unicode.Hiragana, unicode.Katakana, unicode.Hangul) ||
			r >= 0x3000 && r <= 0x303F || r >= 0xFF01 && r <= 0xFF60 {
			w++
		}
	}
	return w
}

// padRight 在 s 後面補空白，讓它佔 width 格。
func padRight(s string, width int) string {
	return s + strings.Repeat(" ", max(width-displayWidth(s), 0))
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

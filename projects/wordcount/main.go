// 專題 2：字詞統計。
//
// 從標準輸入讀文字，算出行數、字數、字元數，以及出現最多次的字詞。
// 用到：bufio.Scanner（第 02 章）、map（第 07 章）、rune（第 08 章）、排序。
//
// 執行：go run ./projects/wordcount -top 3 < README.md
package main

import (
	"bufio"
	"cmp"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Stats 是統計結果。
type Stats struct {
	Lines  int
	Words  int
	Runes  int // 不含換行
	Counts map[string]int
}

// WordCount 是一個字詞和它出現的次數。
type WordCount struct {
	Word  string
	Count int
}

// Count 一行一行讀 r，用空白切出字詞，轉小寫並去掉前後的標點後再計數。
// 注意：中文句子之間通常沒有空白，所以一整串中文會被當成一個字詞。
func Count(r io.Reader) (Stats, error) {
	s := Stats{Counts: make(map[string]int)}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		s.Lines++
		s.Runes += utf8.RuneCountInString(line)
		for _, w := range strings.Fields(line) {
			w = strings.ToLower(strings.TrimFunc(w, unicode.IsPunct))
			if w == "" {
				continue
			}
			s.Words++
			s.Counts[w]++
		}
	}
	return s, sc.Err()
}

// Top 回傳出現次數最多的 n 個字詞；次數相同時依字母順序，讓結果固定。
func Top(counts map[string]int, n int) []WordCount {
	all := make([]WordCount, 0, len(counts))
	for w, c := range counts {
		all = append(all, WordCount{w, c})
	}
	slices.SortFunc(all, func(a, b WordCount) int {
		return cmp.Or(
			cmp.Compare(b.Count, a.Count), // 次數多的在前
			cmp.Compare(a.Word, b.Word),   // 次數一樣就照字母
		)
	})
	return all[:min(n, len(all))]
}

// run 讀 in、寫結果到 out；top 是要列出幾個字詞。
func run(in io.Reader, out io.Writer, top int) error {
	s, err := Count(in)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "行數 %d，字數 %d，字元數 %d\n", s.Lines, s.Words, s.Runes)
	for _, wc := range Top(s.Counts, top) {
		fmt.Fprintf(out, "%-8s %d\n", wc.Word, wc.Count)
	}
	return nil
}

func main() {
	top := flag.Int("top", 5, "列出出現最多次的前幾個字詞")
	flag.Parse()
	if err := run(os.Stdin, os.Stdout, *top); err != nil {
		fmt.Fprintln(os.Stderr, "錯誤：", err)
		os.Exit(1)
	}
}

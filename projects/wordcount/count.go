package main

import (
	"bufio"
	"cmp"
	"io"
	"slices"
	"strings"
	"unicode/utf8"
)

// Stats 是一份輸入的統計結果。
type Stats struct {
	Lines  int
	Words  int
	Runes  int // 不含行尾的換行
	Bytes  int // 含換行，和檔案大小相同
	Counts map[string]int
}

// Count 讀完 r，統計行數、字詞數、字元數和 byte 數。最後一行沒有換行也算一行。
// 用 bufio.Reader 而不是 bufio.Scanner：Scanner 預設遇到超過 64KB 的一行就會停下來回報錯誤。
func Count(r io.Reader, mode HanMode) (Stats, error) {
	s := Stats{Counts: map[string]int{}}
	br := bufio.NewReader(r)
	for {
		line, err := br.ReadString('\n')
		if line != "" {
			s.Lines++
			s.Bytes += len(line)
			text := strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r")
			s.Runes += utf8.RuneCountInString(text)
			for w := range Tokens(text, mode) {
				s.Words++
				s.Counts[w]++
			}
		}
		if err == io.EOF {
			return s, nil
		}
		if err != nil {
			return s, err
		}
	}
}

// Add 把 other 的結果加進 s。s 是零值也可以。
func (s *Stats) Add(other Stats) {
	s.Lines += other.Lines
	s.Words += other.Words
	s.Runes += other.Runes
	s.Bytes += other.Bytes
	if s.Counts == nil {
		s.Counts = map[string]int{}
	}
	for w, c := range other.Counts {
		s.Counts[w] += c
	}
}

// WordCount 是一個字詞和它出現的次數。
type WordCount struct {
	Word  string
	Count int
}

// Top 回傳出現次數最多的 n 個字詞；次數相同時依字詞排序，讓結果固定。
func Top(counts map[string]int, n int) []WordCount {
	all := make([]WordCount, 0, len(counts))
	for w, c := range counts {
		all = append(all, WordCount{w, c})
	}
	slices.SortFunc(all, func(a, b WordCount) int {
		return cmp.Or(
			cmp.Compare(b.Count, a.Count), // 次數多的在前
			strings.Compare(a.Word, b.Word),
		)
	})
	return all[:min(max(n, 0), len(all))]
}

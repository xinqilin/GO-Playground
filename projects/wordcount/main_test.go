package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"
	"unicode"
)

func TestTokens(t *testing.T) {
	tests := []struct {
		line string
		mode HanMode
		want []string
	}{
		{"Go is FUN.", HanChar, []string{"go", "is", "fun"}},
		{"排骨飯 x2", HanChar, []string{"排", "骨", "飯", "x2"}},
		{"排骨飯 x2", HanBigram, []string{"排骨", "骨飯", "x2"}},
		{"飯，好吃", HanBigram, []string{"飯", "好吃"}},
		{"Go1.27的排骨飯", HanBigram, []string{"go1", "27", "的排", "排骨", "骨飯"}},
		{"Today's menu", HanChar, []string{"today", "s", "menu"}},
		{"  ，。！  ", HanChar, nil},
		{"", HanBigram, nil},
	}
	for _, tt := range tests {
		got := slices.Collect(Tokens(tt.line, tt.mode))
		if !slices.Equal(got, tt.want) {
			t.Errorf("Tokens(%q, %v) = %q，想要 %q", tt.line, tt.mode, got, tt.want)
		}
	}
}

// TestTokensStop 確認呼叫端 break 之後，Tokens 不會再呼叫 yield（否則會 panic）。
func TestTokensStop(t *testing.T) {
	for _, mode := range []HanMode{HanChar, HanBigram} {
		var got []string
		for w := range Tokens("排骨飯 go 雞腿飯", mode) {
			got = append(got, w)
			if len(got) == 2 {
				break
			}
		}
		if len(got) != 2 {
			t.Errorf("mode %v：拿到 %q", mode, got)
		}
	}
}

func FuzzTokens(f *testing.F) {
	for _, s := range []string{"Go is fun", "排骨飯今天售完", "Today's 菜單：pork", "\xff\xfe"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, line string) {
		hans := 0
		for _, r := range line {
			if isHan(r) {
				hans++
			}
		}
		chars := 0
		for w := range Tokens(line, HanChar) {
			if w == "" || strings.ContainsFunc(w, unicode.IsSpace) {
				t.Fatalf("Tokens(%q) 產生了 %q", line, w)
			}
			if isHan([]rune(w)[0]) {
				chars++
			}
		}
		if chars != hans {
			t.Errorf("Tokens(%q, HanChar) 有 %d 個漢字字詞，原文有 %d 個漢字", line, chars, hans)
		}
	})
}

func TestCountMatchesFileSize(t *testing.T) {
	for _, name := range []string{"testdata/notice.txt", "testdata/menu.txt"} {
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		s, err := Count(strings.NewReader(string(data)), HanChar)
		if err != nil {
			t.Fatal(err)
		}
		if s.Bytes != len(data) {
			t.Errorf("%s：Bytes = %d，檔案有 %d bytes", name, s.Bytes, len(data))
		}
	}
}

func TestCountLines(t *testing.T) {
	tests := []struct {
		in           string
		lines, runes int
	}{
		{"", 0, 0},
		{"a", 1, 1},
		{"a\n", 1, 1},
		{"a\nb", 2, 2},
		{"排骨\r\n飯\r\n", 2, 3},
		{"\n\n", 2, 0},
	}
	for _, tt := range tests {
		s, err := Count(strings.NewReader(tt.in), HanChar)
		if err != nil || s.Lines != tt.lines || s.Runes != tt.runes || s.Bytes != len(tt.in) {
			t.Errorf("Count(%q) = %d 行 %d 字元 %d bytes，%v；想要 %d 行 %d 字元 %d bytes", tt.in, s.Lines, s.Runes, s.Bytes, err, tt.lines, tt.runes, len(tt.in))
		}
	}
}

// TestLongLine 確認超過 64KB 的一行也讀得完（bufio.Scanner 預設會失敗）。
func TestLongLine(t *testing.T) {
	line := strings.Repeat("飯", 100_000)
	s, err := Count(strings.NewReader(line+"\n"), HanChar)
	if err != nil || s.Counts["飯"] != 100_000 {
		t.Errorf("Count = %d 個「飯」，%v", s.Counts["飯"], err)
	}
}

func TestStatsAdd(t *testing.T) {
	var total Stats
	total.Add(Stats{Lines: 1, Words: 2, Runes: 3, Bytes: 4, Counts: map[string]int{"飯": 2}})
	total.Add(Stats{Lines: 1, Words: 1, Runes: 1, Bytes: 2, Counts: map[string]int{"飯": 1, "go": 1}})
	if total.Lines != 2 || total.Words != 3 || total.Runes != 4 || total.Bytes != 6 || total.Counts["飯"] != 3 || total.Counts["go"] != 1 {
		t.Errorf("total = %+v", total)
	}
}

func TestTop(t *testing.T) {
	counts := map[string]int{"b": 2, "a": 2, "飯": 5, "d": 1}
	want := []WordCount{{"飯", 5}, {"a", 2}, {"b", 2}}
	if got := Top(counts, 3); !slices.Equal(got, want) {
		t.Errorf("Top(3) = %v，想要 %v", got, want)
	}
	if got := Top(counts, 10); len(got) != 4 {
		t.Errorf("Top(10) 有 %d 個，想要全部 4 個", len(got))
	}
	if got := Top(counts, -1); len(got) != 0 {
		t.Errorf("Top(-1) = %v，想要空的", got)
	}
}

func TestRunExitCodes(t *testing.T) {
	tests := []struct {
		args       []string
		wantStatus int
		wantErr    string
	}{
		{[]string{"testdata/notice.txt"}, 0, ""},
		{[]string{"testdata/notice.txt", "testdata/nosuch.txt"}, 1, "wordcount：open testdata/nosuch.txt"},
		{[]string{"-han", "word", "testdata/notice.txt"}, 2, `invalid value "word" for flag -han`},
	}
	for _, tt := range tests {
		var stdout, stderr strings.Builder
		if got := run(tt.args, nil, &stdout, &stderr); got != tt.wantStatus || !strings.Contains(stderr.String(), tt.wantErr) {
			t.Errorf("run(%q) = %d，stderr %q；想要 %d 和 %q", tt.args, got, stderr.String(), tt.wantStatus, tt.wantErr)
		}
	}
}

func Example_run() {
	run([]string{"testdata/notice.txt", "testdata/menu.txt"}, nil, os.Stdout, os.Stdout)
	// Output:
	// testdata/notice.txt：3 行，25 個字詞，29 個字元，90 bytes
	// testdata/menu.txt：3 行，27 個字詞，87 個字元，123 bytes
	// 合計：6 行，52 個字詞，116 個字元，213 bytes
	// 出現最多的字詞：
	//   飯    6
	//   rice  3
	//   排    3
	//   骨    3
	//   chop  2
}

func Example_runBigram() {
	run([]string{"-han", "bigram", "-top", "3", "testdata/notice.txt"}, nil, os.Stdout, os.Stdout)
	// Output:
	// testdata/notice.txt：3 行，21 個字詞，29 個字元，90 bytes
	// 出現最多的字詞：
	//   排骨  2
	//   骨飯  2
	//   今天  1
}

func Example_runStdin() {
	run(nil, strings.NewReader("Go 很好玩，Go 很快\n"), os.Stdout, os.Stdout)
	// Output:
	// 標準輸入：1 行，7 個字詞，12 個字元，25 bytes
	// 出現最多的字詞：
	//   go  2
	//   很  2
	//   好  1
	//   快  1
	//   玩  1
}

func Example_runMissing() {
	status := run([]string{"-top", "0", "testdata/notice.txt", "testdata/nosuch.txt"}, nil, os.Stdout, os.Stdout)
	fmt.Println("exit code", status)
	// Output:
	// testdata/notice.txt：3 行，25 個字詞，29 個字元，90 bytes
	// wordcount：open testdata/nosuch.txt: no such file or directory
	// 合計：3 行，25 個字詞，29 個字元，90 bytes
	// exit code 1
}

func Example_tokens() {
	line := "Today's 便當：排骨飯 x2"
	fmt.Printf("%q\n", slices.Collect(Tokens(line, HanChar)))
	fmt.Printf("%q\n", slices.Collect(Tokens(line, HanBigram)))
	// Output:
	// ["today" "s" "便" "當" "排" "骨" "飯" "x2"]
	// ["today" "s" "便當" "排骨" "骨飯" "x2"]
}

// Example_scannerLimit 示範 bufio.Scanner 遇到超過 64KB 的一行會停下來。
func Example_scannerLimit() {
	long := strings.Repeat("飯", 30_000) // 90,000 bytes
	sc := bufio.NewScanner(strings.NewReader(long + "\n"))
	for sc.Scan() {
		fmt.Println("讀到一行")
	}
	fmt.Printf("錯誤：%v\n", sc.Err())
	// Output:
	// 錯誤：bufio.Scanner: token too long
}

func Example_padRight() {
	for _, w := range []string{"go", "排骨飯", "rice"} {
		fmt.Printf("|%-6s|%s|\n", w, padRight(w, 6))
	}
	// Output:
	// |go    |go    |
	// |排骨飯   |排骨飯|
	// |rice  |rice  |
}

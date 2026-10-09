package main

import (
	"os"
	"slices"
	"strings"
	"testing"
)

const sample = `Go is fun. Go is fast!
The gopher says: go, go, GO.
`

func Example() {
	if err := run(strings.NewReader(sample), os.Stdout, 3); err != nil {
		panic(err)
	}
	// Output:
	// 行數 2，字數 12，字元數 50
	// go       5
	// is       2
	// fast     1
}

func TestCount(t *testing.T) {
	s, err := Count(strings.NewReader(sample))
	if err != nil {
		t.Fatal(err)
	}
	if s.Lines != 2 || s.Words != 12 {
		t.Errorf("lines=%d words=%d, want 2 and 12", s.Lines, s.Words)
	}
	if s.Counts["go"] != 5 {
		t.Errorf(`counts["go"] = %d, want 5（大小寫、標點都要正規化）`, s.Counts["go"])
	}
}

func TestTop(t *testing.T) {
	counts := map[string]int{"b": 2, "a": 2, "c": 5, "d": 1}
	got := Top(counts, 3)
	want := []WordCount{{"c", 5}, {"a", 2}, {"b", 2}}
	if !slices.Equal(got, want) {
		t.Errorf("Top = %v, want %v", got, want)
	}
	if got := Top(counts, 10); len(got) != 4 {
		t.Errorf("n 大於字詞數時應該全部回傳，got %d 個", len(got))
	}
}

func TestEmptyInput(t *testing.T) {
	var out strings.Builder
	if err := run(strings.NewReader(""), &out, 5); err != nil {
		t.Fatal(err)
	}
	if got := out.String(); got != "行數 0，字數 0，字元數 0\n" {
		t.Errorf("got %q", got)
	}
}

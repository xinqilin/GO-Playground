package main

import (
	"errors"
	"strings"
	"testing"
)

func Example_printSpacing() {
	printSpacing()
	// Output:
	// 便當3 85元
	// 便當 3 85 元
}

func Example_sprintDemo() {
	sprintDemo()
	// Output:
	// 排骨飯 x2：170 元
}

func Example_fprintDemo() {
	fprintDemo()
	// Output:
	// 寫到標準輸出
	// Builder 裡有：42 號請到 3 號櫃台
}

func Example_verbsGeneral() {
	verbsGeneral()
	// Output:
	// {排骨飯 85} | {Name:排骨飯 Price:85}
	// main.item{Name:"排骨飯", Price:85}
	// main.item
	// map[排骨:85 滷蛋:15 雞腿:90]
	// [A01 A02] ["A01" "A02"]
}

func Example_verbsInt() {
	verbsInt()
	// Output:
	// 255 +255 11111111 377 ff FF 0xff
	// 飯 '飯' U+98EF 39151
}

func Example_verbsFloat() {
	verbsFloat()
	// Output:
	// 1234.567800 | 1234.57 | 1235
	// 1.234568e+03 | 1234.5678
	// 1.2345e-05 | 1e+21
}

func Example_verbsString() {
	verbsString()
	// Output:
	// Go語 | "Go語" | 476fe8aa9e | 47 6f e8 aa 9e
	// true | false
	// 85%
}

func Example_widthBasics() {
	widthBasics()
	// Output:
	// |    85|85    |000085|
	// |    3.14|3.14    |3.1|
	// |    Go|Go    |Go|
	// |    85|
}

func Example_fareTable() {
	fareTable()
	// Output:
	// 票價  目的地
	//   20  台北車站
	//   35  南港展覽館
	//   65  淡水
}

func Example_widthCJK() {
	widthCJK()
	// Output:
	// |Go    |
	// |排骨飯   |
}

func Example_stringerDemo() {
	stringerDemo()
	// Output:
	// 大份
	// 大份 大份 2
	// Size(7)
	// [小份 大份]
}

func Example_streams() {
	streams()
	// Output:
	// 結果：這一行寫到 stdout
}

func Example_scanLines() {
	scanLines()
	// Output:
	// 第 1 行："排骨飯 85"
	// 第 2 行："雞腿飯 90"
	// 第 3 行：""
	// 第 4 行："滷蛋 15"
}

func Example_scanWords() {
	scanWords()
	// Output:
	// 4 [排骨飯 2 滷蛋 3]
}

func Example_scanTooLong() {
	scanTooLong()
	// Output:
	// 讀到 12 bytes
	// Err: bufio.Scanner: token too long
}

func Example_scanBigBuffer() {
	scanBigBuffer()
	// Output:
	// 讀到 70000 bytes
	// 讀到 12 bytes
	// Err: <nil>
}

func Example_sscanDemo() {
	sscanDemo()
	// Output:
	// 3 <nil> 排骨飯 170
	// 1 expected integer
}

func Example_numberLinesDemo() {
	numberLinesDemo()
	// Output:
	// 01 排骨飯
	// 02 雞腿飯
	// 03 滷蛋
	// "01 紅茶\n02 綠茶\n"
}

func Example_bufferedWriter() {
	bufferedWriter()
	// Output:
	// Flush 前：""
	// Flush 後："叫號：A01"
}

func Example_copyDemo() {
	copyDemo()
	// Output:
	// 整份菜單原封不動地複製過去
	// 40 <nil>
}

func TestSizeString(t *testing.T) {
	tests := []struct {
		s    Size
		want string
	}{
		{Small, "小份"}, {Medium, "中份"}, {Large, "大份"}, {Size(-1), "Size(-1)"},
	}
	for _, tt := range tests {
		if got := tt.s.String(); got != tt.want {
			t.Errorf("Size(%d).String() = %q, want %q", int(tt.s), got, tt.want)
		}
	}
}

// failWriter 每次寫入都失敗，用來確認錯誤有被傳回來。
type failWriter struct{}

func (failWriter) Write(p []byte) (int, error) { return 0, errors.New("磁碟滿了") }

func TestNumberLines(t *testing.T) {
	var sb strings.Builder
	if err := numberLines(strings.NewReader("a\r\nb\n\nc"), &sb); err != nil {
		t.Fatal(err)
	}
	if got, want := sb.String(), "01 a\n02 b\n03 \n04 c\n"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if err := numberLines(strings.NewReader("a\n"), failWriter{}); err == nil {
		t.Error("寫入失敗時應該回傳錯誤")
	}
}

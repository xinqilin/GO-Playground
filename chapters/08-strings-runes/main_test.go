package main

import (
	"strings"
	"testing"
)

func Example_stringBasics() {
	stringBasics()
	// Output:
	// 9 72 H
	// Go Hello
	// Hello，便當店
	// true true true
}

func Example_literals() {
	literals()
	// Output:
	// 第一行
	// 第二行
	// 價格	85
	// C:\menu\new.txt
	// SELECT name, price
	// FROM menu
	// WHERE price < 100
	// 語言
}

func Example_bytesVsRunes() {
	bytesVsRunes()
	// Output:
	// 8 4
	// 71 232
	// 47 6f e8 aa 9e e8 a8 80
}

func Example_runeLiteral() {
	runeLiteral()
	// Output:
	// 35486
	// 語 U+8A9E int32
	// 3 1
	// 語言
}

func Example_cutInMiddle() {
	cutInMiddle()
	// Output:
	// "Go\xe8" false
	// "Go語" true
}

func Example_rangeString() {
	rangeString()
	// Output:
	// 0 G 71
	// 1 o 111
	// 2 語 35486
	// 5 言 35328
}

func Example_rangeInvalid() {
	rangeInvalid()
	// Output:
	// 0 U+0041
	// 1 U+FFFD
	// 2 U+0042
}

func Example_conversions() {
	conversions()
	// Output:
	// 8 4
	// [71 111] [35486 35328]
	// 語 語言
	// Go語言 go語言
}

func Example_editRunes() {
	editRunes()
	// Output:
	// 中杯紅茶
}

func Example_numberToString() {
	numberToString()
	// Output:
	// 65 A
	// 65元
}

func Example_searchFuncs() {
	searchFuncs()
	// Output:
	// true true true
	// 9 -1
	// 2
	// true
}

func Example_splitJoin() {
	splitJoin()
	// Output:
	// ["排骨" "雞腿" "" "素食"]
	// ["排骨" "雞腿" "素食"]
	// 排骨、雞腿、素食
}

func Example_trimCut() {
	trimCut()
	// Output:
	// "排骨飯"
	// AB12345678 menu
	// 特價
	// 排骨 2 true
	// false
}

func Example_replaceFuncs() {
	replaceFuncs()
	// Output:
	// 0912345678
	// a+b-c
	// ==========
	// GO 語言 hello
	// 臺北臺中
}

func Example_builderDemo() {
	builderDemo()
	// Output:
	// 排骨 85
	// 雞腿 90
	// 合計 175
}

func Example_parseNumbers() {
	parseNumbers()
	// Output:
	// 90 <nil>
	// strconv.Atoi: parsing "85元": invalid syntax
	// 7
	// 2147483647 strconv.ParseInt: parsing "99999999999": value out of range
	// 255 true
}

func Example_formatNumbers() {
	formatNumbers()
	// Output:
	// 85元
	// 11111111 ff
	// 3.14
	// "排骨\n飯" "\u6392\u9aa8"
}

func Example_classify() {
	classify()
	// Output:
	// 'A' 字母=true 數字=false 漢字=false 空白=false 標點=false
	// '9' 字母=false 數字=true 漢字=false 空白=false 標點=false
	// '語' 字母=true 數字=false 漢字=true 空白=false 標點=false
	// ' ' 字母=false 數字=false 漢字=false 空白=true 標點=false
	// '，' 字母=false 數字=false 漢字=false 空白=false 標點=true
}

func Example_halfwidth() {
	halfwidth()
	// Output:
	// AB-12345678
}

func Example_alignedMenu() {
	alignedMenu()
	// Output:
	// |Coffee  |
	// |排骨飯     |
	// |雞腿便當    |
	// |Coffee  |
	// |排骨飯  |
	// |雞腿便當|
}

func TestDisplayWidth(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{"", 0},
		{"Coffee", 6},
		{"排骨飯", 6},
		{"A餐，85元", 9},
		{"ＡＢ", 4},
		{"すし", 4},
	}
	for _, tt := range tests {
		if got := displayWidth(tt.s); got != tt.want {
			t.Errorf("displayWidth(%q) = %d, want %d", tt.s, got, tt.want)
		}
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		s     string
		width int
		want  string
	}{
		{"排骨飯", 8, "排骨飯  "},
		{"Coffee", 8, "Coffee  "},
		{"雞腿便當", 8, "雞腿便當"},
		{"雞腿便當", 4, "雞腿便當"},
	}
	for _, tt := range tests {
		if got := padRight(tt.s, tt.width); got != tt.want {
			t.Errorf("padRight(%q, %d) = %q, want %q", tt.s, tt.width, got, tt.want)
		}
	}
}

func TestToHalfwidth(t *testing.T) {
	if got, want := toHalfwidth("ＡＢ１２３　ｇｏ！"), "AB123 go!"; got != want {
		t.Errorf("toHalfwidth = %q, want %q", got, want)
	}
	if got := toHalfwidth("排骨飯、85"); got != "排骨飯、85" {
		t.Errorf("中文字和「、」不是全形英數字，不應該改變，結果是 %q", got)
	}
}

func TestReplaceAt(t *testing.T) {
	if got := replaceAt("Go語言", 3, '法'); got != "Go語法" {
		t.Errorf("replaceAt = %q, want %q", got, "Go語法")
	}
}

func TestConcat(t *testing.T) {
	parts := []string{"排骨", "、", "雞腿"}
	if a, b := concatPlus(parts), concatBuilder(parts); a != "排骨、雞腿" || b != a {
		t.Errorf("concatPlus = %q, concatBuilder = %q", a, b)
	}
	many := strings.Split(strings.Repeat("便當,", 100), ",")
	if allocs := testing.AllocsPerRun(10, func() { concatBuilder(many) }); allocs != 1 {
		t.Errorf("concatBuilder 配置了 %v 次，want 1", allocs)
	}
	if allocs := testing.AllocsPerRun(10, func() { concatPlus(many) }); allocs < 50 {
		t.Errorf("concatPlus 只配置了 %v 次，預期每次串接都會配置", allocs)
	}
}

func BenchmarkConcatPlus(b *testing.B) {
	parts := strings.Split(strings.Repeat("便當,", 1000), ",")
	for b.Loop() {
		concatPlus(parts)
	}
}

func BenchmarkConcatBuilder(b *testing.B) {
	parts := strings.Split(strings.Repeat("便當,", 1000), ",")
	for b.Loop() {
		concatBuilder(parts)
	}
}

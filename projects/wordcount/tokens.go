package main

import (
	"fmt"
	"iter"
	"strings"
	"unicode"
	"unicode/utf8"
)

// HanMode 決定連續的漢字怎麼切成字詞。
type HanMode int

const (
	HanChar   HanMode = iota // 每個漢字是一個字詞
	HanBigram                // 每兩個相鄰的漢字是一個字詞
)

// MarshalText 讓 HanMode 可以用在 flag.TextVar，也決定 -help 顯示的預設值。
func (m HanMode) MarshalText() ([]byte, error) {
	switch m {
	case HanChar:
		return []byte("char"), nil
	case HanBigram:
		return []byte("bigram"), nil
	}
	return nil, fmt.Errorf("不認得的 HanMode %d", int(m))
}

// UnmarshalText 解析 -han 旗標的值。
func (m *HanMode) UnmarshalText(text []byte) error {
	switch string(text) {
	case "char":
		*m = HanChar
	case "bigram":
		*m = HanBigram
	default:
		return fmt.Errorf("要用 char 或 bigram，收到 %q", text)
	}
	return nil
}

// Tokens 依序產生 line 裡的字詞：
//   - 連續的字母和數字是一個字詞，轉成小寫。
//   - 連續的漢字依 mode 切成單字或雙字。
//   - 其他字元（空白、標點、符號）都是分隔。
func Tokens(line string, mode HanMode) iter.Seq[string] {
	return func(yield func(string) bool) {
		rest := line
		for rest != "" {
			r, size := utf8.DecodeRuneInString(rest)
			switch {
			case isHan(r):
				end := runEnd(rest, isHan)
				if !yieldHan(rest[:end], mode, yield) {
					return
				}
				rest = rest[end:]
			case isWordRune(r):
				end := runEnd(rest, isWordRune)
				if !yield(strings.ToLower(rest[:end])) {
					return
				}
				rest = rest[end:]
			default:
				rest = rest[size:]
			}
		}
	}
}

func isHan(r rune) bool { return unicode.Is(unicode.Han, r) }

// isWordRune 回傳 r 是不是英數字詞的一部分。漢字也是 Letter，所以要先排除。
func isWordRune(r rune) bool {
	return !isHan(r) && (unicode.IsLetter(r) || unicode.IsDigit(r))
}

// runEnd 回傳 s 開頭連續符合 keep 的字元到哪裡結束（byte 位置）。
func runEnd(s string, keep func(rune) bool) int {
	for i, r := range s {
		if !keep(r) {
			return i
		}
	}
	return len(s)
}

// yieldHan 把一段連續的漢字依 mode 切開交給 yield；yield 回傳 false 時跟著回傳 false。
// 雙字模式下，只有一個字的段落就是那一個字。
func yieldHan(run string, mode HanMode, yield func(string) bool) bool {
	runes := []rune(run)
	if mode == HanChar || len(runes) == 1 {
		for _, r := range runes {
			if !yield(string(r)) {
				return false
			}
		}
		return true
	}
	for i := range len(runes) - 1 {
		if !yield(string(runes[i : i+2])) {
			return false
		}
	}
	return true
}

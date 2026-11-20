//go:build solution

package toolbox

import (
	"errors"
	"go/ast"
	"go/build/constraint"
	"go/parser"
	"go/token"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"
)

// BinaryName 回傳交叉編譯時執行檔的檔名：「名稱-GOOS-GOARCH」，Windows 再加上 .exe。
func BinaryName(name, goos, goarch string) string {
	out := name + "-" + goos + "-" + goarch
	if goos == "windows" {
		out += ".exe"
	}
	return out
}

// Matches 回傳 //go:build 後面的條件 expr 在 tags 這組標籤下成不成立。
func Matches(expr string, tags []string) (bool, error) {
	e, err := constraint.Parse("//go:build " + expr)
	if err != nil {
		return false, err
	}
	return e.Eval(func(tag string) bool { return slices.Contains(tags, tag) }), nil
}

// FileMatches 回傳 name 這個 .go 檔在 goos/goarch 下會不會被編譯，只看檔名。
func FileMatches(name, goos, goarch string) bool {
	if strings.HasPrefix(name, "_") || strings.HasPrefix(name, ".") {
		return false
	}
	name, _, _ = strings.Cut(name, ".")
	name = strings.TrimSuffix(name, "_test")
	_, rest, ok := strings.Cut(name, "_")
	if !ok {
		return true
	}
	parts := strings.Split(rest, "_")
	n := len(parts)
	if n >= 2 && KnownOS[parts[n-2]] && KnownArch[parts[n-1]] {
		return parts[n-2] == goos && parts[n-1] == goarch
	}
	last := parts[n-1]
	if KnownOS[last] {
		return last == goos
	}
	if KnownArch[last] {
		return last == goarch
	}
	return true
}

// Version 決定程式要顯示的版本。
func Version(info *debug.BuildInfo, injected string) string {
	if injected != "" && injected != "dev" {
		return injected
	}
	if info == nil {
		return "dev"
	}
	if v := info.Main.Version; v != "" && v != "(devel)" {
		return v
	}
	var rev string
	var dirty bool
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			rev = s.Value
		case "vcs.modified":
			dirty = s.Value == "true"
		}
	}
	if rev == "" {
		return "dev"
	}
	if len(rev) > 7 {
		rev = rev[:7]
	}
	if dirty {
		rev += "-dirty"
	}
	return rev
}

// GenerateArgs 解析一行 //go:generate 指令，回傳要執行的指令和參數。
func GenerateArgs(line string) ([]string, error) {
	rest, ok := strings.CutPrefix(line, "//go:generate")
	if !ok || rest == "" || (rest[0] != ' ' && rest[0] != '\t') {
		return nil, ErrNotDirective
	}
	var args []string
	for {
		rest = strings.TrimLeft(rest, " \t")
		if rest == "" {
			break
		}
		if rest[0] != '"' {
			end := strings.IndexAny(rest, " \t")
			if end < 0 {
				end = len(rest)
			}
			args = append(args, rest[:end])
			rest = rest[end:]
			continue
		}
		end := closingQuote(rest)
		if end < 0 {
			return nil, errors.New("引號沒有成對")
		}
		arg, err := strconv.Unquote(rest[:end+1])
		if err != nil {
			return nil, err
		}
		rest = rest[end+1:]
		if rest != "" && rest[0] != ' ' && rest[0] != '\t' {
			return nil, errors.New("右引號後面要有空白")
		}
		args = append(args, arg)
	}
	if len(args) == 0 {
		return nil, errors.New("沒有任何參數")
	}
	return args, nil
}

// closingQuote 回傳 s（以 " 開頭）裡右引號的位置，跳過反斜線跳脫的字元；找不到時回傳 -1。
func closingQuote(s string) int {
	for i := 1; i < len(s); i++ {
		switch s[i] {
		case '\\':
			i++
		case '"':
			return i
		}
	}
	return -1
}

// IsGenerated 回傳 src 是不是程式產生的 Go 原始檔。
func IsGenerated(src []byte) bool {
	f, err := parser.ParseFile(token.NewFileSet(), "", src, parser.PackageClauseOnly|parser.ParseComments)
	if err != nil {
		return false
	}
	return ast.IsGenerated(f)
}

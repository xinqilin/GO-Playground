package toolbox

import (
	"errors"
	"go/build"
	"io"
	"os"
	"runtime/debug"
	"slices"
	"strings"
	"testing"
)

func TestBinaryName(t *testing.T) {
	tests := []struct{ name, goos, goarch, want string }{
		{"pos", "linux", "arm64", "pos-linux-arm64"},
		{"pos", "darwin", "amd64", "pos-darwin-amd64"},
		{"pos", "windows", "amd64", "pos-windows-amd64.exe"},
		{"點餐機", "windows", "386", "點餐機-windows-386.exe"},
	}
	for _, tt := range tests {
		if got := BinaryName(tt.name, tt.goos, tt.goarch); got != tt.want {
			t.Errorf("BinaryName(%q, %q, %q) = %q，想要 %q", tt.name, tt.goos, tt.goarch, got, tt.want)
		}
	}
}

func TestMatches(t *testing.T) {
	tests := []struct {
		expr string
		tags []string
		want bool
	}{
		{"linux && !cgo", []string{"linux"}, true},
		{"linux && !cgo", []string{"linux", "cgo"}, false},
		{"pro || enterprise", []string{"enterprise"}, true},
		{"pro || enterprise", nil, false},
		{"!pro", nil, true},
		{"(linux || darwin) && arm64", []string{"darwin", "arm64"}, true},
		{"(linux || darwin) && arm64", []string{"windows", "arm64"}, false},
	}
	for _, tt := range tests {
		got, err := Matches(tt.expr, tt.tags)
		if err != nil || got != tt.want {
			t.Errorf("Matches(%q, %q) = %v, %v，想要 %v, nil", tt.expr, tt.tags, got, err, tt.want)
		}
	}
	for _, expr := range []string{"linux &&", "pro enterprise", "(linux", ""} {
		if _, err := Matches(expr, nil); err == nil {
			t.Errorf("Matches(%q) 想要錯誤，得到 nil", expr)
		}
	}
}

var fileNames = []string{
	"main.go", "linux.go", "open_darwin.go", "open_windows.go", "open_windows_test.go",
	"price_amd64.go", "price_linux_amd64.go", "price_linux_arm64_test.go",
	"price_arm64_linux.go", "tools_unix.go", "menu_gen.go", "menu_linux.pb.go",
	"linux_test.go", "_scratch.go", ".hidden.go",
}

type platform struct{ goos, goarch string }

var platforms = []platform{{"darwin", "arm64"}, {"linux", "amd64"}, {"linux", "arm64"}, {"windows", "amd64"}}

func TestFileMatches(t *testing.T) {
	// 每個檔名會在哪些平台被編譯（platforms 的索引）。
	want := map[string][]int{
		"main.go":                   {0, 1, 2, 3},
		"linux.go":                  {0, 1, 2, 3},
		"open_darwin.go":            {0},
		"open_windows.go":           {3},
		"open_windows_test.go":      {3},
		"price_amd64.go":            {1, 3},
		"price_linux_amd64.go":      {1},
		"price_linux_arm64_test.go": {2},
		"price_arm64_linux.go":      {1, 2},
		"tools_unix.go":             {0, 1, 2, 3},
		"menu_gen.go":               {0, 1, 2, 3},
		"menu_linux.pb.go":          {1, 2},
		"linux_test.go":             {0, 1, 2, 3},
		"_scratch.go":               {},
		".hidden.go":                {},
	}
	for _, name := range fileNames {
		for i, p := range platforms {
			w := slices.Contains(want[name], i)
			if got := FileMatches(name, p.goos, p.goarch); got != w {
				t.Errorf("FileMatches(%q, %q, %q) = %v，想要 %v", name, p.goos, p.goarch, got, w)
			}
		}
	}
}

// TestFileMatchesTable 確認上面的答案和 go/build 本身的判斷一致。
func TestFileMatchesTable(t *testing.T) {
	for _, name := range fileNames {
		for _, p := range platforms {
			ctxt := build.Context{
				GOOS: p.goos, GOARCH: p.goarch, Compiler: "gc",
				OpenFile: func(string) (io.ReadCloser, error) {
					return io.NopCloser(strings.NewReader("package p\n")), nil
				},
			}
			w, err := ctxt.MatchFile(".", name)
			if err != nil {
				t.Fatal(err)
			}
			if got := FileMatches(name, p.goos, p.goarch); got != w {
				t.Errorf("FileMatches(%q, %q, %q) = %v，go/build 是 %v", name, p.goos, p.goarch, got, w)
			}
		}
	}
}

func TestVersion(t *testing.T) {
	const rev = "e29293b9197c5e5eb11bddb0fa7a0fee9bf0680b"
	vcs := func(modified string) []debug.BuildSetting {
		return []debug.BuildSetting{
			{Key: "GOOS", Value: "linux"},
			{Key: "vcs.revision", Value: rev},
			{Key: "vcs.modified", Value: modified},
		}
	}
	tests := []struct {
		name     string
		info     *debug.BuildInfo
		injected string
		want     string
	}{
		{"注入的版本", nil, "1.2.3", "1.2.3"},
		{"注入的版本優先", &debug.BuildInfo{Main: debug.Module{Version: "v1.4.0"}}, "2.0.0", "2.0.0"},
		{"module 版本", &debug.BuildInfo{Main: debug.Module{Version: "v1.4.0"}}, "dev", "v1.4.0"},
		{"commit", &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}, Settings: vcs("false")}, "dev", "e29293b"},
		{"有改動", &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}, Settings: vcs("true")}, "", "e29293b-dirty"},
		{"短的 revision", &debug.BuildInfo{Settings: []debug.BuildSetting{{Key: "vcs.revision", Value: "abc"}}}, "", "abc"},
		{"什麼都沒有", &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}}, "dev", "dev"},
		{"沒有建置資訊", nil, "", "dev"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Version(tt.info, tt.injected); got != tt.want {
				t.Errorf("Version(…, %q) = %q，想要 %q", tt.injected, got, tt.want)
			}
		})
	}
}

func TestGenerateArgs(t *testing.T) {
	tests := []struct {
		line string
		want []string
	}{
		{"//go:generate go run ./genmenu -in menu.csv -out menu_gen.go", []string{"go", "run", "./genmenu", "-in", "menu.csv", "-out", "menu_gen.go"}},
		{"//go:generate\tstringer -type=Dish", []string{"stringer", "-type=Dish"}},
		{`//go:generate echo "排骨飯 便當" "a\tb" x`, []string{"echo", "排骨飯 便當", "a\tb", "x"}},
		{"//go:generate   go \t version  ", []string{"go", "version"}},
		{`//go:generate echo ""`, []string{"echo", ""}},
	}
	for _, tt := range tests {
		got, err := GenerateArgs(tt.line)
		if err != nil || !slices.Equal(got, tt.want) {
			t.Errorf("GenerateArgs(%q) = %q, %v，想要 %q, nil", tt.line, got, err, tt.want)
		}
	}
	for _, line := range []string{"// go:generate echo hi", " //go:generate echo hi", "//go:generateecho hi", "//go:generate", "package main"} {
		if _, err := GenerateArgs(line); !errors.Is(err, ErrNotDirective) {
			t.Errorf("GenerateArgs(%q) 的錯誤 = %v，想要 ErrNotDirective", line, err)
		}
	}
	for _, line := range []string{`//go:generate echo "abc`, `//go:generate echo "abc"def`, `//go:generate echo "\q"`, `//go:generate echo "a\"`, "//go:generate   "} {
		if _, err := GenerateArgs(line); err == nil || errors.Is(err, ErrNotDirective) {
			t.Errorf("GenerateArgs(%q) 的錯誤 = %v，想要 ErrNotDirective 以外的錯誤", line, err)
		}
	}
}

func TestIsGenerated(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want bool
	}{
		{"產生的檔案", "// Code generated by genmenu from menu.csv; DO NOT EDIT.\n\npackage main\n", true},
		{"前面有 build 條件", "//go:build linux\n\n// Code generated by stringer; DO NOT EDIT.\n\npackage main\n", true},
		{"和 package 說明在一起", "// Code generated by genmenu; DO NOT EDIT.\n\n// Package menu 是菜單。\npackage menu\n", true},
		{"在 package 後面", "package main\n\n// Code generated by genmenu; DO NOT EDIT.\nvar x = 1\n", false},
		{"只是提到", "// Package lint 檢查產生的檔案。\npackage lint\n\nconst marker = \"DO NOT EDIT\"\n", false},
		{"少了句點", "// Code generated by hand. DO NOT EDIT\npackage main\n", false},
		{"不在行首", "// 這不是 // Code generated by x DO NOT EDIT.\npackage main\n", false},
		{"手寫的檔案", "// Package menu 是菜單。\npackage menu\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsGenerated([]byte(tt.src)); got != tt.want {
				t.Errorf("IsGenerated = %v，想要 %v\n%s", got, tt.want, tt.src)
			}
		})
	}
	src, err := os.ReadFile("../../chapters/32-toolchain/menu_gen.go")
	if err != nil {
		t.Fatal(err)
	}
	if !IsGenerated(src) {
		t.Error("第 32 章的 menu_gen.go 應該算是產生的檔案")
	}
}

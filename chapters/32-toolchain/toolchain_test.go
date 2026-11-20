package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// 這個檔案實際執行 go build、go list、go vet、go fix，確認頁面上貼的指令和輸出是真的。
// 要編譯好幾次，加 -short 可以跳過。

// goTool 執行 go 指令，回傳 stdout 和 stderr 合併的內容。env 會加在目前的環境變數後面。
func goTool(t *testing.T, env []string, args ...string) (string, error) {
	t.Helper()
	if testing.Short() {
		t.Skip("-short：跳過要執行 go 指令的測試")
	}
	cmd := exec.Command("go", args...)
	cmd.Env = append(os.Environ(), env...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func TestFileSelection(t *testing.T) {
	tests := []struct {
		env  []string
		tags string
		want string
	}{
		{[]string{"GOOS=darwin"}, "", "[edition_free.go main.go menu.go menu_gen.go open_darwin.go] [edition_pro.go open_other.go open_windows.go]"},
		{[]string{"GOOS=windows"}, "", "[edition_free.go main.go menu.go menu_gen.go open_windows.go] [edition_pro.go open_darwin.go open_other.go]"},
		{[]string{"GOOS=linux"}, "pro", "[edition_pro.go main.go menu.go menu_gen.go open_other.go] [edition_free.go open_darwin.go open_windows.go]"},
	}
	for _, tt := range tests {
		out, err := goTool(t, tt.env, "list", "-tags="+tt.tags, "-f", "{{.GoFiles}} {{.IgnoredGoFiles}}", ".")
		if err != nil {
			t.Fatalf("%v：%v\n%s", tt.env, err, out)
		}
		if !strings.Contains(out, tt.want) {
			t.Errorf("%v -tags=%q：\n得到 %s想要 %s", tt.env, tt.tags, out, tt.want)
		}
	}
}

func TestLdflags(t *testing.T) {
	app := filepath.Join(t.TempDir(), "app")
	if out, err := goTool(t, nil, "build", "-tags", "pro", "-ldflags", "-X main.version=1.2.3", "-o", app, "."); err != nil {
		t.Fatalf("go build：%v\n%s", err, out)
	}
	out, err := exec.Command(app).Output()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "版本 1.2.3，專業版\n") {
		t.Errorf("執行結果沒有注入的版本：\n%s", out)
	}
	info, err := goTool(t, nil, "version", "-m", app)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"\tbuild\t-ldflags=\"-X main.version=1.2.3\"\n", "\tbuild\t-tags=pro\n"} {
		if !strings.Contains(info, want) {
			t.Errorf("go version -m 沒有 %q：\n%s", want, info)
		}
	}
}

func TestCrossCompile(t *testing.T) {
	tests := []struct {
		goos, goarch string
		magic        []byte
	}{
		{"linux", "amd64", []byte("\x7fELF")},
		{"windows", "amd64", []byte("MZ")},
	}
	for _, tt := range tests {
		bin := filepath.Join(t.TempDir(), "app")
		env := []string{"GOOS=" + tt.goos, "GOARCH=" + tt.goarch, "CGO_ENABLED=0"}
		if out, err := goTool(t, env, "build", "-o", bin, "."); err != nil {
			t.Fatalf("%s/%s：%v\n%s", tt.goos, tt.goarch, err, out)
		}
		data, err := os.ReadFile(bin)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.HasPrefix(data, tt.magic) {
			t.Errorf("%s/%s 的執行檔開頭是 %q，想要 %q", tt.goos, tt.goarch, data[:4], tt.magic)
		}
	}
}

func TestVetSubset(t *testing.T) {
	out, err := goTool(t, nil, "-C", "testdata/vetsubset", "test", ".")
	if err != nil {
		t.Fatalf("go test 應該通過：%v\n%s", err, out)
	}
	out, err = goTool(t, nil, "-C", "testdata/vetsubset", "vet", ".")
	if err == nil {
		t.Fatal("go vet 應該失敗")
	}
	const want = "queue.go:14:9: Next passes lock by value: example.com/vetsubset.Queue contains sync.Mutex\n"
	if !strings.Contains(out, want) {
		t.Errorf("go vet 的輸出：\n%s想要包含：\n%s", out, want)
	}
}

func TestGoFix(t *testing.T) {
	out, err := goTool(t, nil, "-C", "testdata/oldstyle", "fix", "-diff", "./...")
	var exit *exec.ExitError
	if err != nil && !errors.As(err, &exit) {
		t.Fatal(err)
	}
	for _, want := range []string{
		"+func Lookup(menu map[string]any, name string) any {\n",
		"+\tslices.Sort(prices)\n",
		"+\treturn slices.Contains(menu, name)\n",
		"+\tif after, ok := strings.CutPrefix(line, \"備註：\"); ok {\n",
		"+\t\twg.Go(func() {\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("go fix -diff 沒有建議：%q\n%s", want, out)
		}
	}
}

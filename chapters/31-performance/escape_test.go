package main

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestEscapeAnalysis 實際執行 go build -gcflags=-m，確認頁面上逃逸分析的結果沒有過時。
func TestEscapeAnalysis(t *testing.T) {
	if testing.Short() {
		t.Skip("會執行 go build，-short 時略過")
	}
	bin := filepath.Join(t.TempDir(), "perf")
	out, err := exec.Command("go", "build", "-gcflags=-m", "-o", bin, ".").CombinedOutput()
	if err != nil {
		t.Fatalf("go build：%v\n%s", err, out)
	}
	for _, want := range []string{
		"escape.go:13:2: moved to heap: t",
		"escape.go:29:15: nos does not escape",
		"escape.go:41:13: make([]int, n) escapes to heap",
		"escape.go:49:60: make([]int, 3) does not escape",
	} {
		if !strings.Contains(string(out), want) {
			t.Errorf("輸出裡沒有 %q", want)
		}
	}
}

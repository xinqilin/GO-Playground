package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestKioskLayout 實際執行頁面上對 testdata/kiosk 下的指令，確認輸出和頁面一致。
func TestKioskLayout(t *testing.T) {
	if testing.Short() {
		t.Skip("-short：跳過要執行 go 指令的測試")
	}
	run := func(args ...string) string {
		t.Helper()
		out, err := exec.Command("go", append([]string{"-C", "testdata/kiosk"}, args...)...).CombinedOutput()
		if err != nil {
			t.Fatalf("go %s：%v\n%s", strings.Join(args, " "), err, out)
		}
		return string(out)
	}

	const pkgs = "example.com/kiosk/cmd/board\nexample.com/kiosk/cmd/kiosk\nexample.com/kiosk/internal/display\nexample.com/kiosk/internal/queue\n"
	if got := run("list", "./..."); !strings.Contains(got, pkgs) {
		t.Errorf("go list ./... =\n%s想要\n%s", got, pkgs)
	}
	if got, want := run("run", "./cmd/kiosk"), "您的號碼：A001\n您的號碼：A002\n您的號碼：A003\n"; !strings.Contains(got, want) {
		t.Errorf("go run ./cmd/kiosk =\n%s想要\n%s", got, want)
	}
	if got, want := run("run", "./cmd/board"), "請 A001 到櫃台\n請 A002 到櫃台\n目前沒有人在等\n"; !strings.Contains(got, want) {
		t.Errorf("go run ./cmd/board =\n%s想要\n%s", got, want)
	}

	bin := t.TempDir()
	run("build", "-o", bin+string(filepath.Separator), "./cmd/...")
	for _, name := range []string{"kiosk", "board"} {
		matches, _ := filepath.Glob(filepath.Join(bin, name+"*"))
		if len(matches) != 1 {
			t.Errorf("go build -o 目錄/ ./cmd/... 沒有產生 %s：%v", name, matches)
		}
	}
	if _, err := os.Stat(filepath.Join(bin, "display")); err == nil {
		t.Error("internal/display 不是 main package，不應該產生執行檔")
	}
}

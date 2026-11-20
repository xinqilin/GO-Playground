package main

import (
	"os/exec"
	"strings"
	"testing"
)

// goCmd 在 testdata 底下的 dir 執行 go 指令，回傳合併的輸出和有沒有成功。
// testdata 目錄會被 ./... 略過，所以這些小 module 不會影響 repo 本身的建置。
func goCmd(t *testing.T, dir string, args ...string) (string, bool) {
	t.Helper()
	if testing.Short() {
		t.Skip("會執行 go 指令，-short 時略過")
	}
	cmd := exec.Command("go", append([]string{"-C", "testdata/" + dir}, args...)...)
	out, err := cmd.CombinedOutput()
	if _, ok := err.(*exec.ExitError); err != nil && !ok {
		t.Fatalf("執行 go 失敗：%v", err)
	}
	return string(out), err == nil
}

// TestWorkspace：shop 和 price 是兩個 module，go.work 讓 shop 直接用本機的 price。
func TestWorkspace(t *testing.T) {
	out, ok := goCmd(t, "workspace", "run", "./shop")
	want := "排骨飯 90 元\n滷肉飯 50 元\n"
	if !ok || out != want {
		t.Errorf("go run ./shop：\n%s\nwant\n%s", out, want)
	}
}

// TestReplace：沒有 go.work，靠 go.mod 裡的 replace 用本機的 price。
func TestReplace(t *testing.T) {
	out, ok := goCmd(t, "replace", "run", ".")
	if want := "排骨飯 90 元\n"; !ok || out != want {
		t.Errorf("go run .：\n%s\nwant\n%s", out, want)
	}
}

// TestInternalBoundary：kitchen 可以用自己的 internal，counter 不行。
func TestInternalBoundary(t *testing.T) {
	if out, ok := goCmd(t, "internalerr", "build", "./kitchen"); !ok {
		t.Errorf("kitchen 應該編譯得過：\n%s", out)
	}
	out, ok := goCmd(t, "internalerr", "build", "./...")
	if ok || !strings.Contains(out, "use of internal package example.com/shop/kitchen/internal/recipe not allowed") {
		t.Errorf("counter 應該因為 internal 編譯失敗：\n%s", out)
	}
}

// TestImportCycle：order 和 kitchen 互相 import，編譯失敗。
func TestImportCycle(t *testing.T) {
	out, ok := goCmd(t, "cycle", "build", "./...")
	if ok || !strings.Contains(out, "import cycle not allowed") {
		t.Errorf("應該因為 import 循環編譯失敗：\n%s", out)
	}
}

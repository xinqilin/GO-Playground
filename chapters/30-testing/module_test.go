package main

import (
	"os/exec"
	"strings"
	"testing"
)

// goTest 在 testdata 底下的 dir 執行 go test，回傳合併的輸出和有沒有成功。
// testdata 目錄會被 ./... 略過，這兩個故意失敗的小 module 不會讓 repo 的測試變紅。
func goTest(t *testing.T, dir string) (string, bool) {
	t.Helper()
	if testing.Short() {
		t.Skip("會執行 go 指令，-short 時略過")
	}
	out, err := exec.Command("go", "-C", "testdata/"+dir, "test").CombinedOutput()
	if _, ok := err.(*exec.ExitError); err != nil && !ok {
		t.Fatalf("執行 go 失敗：%v", err)
	}
	return string(out), err == nil
}

// TestFailingOutput：頁面上「測試失敗長什麼樣子」的輸出。
func TestFailingOutput(t *testing.T) {
	out, ok := goTest(t, "failing")
	for _, want := range []string{
		"--- FAIL: TestMemberPrice",
		"price_test.go:18: MemberPrice(95) = 81，want 86",
		"price_test.go:18: MemberPrice(55) = 45，want 50",
		"--- FAIL: TestMenuPrices",
		"price_test.go:27: 讀不到菜單：open testdata/menu.txt: no such file or directory",
	} {
		if ok || !strings.Contains(out, want) {
			t.Errorf("輸出裡應該有 %q：\n%s", want, out)
		}
	}
}

// TestFuzzCorpus：fuzz 找到的輸入存在 testdata/fuzz 裡，之後每次 go test 都會重跑它。
func TestFuzzCorpus(t *testing.T) {
	out, ok := goTest(t, "fuzzbug")
	for _, want := range []string{
		"--- FAIL: FuzzNormalizePlate/771e938e4458e983",
		"slice bounds out of range [:-3]",
	} {
		if ok || !strings.Contains(out, want) {
			t.Errorf("輸出裡應該有 %q：\n%s", want, out)
		}
	}
}

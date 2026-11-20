package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// TestMenuUpToDate 確認 commit 進來的 menu_gen.go 和 menu.csv 一致。
// 改了 CSV 卻忘了 go generate，這個測試就會失敗。
func TestMenuUpToDate(t *testing.T) {
	f, err := os.Open("../menu.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	want, err := generate(f, "menu.csv")
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile("../menu_gen.go")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("menu_gen.go 過期了，請執行 go generate ./chapters/32-toolchain")
	}
}

func TestGenerateErrors(t *testing.T) {
	tests := map[string]string{
		"空的":   "菜名,價格\n",
		"價格不對": "菜名,價格\n排骨飯,九十\n",
		"欄位數":  "菜名,價格\n排骨飯\n",
	}
	for name, csv := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := generate(strings.NewReader(csv), "menu.csv"); err == nil {
				t.Error("想要錯誤，得到 nil")
			}
		})
	}
}

package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// makeTree 在 dir 底下建立幾個發票檔：檔名是發票號碼，內容是金額。
// MkdirAll 會把中間缺少的目錄一次建好，目錄已經存在也不算錯誤。
func makeTree(dir string) error {
	files := map[string]string{
		"2026-10/A01.txt":       "90",
		"2026-10/A02.txt":       "100",
		"2026-11/A03.txt":       "85",
		"2026-11/draft/A04.txt": "120",
		"notes.md":              "月底對帳\n",
	}
	for name, body := range files {
		p := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// walkDemo 走訪整棵目錄樹，印出每個檔案的大小；名叫 draft 的目錄整個跳過。
func walkDemo() {
	dir, done := scratch()
	defer done()
	if err := makeTree(dir); err != nil {
		fmt.Println("建立失敗：", err)
		return
	}

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err // 這個目錄讀不了：回傳錯誤，整個走訪就停下來
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		if d.IsDir() {
			if d.Name() == "draft" {
				fmt.Println("跳過", rel+"/")
				return filepath.SkipDir
			}
			fmt.Println(rel + "/")
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		fmt.Println(rel, info.Size(), "byte")
		return nil
	})
	fmt.Println("錯誤：", err)
}

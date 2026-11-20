package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strconv"
	"strings"
)

// receiptTotal 加總 fsys 裡所有 .txt 發票的金額，名叫 draft 的目錄整個跳過。
// 參數是 fs.FS 而不是目錄路徑：真的目錄、embed 進來的檔案、測試用的 MapFS 都能傳進來。
// fs.FS 裡的路徑一律用 / 分隔，而且是相對路徑，所以用 path 套件，不用 filepath。
func receiptTotal(fsys fs.FS) (int, error) {
	total := 0
	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() == "draft" {
			return fs.SkipDir
		}
		if d.IsDir() || path.Ext(p) != ".txt" {
			return nil
		}

		data, err := fs.ReadFile(fsys, p)
		if err != nil {
			return err
		}
		n, err := strconv.Atoi(strings.TrimSpace(string(data)))
		if err != nil {
			return fmt.Errorf("%s：%w", p, err)
		}
		total += n
		return nil
	})
	return total, err
}

// fsDemo 把磁碟上的目錄包成 fs.FS，交給 receiptTotal。
func fsDemo() {
	dir, done := scratch()
	defer done()
	if err := makeTree(dir); err != nil {
		fmt.Println("建立失敗：", err)
		return
	}

	total, err := receiptTotal(os.DirFS(dir))
	fmt.Println("磁碟上的發票合計：", total, err)
}

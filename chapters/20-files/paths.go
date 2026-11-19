package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// receiptPath 依照月份和發票號碼，組出發票檔的位置。
func receiptPath(base, month, number string) string {
	return filepath.Join(base, month, number+".txt")
}

// receiptNumber 從發票檔的位置取回號碼：最後一段去掉副檔名。
func receiptNumber(p string) string {
	return strings.TrimSuffix(filepath.Base(p), filepath.Ext(p))
}

// pathsDemo 組路徑、拆路徑都交給 filepath，它知道這個作業系統用哪一種分隔符號。
// 印出來之前用 ToSlash 統一成 /，所以這段輸出在 Windows 上也一樣。
func pathsDemo() {
	p := receiptPath("receipts", "2026-11", "A03")
	fmt.Println(filepath.ToSlash(p), "號碼", receiptNumber(p))
	fmt.Println("同一個月的放在", filepath.ToSlash(filepath.Dir(p)))

	// 設定檔裡的目錄常常寫得不乾淨，Join 會順便整理好
	messy := receiptPath("./receipts/", "2026-11/../2026-10", "A01")
	fmt.Println(filepath.ToSlash(messy), "號碼", receiptNumber(messy))
}

// traversalDemo 檔名是使用者給的，直接 Join 就可能跑出 uploads 目錄。
func traversalDemo() {
	for _, name := range []string{"A01.jpg", "../config.json", "photos/../../../etc/passwd"} {
		p := filepath.Join("uploads", name)
		fmt.Printf("%s → %s，IsLocal=%v\n", name, filepath.ToSlash(p), filepath.IsLocal(name))
	}
}

// rootDemo 用 os.Root 把所有檔案操作關在一個目錄裡，名稱想跑出去就直接失敗。
func rootDemo() {
	dir, done := scratch()
	defer done()

	root, err := os.OpenRoot(dir)
	if err != nil {
		fmt.Println("開不了目錄：", err)
		return
	}
	defer root.Close()

	if err := root.WriteFile("A01.jpg", []byte("照片"), 0o644); err != nil {
		fmt.Println("寫入失敗：", err)
		return
	}
	data, err := root.ReadFile("A01.jpg")
	fmt.Println(string(data), err)

	for _, name := range []string{"../config.json", "/etc/passwd"} {
		_, err := root.Open(name)
		fmt.Println(err)
	}
}

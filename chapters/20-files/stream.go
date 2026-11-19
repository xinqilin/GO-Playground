package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// copyFile 把 src 複製成 dst。io.Copy 每次只搬一小塊，檔案再大也只用固定大小的記憶體。
// 先開來源再建目的檔：來源打不開時，不會留下一個空的目的檔。
func copyFile(dst, src string) (n int64, err error) {
	in, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer func() {
		err = errors.Join(err, out.Close())
	}()

	return io.Copy(out, in)
}

// copyFileDemo 備份菜單，再試一次來源不存在的情況。
func copyFileDemo() {
	dir, done := scratch()
	defer done()

	src := filepath.Join(dir, "menu.txt")
	if err := os.WriteFile(src, []byte("排骨飯 90\n雞腿飯 100\n"), 0o644); err != nil {
		fmt.Println("寫入失敗：", err)
		return
	}
	n, err := copyFile(filepath.Join(dir, "menu.bak"), src)
	fmt.Println("複製了", n, "byte", err)

	dst := filepath.Join(dir, "lost.bak")
	_, err = copyFile(dst, filepath.Join(dir, "lost.txt"))
	fmt.Println("來源不存在？", errors.Is(err, fs.ErrNotExist))
	_, err = os.Stat(dst)
	fmt.Println("目的檔不存在？", errors.Is(err, fs.ErrNotExist))
}

// copyWithSum 把 src 全部複製到 dst，同時算出內容的 SHA-256。
// TeeReader 讓讀出來的每一塊也順便寫進 h，整份資料只要走過一次。
func copyWithSum(dst io.Writer, src io.Reader) (int64, string, error) {
	h := sha256.New()
	n, err := io.Copy(dst, io.TeeReader(src, h))
	return n, hex.EncodeToString(h.Sum(nil)), err
}

// mergeDemo 把三天的營業額接成一份，同時算出合併後內容的指紋。
// MultiReader 讀完第一個才讀下一個，對外看起來就是一個 Reader。
func mergeDemo() {
	days := io.MultiReader(
		strings.NewReader("週一 12800\n"),
		strings.NewReader("週二 9650\n"),
		strings.NewReader("週三 15020\n"),
	)

	var merged strings.Builder
	n, sum, err := copyWithSum(&merged, days)
	fmt.Print(merged.String())
	fmt.Println(n, "byte，SHA-256 開頭", sum[:12], err)
}

// previewDemo 只看最前面 n 個 byte。LimitReader 數的是 byte 不是字，
// 切在中文字中間，就會拿到半個字。
func previewDemo() {
	const calls = "A01 排骨飯\nA02 雞腿飯\n"
	for _, n := range []int64{14, 12} {
		b, err := io.ReadAll(io.LimitReader(strings.NewReader(calls), n))
		fmt.Printf("前 %d byte：%q %v\n", n, b, err)
	}
}

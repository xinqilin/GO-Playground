package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// scratch 開一個新的暫存目錄給示範用，回傳的 done 會把它整個刪掉。
// MkdirTemp 的第一個參數是空字串時，用系統的暫存目錄（os.TempDir）；名稱裡的 * 會換成亂數。
func scratch() (dir string, done func()) {
	dir, err := os.MkdirTemp("", "ch20-*")
	if err != nil {
		panic(err) // 連暫存目錄都開不了，示範沒辦法繼續
	}
	return dir, func() { os.RemoveAll(dir) }
}

// saveAtomic 用 write 產生新內容，整份寫好才換掉 path。
// 先寫進同一個目錄裡的暫存檔，成功才 Rename 過去；中途出錯就刪掉暫存檔，原本的檔案完全沒動到。
func saveAtomic(path string, write func(io.Writer) error) (err error) {
	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tmp.Close()
			os.Remove(tmp.Name())
		}
	}()

	if err = write(tmp); err != nil {
		return err
	}
	if err = tmp.Sync(); err != nil { // 確定資料真的寫到磁碟上，再換名字
		return err
	}
	if err = tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

// atomicDemo 存兩次價目表：第一次成功，第二次寫到一半出錯。
func atomicDemo() {
	dir, done := scratch()
	defer done()
	path := filepath.Join(dir, "prices.txt")

	err := saveAtomic(path, func(w io.Writer) error {
		_, err := io.WriteString(w, "排骨飯 90\n雞腿飯 100\n")
		return err
	})
	fmt.Println("第一次存檔：", err)

	err = saveAtomic(path, func(w io.Writer) error {
		if _, err := io.WriteString(w, "排骨飯 95\n"); err != nil {
			return err
		}
		return errors.New("第二行的價格不是數字")
	})
	fmt.Println("第二次存檔：", err)

	data, err := os.ReadFile(path)
	fmt.Printf("檔案內容：%q %v\n", data, err)

	entries, err := os.ReadDir(dir)
	for _, e := range entries {
		fmt.Println("目錄裡有：", e.Name())
	}
	fmt.Println("錯誤：", err)
}

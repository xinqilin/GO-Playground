package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// wholeDemo 把整份菜單一次寫進檔案，再一次讀回來。
func wholeDemo() {
	dir, done := scratch()
	defer done()

	path := filepath.Join(dir, "menu.txt")
	err := os.WriteFile(path, []byte("排骨飯 90\n雞腿飯 100\n"), 0o644)
	if err != nil {
		fmt.Println("寫入失敗：", err)
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("讀取失敗：", err)
		return
	}
	fmt.Printf("讀回 %d byte：\n%s", len(data), data)
}

// loadNotice 讀今天的公告；檔案不存在不算錯誤，代表今天沒有公告。
func loadNotice(path string) (string, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return "（今天沒有公告）", nil
	}
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// missingDemo 看看檔案不存在時，拿到的錯誤長什麼樣子。
func missingDemo() {
	dir, done := scratch()
	defer done()

	path := filepath.Join(dir, "notice.txt")
	_, err := os.ReadFile(path)
	fmt.Println("不存在？", errors.Is(err, fs.ErrNotExist))

	if pe, ok := errors.AsType[*fs.PathError](err); ok {
		fmt.Printf("Op=%s 檔名=%s Err=%v\n", pe.Op, filepath.Base(pe.Path), pe.Err)
	}

	notice, err := loadNotice(path)
	fmt.Println(notice, err)
}

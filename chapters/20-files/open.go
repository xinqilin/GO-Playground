package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// appendCall 把一筆叫號紀錄加在 path 的最後面；檔案不存在就建立。
// 寫入的檔案要檢查 Close 的錯誤，所以回傳值要有名字，defer 才改得到它。
func appendCall(path, line string) (err error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	_, err = fmt.Fprintln(f, line)
	return err
}

// countCalls 一行一行讀 path，數有幾筆叫號（空白行不算）。
// 只讀不寫的檔案，Close 失敗也不會讓資料不見，defer f.Close() 就夠了。
func countCalls(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	n := 0
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if strings.TrimSpace(sc.Text()) != "" {
			n++
		}
	}
	return n, sc.Err()
}

// openDemo 一筆一筆附加叫號紀錄，再讀回來數一數。
func openDemo() {
	dir, done := scratch()
	defer done()

	path := filepath.Join(dir, "calls.log")
	for _, line := range []string{"A01 排骨飯", "A02 雞腿飯", "", "A03 魚排飯"} {
		if err := appendCall(path, line); err != nil {
			fmt.Println("寫入失敗：", err)
			return
		}
	}

	n, err := countCalls(path)
	fmt.Println("叫了", n, "號", err)
}

package main

import (
	"errors"
	"fmt"
	"io"
)

// receiptFile 是假的檔案：寫入一定成功，Close 時回傳 closeErr。
type receiptFile struct {
	closeErr error
}

func (f *receiptFile) Write(p []byte) (int, error) {
	return len(p), nil
}

func (f *receiptFile) Close() error {
	return f.closeErr
}

// saveReceipt 寫完之後關檔。回傳值有名字，defer 才能把 Close 的錯誤併進去。
func saveReceipt(w io.WriteCloser, total int) (err error) {
	defer func() {
		err = errors.Join(err, w.Close())
	}()
	_, err = fmt.Fprintf(w, "合計 %d 元\n", total)
	return err
}

func closeDemo() {
	fmt.Println(saveReceipt(&receiptFile{}, 290))
	fmt.Println(saveReceipt(&receiptFile{closeErr: errors.New("磁碟滿了，收據沒有存進去")}, 290))
}

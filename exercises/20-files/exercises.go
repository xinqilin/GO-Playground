//go:build !solution

package cabinet

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
)

// LoadOrInit 讀 path 的內容。檔案還不存在時，用 initial 建立它（權限 0o644），回傳 initial。
// 其他錯誤（例如 path 是一個目錄）照樣回傳，也不要建立檔案。
func LoadOrInit(path, initial string) (string, error) {
	// TODO: 實作
	return "", nil
}

// Merge 依照順序把 srcs 的內容接起來寫進 dst（dst 已經存在就覆蓋），回傳寫了多少 byte。
// 先把所有來源都打開：任何一個打不開，就回傳那個錯誤，而且不要建立 dst。
// 每個打開的檔案都要關掉，dst 的 Close 錯誤也要回傳。
func Merge(dst string, srcs ...string) (int64, error) {
	// TODO: 實作
	return 0, nil
}

// Usage 算出 fsys 裡每個最上層目錄底下所有檔案的大小總和（包含更深的子目錄）。
// 直接放在最上層的檔案，算在 "." 底下；沒有任何檔案的目錄不出現在結果裡。
func Usage(fsys fs.FS) (map[string]int64, error) {
	// TODO: 實作
	return nil, nil
}

// Rotate 輪替日誌檔：path.(keep-1) 改名成 path.keep，…，path.1 改名成 path.2，
// 最後 path 改名成 path.1；原本的 path.keep 被蓋掉。
// 中間缺了某個編號就跳過它；path 本身不存在時什麼都不做，回傳 nil。keep 至少是 1。
func Rotate(path string, keep int) error {
	// TODO: 實作
	return nil
}

// SaveUpload 把 r 的內容存成 dir 底下的 name，最多 limit 個 byte。
//   - name 想跑出 dir（例如 "../x"、"/etc/passwd"、指向外面的符號連結）時回傳錯誤。
//   - name 已經存在時不覆蓋，回傳的錯誤要讓 errors.Is(err, fs.ErrExist) 成立。
//   - 內容超過 limit 時回傳 ErrTooLarge，而且不留下寫到一半的檔案。
func SaveUpload(dir, name string, r io.Reader, limit int64) error {
	// TODO: 實作
	return nil
}

// WriteReport 把每一行加上編號寫進 path，例如 "1. 排骨飯 12"。
//
// 找錯：檔案建立了，內容卻是空的；Close 的錯誤也被丟掉了。
func WriteReport(path string, rows []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for i, row := range rows {
		if _, err := fmt.Fprintf(w, "%d. %s\n", i+1, row); err != nil {
			return err
		}
	}
	return nil
}

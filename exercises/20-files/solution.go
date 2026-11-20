//go:build solution

package cabinet

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

// LoadOrInit 只有在 errors.Is(err, fs.ErrNotExist) 時才建立檔案，其他錯誤原樣回傳。
func LoadOrInit(path, initial string) (string, error) {
	data, err := os.ReadFile(path)
	if err == nil {
		return string(data), nil
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return "", err
	}
	if err := os.WriteFile(path, []byte(initial), 0o644); err != nil {
		return "", err
	}
	return initial, nil
}

// Merge 先打開所有來源，再建立 dst，用 MultiReader 把來源接成一個 Reader 一次 Copy 完。
func Merge(dst string, srcs ...string) (n int64, err error) {
	readers := make([]io.Reader, 0, len(srcs))
	for _, src := range srcs {
		f, err := os.Open(src)
		if err != nil {
			return 0, err
		}
		defer f.Close() // 每個來源都要開著直到 Copy 結束，所以這裡在迴圈裡 defer 是對的
		readers = append(readers, f)
	}

	out, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer func() {
		err = errors.Join(err, out.Close())
	}()

	return io.Copy(out, io.MultiReader(readers...))
}

// Usage 走訪整個 fsys，用路徑的第一段當 key；路徑裡沒有 / 的檔案就在最上層。
func Usage(fsys fs.FS) (map[string]int64, error) {
	usage := map[string]int64{}
	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		top, _, nested := strings.Cut(p, "/")
		if !nested {
			top = "."
		}
		usage[top] += info.Size()
		return nil
	})
	if err != nil {
		return nil, err
	}
	return usage, nil
}

// Rotate 從編號最大的開始往後挪，才不會把還沒挪的檔案蓋掉。
func Rotate(path string, keep int) error {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}
	for i := keep - 1; i >= 1; i-- {
		err := os.Rename(fmt.Sprintf("%s.%d", path, i), fmt.Sprintf("%s.%d", path, i+1))
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	return os.Rename(path, path+".1")
}

// SaveUpload 用 os.Root 擋住跑出 dir 的名稱，O_EXCL 擋住已經存在的檔案，
// 多讀 1 個 byte 來判斷內容有沒有超過 limit。
func SaveUpload(dir, name string, r io.Reader, limit int64) (err error) {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return err
	}
	defer root.Close()

	f, err := root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, f.Close())
		if err != nil {
			root.Remove(name)
		}
	}()

	n, err := io.Copy(f, io.LimitReader(r, limit+1))
	if err != nil {
		return err
	}
	if n > limit {
		return ErrTooLarge
	}
	return nil
}

// WriteReport 寫完要 Flush，把 bufio.Writer 裡還沒寫出去的資料交給檔案；
// Close 的錯誤用 errors.Join 併進回傳值。
func WriteReport(path string, rows []string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	w := bufio.NewWriter(f)
	for i, row := range rows {
		if _, err := fmt.Fprintf(w, "%d. %s\n", i+1, row); err != nil {
			return err
		}
	}
	return w.Flush()
}

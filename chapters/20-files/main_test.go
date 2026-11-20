package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func Example_wholeDemo() {
	wholeDemo()
	// Output:
	// 讀回 27 byte：
	// 排骨飯 90
	// 雞腿飯 100
}

func Example_missingDemo() {
	missingDemo()
	// Output:
	// 不存在？ true
	// Op=open 檔名=notice.txt Err=no such file or directory
	// （今天沒有公告） <nil>
}

func Example_openDemo() {
	openDemo()
	// Output:
	// 叫了 3 號 <nil>
}

func Example_copyFileDemo() {
	copyFileDemo()
	// Output:
	// 複製了 27 byte <nil>
	// 來源不存在？ true
	// 目的檔不存在？ true
}

func Example_mergeDemo() {
	mergeDemo()
	// Output:
	// 週一 12800
	// 週二 9650
	// 週三 15020
	// 38 byte，SHA-256 開頭 36699ab88405 <nil>
}

func Example_previewDemo() {
	previewDemo()
	// Output:
	// 前 14 byte："A01 排骨飯\n" <nil>
	// 前 12 byte："A01 排骨\xe9\xa3" <nil>
}

func Example_pathsDemo() {
	pathsDemo()
	// Output:
	// receipts/2026-11/A03.txt 號碼 A03
	// 同一個月的放在 receipts/2026-11
	// receipts/2026-10/A01.txt 號碼 A01
}

func Example_traversalDemo() {
	traversalDemo()
	// Output:
	// A01.jpg → uploads/A01.jpg，IsLocal=true
	// ../config.json → config.json，IsLocal=false
	// photos/../../../etc/passwd → ../etc/passwd，IsLocal=false
}

func Example_rootDemo() {
	rootDemo()
	// Output:
	// 照片 <nil>
	// openat ../config.json: path escapes from parent
	// openat /etc/passwd: path escapes from parent
}

func Example_walkDemo() {
	walkDemo()
	// Output:
	// ./
	// 2026-10/
	// 2026-10/A01.txt 2 byte
	// 2026-10/A02.txt 3 byte
	// 2026-11/
	// 2026-11/A03.txt 2 byte
	// 跳過 2026-11/draft/
	// notes.md 13 byte
	// 錯誤： <nil>
}

func Example_fsDemo() {
	fsDemo()
	// Output:
	// 磁碟上的發票合計： 275 <nil>
}

func Example_embedDemo() {
	embedDemo()
	// Output:
	// 歡迎光臨，今日公休：週一
	// embed 進來的發票合計： 205 <nil>
	// [receipts/2026-12/B01.txt receipts/2027-01/B02.txt receipts/draft/B03.txt] <nil>
}

func Example_atomicDemo() {
	atomicDemo()
	// Output:
	// 第一次存檔： <nil>
	// 第二次存檔： 第二行的價格不是數字
	// 檔案內容："排骨飯 90\n雞腿飯 100\n" <nil>
	// 目錄裡有： prices.txt
	// 錯誤： <nil>
}

// TestReceiptTotal 用記憶體裡的 MapFS 當檔案系統，不必真的建目錄和檔案。
func TestReceiptTotal(t *testing.T) {
	fsys := fstest.MapFS{
		"2026-10/A01.txt":       {Data: []byte("90\n")},
		"2026-10/notes.md":      {Data: []byte("不是發票")},
		"2026-11/A02.txt":       {Data: []byte("100")},
		"2026-11/draft/A03.txt": {Data: []byte("999")},
	}
	if got, err := receiptTotal(fsys); got != 190 || err != nil {
		t.Errorf("receiptTotal = %d, %v；want 190, nil", got, err)
	}

	fsys["2026-12/A04.txt"] = &fstest.MapFile{Data: []byte("九十")}
	if _, err := receiptTotal(fsys); err == nil || !strings.Contains(err.Error(), "2026-12/A04.txt") {
		t.Errorf("金額不是數字時 err = %v；want 錯誤訊息裡有檔名", err)
	}
}

// TestRootSymlink：目錄裡的符號連結指向外面，直接開打得開，透過 os.Root 就打不開。
func TestRootSymlink(t *testing.T) {
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "secret.txt"), []byte("密碼"), 0o644); err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(dir, "link")); err != nil {
		t.Skip("這個系統不能建立符號連結：", err)
	}

	if _, err := os.ReadFile(filepath.Join(dir, "link", "secret.txt")); err != nil {
		t.Fatalf("直接開應該打得開：%v", err)
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	if _, err := root.ReadFile("link/secret.txt"); err == nil {
		t.Error("透過 os.Root 不應該讀得到目錄外面的檔案")
	}
}

// TestSaveAtomicReplaces：舊檔被整份換掉，目錄裡不會留下暫存檔。
func TestSaveAtomicReplaces(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "prices.txt")
	if err := os.WriteFile(path, []byte("舊的"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := saveAtomic(path, func(w io.Writer) error {
		_, err := io.WriteString(w, "新的")
		return err
	})
	if data, _ := os.ReadFile(path); err != nil || string(data) != "新的" {
		t.Errorf("內容 %q, err %v；want \"新的\", nil", data, err)
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 1 {
		t.Errorf("目錄裡有 %d 個檔案，want 1", len(entries))
	}

	if err := saveAtomic(filepath.Join(dir, "nodir", "x.txt"), func(io.Writer) error { return nil }); err == nil {
		t.Error("目錄不存在時應該回傳錯誤")
	}
}

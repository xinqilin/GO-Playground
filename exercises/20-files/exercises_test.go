package cabinet

import (
	"errors"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

// write 在 dir 底下建立 name，內容是 body。
func write(t *testing.T, dir, name, body string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

// read 讀回 p 的內容；檔案不存在時回傳 "<不存在>"。
func read(t *testing.T, p string) string {
	t.Helper()
	data, err := os.ReadFile(p)
	if errors.Is(err, fs.ErrNotExist) {
		return "<不存在>"
	}
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestLoadOrInit(t *testing.T) {
	dir := t.TempDir()

	p := filepath.Join(dir, "prices.txt")
	got, err := LoadOrInit(p, "排骨飯 90\n")
	if got != "排骨飯 90\n" || err != nil || read(t, p) != "排骨飯 90\n" {
		t.Errorf("檔案不存在：回傳 %q, %v，檔案內容 %q；want 建立檔案並回傳 initial", got, err, read(t, p))
	}

	write(t, dir, "prices.txt", "雞腿飯 100\n")
	got, err = LoadOrInit(p, "排骨飯 90\n")
	if got != "雞腿飯 100\n" || err != nil || read(t, p) != "雞腿飯 100\n" {
		t.Errorf("檔案存在：回傳 %q, %v，檔案內容 %q；want 原本的內容、檔案不變", got, err, read(t, p))
	}

	if _, err := LoadOrInit(dir, "x"); err == nil {
		t.Error("path 是目錄時應該回傳錯誤")
	}
	if _, err := LoadOrInit(filepath.Join(dir, "nodir", "x.txt"), "x"); err == nil {
		t.Error("上層目錄不存在時應該回傳錯誤")
	}
}

func TestMerge(t *testing.T) {
	dir := t.TempDir()
	a := write(t, dir, "mon.txt", "週一 12800\n")
	b := write(t, dir, "tue.txt", "週二 9650\n")
	c := write(t, dir, "wed.txt", "")
	dst := filepath.Join(dir, "week.txt")
	write(t, dir, "week.txt", "舊的內容會被蓋掉，而且比新的長很多很多很多")

	n, err := Merge(dst, a, b, c)
	want := "週一 12800\n週二 9650\n"
	if got := read(t, dst); got != want || n != int64(len(want)) || err != nil {
		t.Errorf("Merge = %d, %v，內容 %q；want %d, nil，%q", n, err, got, len(want), want)
	}

	lost := filepath.Join(dir, "lost.txt")
	_, err = Merge(lost, a, filepath.Join(dir, "missing.txt"))
	if !errors.Is(err, fs.ErrNotExist) || read(t, lost) != "<不存在>" {
		t.Errorf("來源不存在：err = %v，dst %q；want fs.ErrNotExist 而且不建立 dst", err, read(t, lost))
	}
}

func TestUsage(t *testing.T) {
	fsys := fstest.MapFS{
		"2026-10/A01.txt":     {Data: []byte("90")},
		"2026-10/old/A00.txt": {Data: []byte("12345")},
		"2026-11/A03.txt":     {Data: []byte("100")},
		"readme.md":           {Data: []byte("對帳")},
		"empty":               {Mode: fs.ModeDir},
	}
	got, err := Usage(fsys)
	want := map[string]int64{"2026-10": 7, "2026-11": 3, ".": 6}
	if !maps.Equal(got, want) || err != nil {
		t.Errorf("Usage = %v, %v；want %v, nil", got, err, want)
	}
}

func TestRotate(t *testing.T) {
	dir := t.TempDir()
	log := filepath.Join(dir, "calls.log")
	names := func() string {
		var out []string
		for _, s := range []string{"", ".1", ".2", ".3", ".4"} {
			out = append(out, s+"="+read(t, log+s))
		}
		return strings.Join(out, " ")
	}

	write(t, dir, "calls.log", "c")
	write(t, dir, "calls.log.1", "b")
	write(t, dir, "calls.log.2", "a")
	if err := Rotate(log, 2); err != nil {
		t.Fatal(err)
	}
	if got, want := names(), "=<不存在> .1=c .2=b .3=<不存在> .4=<不存在>"; got != want {
		t.Errorf("keep=2：%s\nwant %s", got, want)
	}

	write(t, dir, "calls.log", "d")
	if err := Rotate(log, 4); err != nil {
		t.Fatal(err)
	}
	if got, want := names(), "=<不存在> .1=d .2=c .3=b .4=<不存在>"; got != want {
		t.Errorf("keep=4：%s\nwant %s", got, want)
	}

	if err := Rotate(log, 4); err != nil {
		t.Errorf("path 不存在時應該回傳 nil，got %v", err)
	}
	if got, want := names(), "=<不存在> .1=d .2=c .3=b .4=<不存在>"; got != want {
		t.Errorf("path 不存在時什麼都不該動：%s\nwant %s", got, want)
	}
}

func TestSaveUpload(t *testing.T) {
	dir := t.TempDir()
	uploads := filepath.Join(dir, "uploads")
	if err := os.Mkdir(uploads, 0o755); err != nil {
		t.Fatal(err)
	}

	err := SaveUpload(uploads, "A01.jpg", strings.NewReader("照片"), 6)
	if got := read(t, filepath.Join(uploads, "A01.jpg")); err != nil || got != "照片" {
		t.Errorf("剛好 limit：err = %v，內容 %q；want nil，\"照片\"", err, got)
	}

	err = SaveUpload(uploads, "A01.jpg", strings.NewReader("另一張"), 100)
	if got := read(t, filepath.Join(uploads, "A01.jpg")); !errors.Is(err, fs.ErrExist) || got != "照片" {
		t.Errorf("已經存在：err = %v，內容 %q；want fs.ErrExist，原本的內容", err, got)
	}

	err = SaveUpload(uploads, "big.jpg", strings.NewReader("1234567"), 6)
	if got := read(t, filepath.Join(uploads, "big.jpg")); !errors.Is(err, ErrTooLarge) || got != "<不存在>" {
		t.Errorf("超過 limit：err = %v，檔案 %q；want ErrTooLarge，不留下檔案", err, got)
	}

	for _, name := range []string{"../evil.txt", "/evil.txt"} {
		err := SaveUpload(uploads, name, strings.NewReader("x"), 10)
		if got := read(t, filepath.Join(dir, "evil.txt")); err == nil || got != "<不存在>" {
			t.Errorf("SaveUpload(%q)：err = %v，目錄外的檔案 %q；want 錯誤，而且不建立檔案", name, err, got)
		}
	}

	if err := os.Symlink(dir, filepath.Join(uploads, "link")); err == nil {
		err := SaveUpload(uploads, "link/evil.txt", strings.NewReader("x"), 10)
		if got := read(t, filepath.Join(dir, "evil.txt")); err == nil || got != "<不存在>" {
			t.Errorf("透過符號連結：err = %v，目錄外的檔案 %q；want 錯誤，而且不建立檔案", err, got)
		}
	}
}

func TestWriteReport(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "report.txt")
	err := WriteReport(p, []string{"排骨飯 12", "雞腿飯 8"})
	if got, want := read(t, p), "1. 排骨飯 12\n2. 雞腿飯 8\n"; got != want || err != nil {
		t.Errorf("WriteReport：err = %v，內容 %q；want nil，%q", err, got, want)
	}
}

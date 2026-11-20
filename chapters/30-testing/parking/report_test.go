package parking

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"
)

// update 是這個 package 的測試自己加的旗標：go test ./chapters/30-testing/parking -update
var update = flag.Bool("update", false, "用目前的輸出覆寫 testdata 裡的 golden file")

// reportRecords 是報表測試用的固定資料。
var reportRecords = []Record{
	{"XY-5678", at(19, 22, 30), at(20, 8, 10), 300},
	{"ABC-1234", at(20, 8, 5), at(20, 9, 40), 80},
	{"1234-AB", at(20, 12, 0), at(20, 12, 8), 0},
}

func TestWriteReport(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteReport(&buf, reportRecords); err != nil {
		t.Fatal(err)
	}
	golden := filepath.Join("testdata", "report.golden")
	if *update {
		if err := os.WriteFile(golden, buf.Bytes(), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("%v（第一次請加 -update 產生）", err)
	}
	if got := buf.String(); got != string(want) {
		t.Errorf("報表和 %s 不一樣：\n--- got\n%s--- want\n%s", golden, got, want)
	}
}

// TestFeesInReport 確認報表用的資料本身是對的，golden file 才不會把錯的金額鎖住。
func TestFeesInReport(t *testing.T) {
	for _, r := range reportRecords {
		if fee, _ := Fee(r.In, r.Out); fee != r.Fee {
			t.Errorf("%s 的費用寫成 %d，Fee 算出來是 %d", r.Plate, r.Fee, fee)
		}
	}
}

func TestSaveReport(t *testing.T) {
	dir := t.TempDir() // 測試結束後自動刪除
	path, err := SaveReport(dir, reportRecords)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(dir, "report-2026-11-20.txt"); path != want {
		t.Errorf("path = %q，want %q", path, want)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var want bytes.Buffer
	WriteReport(&want, reportRecords)
	if !bytes.Equal(got, want.Bytes()) {
		t.Errorf("存檔的內容和 WriteReport 不一樣：\n%s", got)
	}
}

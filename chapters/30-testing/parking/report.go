package parking

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// WriteReport 把出場紀錄寫成日報表，日期取第一筆的出場日。
func WriteReport(w io.Writer, recs []Record) error {
	if len(recs) == 0 {
		_, err := fmt.Fprintln(w, "沒有出場紀錄")
		return err
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%s 出場紀錄\n", recs[0].Out.Format("2006-01-02"))
	total := 0
	for _, r := range recs {
		fmt.Fprintf(&b, "%-8s  %s → %s  %3d 元\n",
			r.Plate, r.In.Format("01-02 15:04"), r.Out.Format("15:04"), r.Fee)
		total += r.Fee
	}
	fmt.Fprintf(&b, "共 %d 台，合計 %d 元\n", len(recs), total)
	_, err := io.WriteString(w, b.String())
	return err
}

// SaveReport 把日報表存成 dir 裡的 report-日期.txt，回傳檔案路徑。
func SaveReport(dir string, recs []Record) (string, error) {
	name := "report.txt"
	if len(recs) > 0 {
		name = "report-" + recs[0].Out.Format("2006-01-02") + ".txt"
	}
	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	if err := WriteReport(f, recs); err != nil {
		f.Close()
		return "", err
	}
	return path, f.Close()
}

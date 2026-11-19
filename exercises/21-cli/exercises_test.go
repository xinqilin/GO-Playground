package kiosk

import (
	"context"
	"errors"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// quiet 執行 f，回傳 f 執行期間寫到 os.Stderr 的內容。
func quiet(t *testing.T, f func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stderr
	os.Stderr = w
	f()
	os.Stderr = old
	w.Close()
	out, _ := io.ReadAll(r)
	return string(out)
}

func TestParseBooking(t *testing.T) {
	ok := []struct {
		args []string
		want Booking
	}{
		{[]string{"-name", "王小姐"}, Booking{Name: "王小姐", People: 2, Time: "18:00"}},
		{[]string{"-name=陳先生", "-people", "5", "-time", "19:30", "-outdoor", "要", "兒童椅"},
			Booking{Name: "陳先生", People: 5, Time: "19:30", Outdoor: true, Note: "要 兒童椅"}},
	}
	for _, tt := range ok {
		got, err := ParseBooking(tt.args)
		if got != tt.want || err != nil {
			t.Errorf("ParseBooking(%q) = %+v, %v；want %+v, nil", tt.args, got, err, tt.want)
		}
	}

	bad := [][]string{
		{},
		{"-people", "3"},
		{"-name", "林", "-people", "0"},
		{"-name", "林", "-people", "13"},
		{"-name", "林", "-people", "三"},
		{"-name", "林", "-time", "晚上七點"},
		{"-name", "林", "-vip"},
	}
	printed := quiet(t, func() {
		for _, args := range bad {
			if _, err := ParseBooking(args); err == nil {
				t.Errorf("ParseBooking(%q) 應該回傳錯誤", args)
			}
		}
		if _, err := ParseBooking([]string{"-h"}); !errors.Is(err, flag.ErrHelp) {
			t.Errorf("ParseBooking(-h) 的錯誤 = %v；want flag.ErrHelp", err)
		}
	})
	if printed != "" {
		t.Errorf("ParseBooking 不應該印東西到 os.Stderr，卻印了：\n%s", printed)
	}
}

func TestLoadConfig(t *testing.T) {
	env := func(kv ...string) func(string) string {
		m := map[string]string{}
		for i := 0; i+1 < len(kv); i += 2 {
			m[kv[i]] = kv[i+1]
		}
		return func(k string) string { return m[k] }
	}
	tests := []struct {
		args   []string
		getenv func(string) string
		want   Config
	}{
		{nil, env(), Config{8080, "總店", false}},
		{nil, env("KIOSK_PORT", "9000", "KIOSK_SHOP", "西門店", "KIOSK_DEBUG", "true"), Config{9000, "西門店", true}},
		{[]string{"-port", "7000", "-debug=false"}, env("KIOSK_PORT", "9000", "KIOSK_DEBUG", "1"), Config{7000, "總店", false}},
		{[]string{"-shop", "板橋店"}, env("KIOSK_SHOP", "西門店"), Config{8080, "板橋店", false}},
	}
	for _, tt := range tests {
		got, err := LoadConfig(tt.args, tt.getenv)
		if got != tt.want || err != nil {
			t.Errorf("LoadConfig(%q) = %+v, %v；want %+v, nil", tt.args, got, err, tt.want)
		}
	}

	printed := quiet(t, func() {
		for _, key := range []string{"KIOSK_PORT", "KIOSK_DEBUG"} {
			_, err := LoadConfig(nil, env(key, "很多"))
			if err == nil || !strings.Contains(err.Error(), key) {
				t.Errorf("%s 看不懂時 err = %v；want 錯誤訊息提到 %s", key, err, key)
			}
		}
		if _, err := LoadConfig([]string{"-port", "八"}, env()); err == nil {
			t.Error("-port 八 應該回傳錯誤")
		}
	})
	if printed != "" {
		t.Errorf("LoadConfig 不應該印東西到 os.Stderr，卻印了：\n%s", printed)
	}
}

// cancelAt 是一個 io.Writer：收到第 n 行時，用 cause 取消 ctx，好像使用者剛好按了 Ctrl-C。
type cancelAt struct {
	n      int
	cancel context.CancelCauseFunc
	cause  error
	sb     strings.Builder
}

func (c *cancelAt) Write(p []byte) (int, error) {
	c.sb.Write(p)
	if strings.Count(c.sb.String(), "\n") == c.n {
		c.cancel(c.cause)
	}
	return len(p), nil
}

func TestRelay(t *testing.T) {
	var sb strings.Builder
	n, err := Relay(context.Background(), strings.NewReader("A01\nA02\nA03\n"), &sb)
	if n != 3 || err != nil || sb.String() != "A01\nA02\nA03\n" {
		t.Errorf("沒有取消：%d, %v，寫出 %q；want 3, nil，三行", n, err, sb.String())
	}

	ctx, cancel := context.WithCancelCause(context.Background())
	stopped := errors.New("店長按了 Ctrl-C")
	w := &cancelAt{n: 2, cancel: cancel, cause: stopped}
	n, err = Relay(ctx, strings.NewReader("A01\nA02\nA03\nA04\n"), w)
	if n != 2 || !errors.Is(err, stopped) || w.sb.String() != "A01\nA02\n" {
		t.Errorf("第 2 行之後取消：%d, %v，寫出 %q；want 2、取消的原因、只有兩行", n, err, w.sb.String())
	}
}

func TestCount(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.txt")
	b := filepath.Join(dir, "b.txt")
	if err := os.WriteFile(a, []byte("排骨飯 90\n雞腿飯 100\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(b, []byte("滷肉飯 60\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(dir, "nope.txt")

	tests := []struct {
		args      []string
		stdin     string
		code      int
		stdout    string
		stderrHas string
	}{
		{nil, "一\n二\n三\n", 0, "3 -\n", ""},
		{[]string{"-w"}, "排骨飯 90\n雞腿飯 100\n", 0, "4 -\n", ""},
		{[]string{a, "-"}, "x\n", 0, "2 " + a + "\n1 -\n3 總計\n", ""},
		{[]string{"-w", a, missing, b}, "", 1, "4 " + a + "\n2 " + b + "\n6 總計\n", "nope.txt"},
		{[]string{"-x"}, "", 2, "", "-x"},
		{[]string{"-h"}, "", 0, "", "-w"},
	}
	for _, tt := range tests {
		var stdout, stderr strings.Builder
		code := Count(tt.args, strings.NewReader(tt.stdin), &stdout, &stderr)
		if code != tt.code || stdout.String() != tt.stdout || !strings.Contains(stderr.String(), tt.stderrHas) ||
			(tt.stderrHas == "" && stderr.Len() > 0) {
			t.Errorf("Count(%q)：結束碼 %d，stdout %q，stderr %q；want %d，%q，stderr 包含 %q",
				tt.args, code, stdout.String(), stderr.String(), tt.code, tt.stdout, tt.stderrHas)
		}
	}
}

func TestReport(t *testing.T) {
	tests := []struct {
		args      []string
		code      int
		stdout    string
		stderrHas string
	}{
		{[]string{"-n", "2", "排骨飯", "雞腿飯", "滷肉飯"}, 0, "排骨飯\n雞腿飯\n", ""},
		{[]string{"排骨飯"}, 0, "排骨飯\n", ""},
		{[]string{"-n", "兩個", "排骨飯"}, 2, "", "invalid value"},
		{[]string{}, 2, "", "至少"},
		{[]string{"-h"}, 0, "", "-n"},
	}
	for _, tt := range tests {
		var stdout, stderr strings.Builder
		var code int
		printed := quiet(t, func() { code = Report(tt.args, &stdout, &stderr) })
		if code != tt.code || stdout.String() != tt.stdout || !strings.Contains(stderr.String(), tt.stderrHas) || printed != "" {
			t.Errorf("Report(%q)：結束碼 %d，stdout %q，stderr %q，os.Stderr %q；want %d，%q，stderr 包含 %q",
				tt.args, code, stdout.String(), stderr.String(), printed, tt.code, tt.stdout, tt.stderrHas)
		}
	}
}

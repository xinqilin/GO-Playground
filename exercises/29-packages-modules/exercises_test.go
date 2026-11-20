package modver

import (
	"maps"
	"slices"
	"testing"
)

func TestParse(t *testing.T) {
	good := []struct {
		in   string
		want Version
	}{
		{"v1.2.3", Version{1, 2, 3, ""}},
		{"v0.0.0", Version{0, 0, 0, ""}},
		{"v10.20.30", Version{10, 20, 30, ""}},
		{"v1.0.0-rc.1", Version{1, 0, 0, "rc.1"}},
		{"v0.0.0-20230129092748-24d4a6f8daec", Version{0, 0, 0, "20230129092748-24d4a6f8daec"}},
	}
	for _, tt := range good {
		got, err := Parse(tt.in)
		if err != nil || got != tt.want {
			t.Errorf("Parse(%q) = %+v, %v，want %+v", tt.in, got, err, tt.want)
		}
	}
	for _, bad := range []string{"", "1.2.3", "v1.2", "v1.2.3.4", "v1.02.3", "v1.2.x", "v1.2.3-", "v1.2.3-rc_1", "v-1.2.3", "v1.+2.3", "V1.2.3"} {
		if got, err := Parse(bad); err == nil {
			t.Errorf("Parse(%q) = %+v，want 錯誤", bad, got)
		}
	}
}

func v(major, minor, patch int, pre string) Version {
	return Version{major, minor, patch, pre}
}

func TestCompare(t *testing.T) {
	// 每一組的 a 都比 b 舊。
	older := [][2]Version{
		{v(1, 2, 3, ""), v(1, 2, 4, "")},
		{v(1, 9, 0, ""), v(1, 10, 0, "")},
		{v(1, 99, 99, ""), v(2, 0, 0, "")},
		{v(1, 0, 0, "rc.1"), v(1, 0, 0, "")},
		{v(1, 0, 0, "rc.2"), v(1, 0, 0, "rc.10")},
		{v(1, 0, 0, "alpha"), v(1, 0, 0, "beta")},
		{v(1, 0, 0, "alpha"), v(1, 0, 0, "alpha.1")},
		{v(1, 0, 0, "1"), v(1, 0, 0, "alpha")},
		{v(0, 9, 0, ""), v(1, 0, 0, "rc.1")},
	}
	for _, p := range older {
		a, b := p[0], p[1]
		if got := Compare(a, b); got != -1 {
			t.Errorf("Compare(%v, %v) = %d，want -1", a, b, got)
		}
		if got := Compare(b, a); got != 1 {
			t.Errorf("Compare(%v, %v) = %d，want 1", b, a, got)
		}
	}
	for _, x := range []Version{v(1, 2, 3, ""), v(1, 0, 0, "rc.1")} {
		if got := Compare(x, x); got != 0 {
			t.Errorf("Compare(%v, %v) = %d，want 0", x, x, got)
		}
	}
}

func TestLatest(t *testing.T) {
	available := []Version{
		v(1, 2, 0, ""), v(1, 4, 1, ""), v(1, 10, 0, ""), v(1, 11, 0, "rc.1"),
		v(2, 0, 0, ""), v(2, 1, 0, ""), v(0, 9, 0, ""),
	}
	tests := []struct {
		current Version
		want    Version
		wantOK  bool
	}{
		{v(1, 3, 0, ""), v(1, 10, 0, ""), true}, // 不會跳到 v2，也不會選 rc
		{v(2, 0, 0, ""), v(2, 1, 0, ""), true},
		{v(2, 1, 0, ""), v(2, 1, 0, ""), false},
		{v(0, 5, 0, ""), v(0, 9, 0, ""), true},
		{v(3, 0, 0, ""), v(3, 0, 0, ""), false},
	}
	for _, tt := range tests {
		got, ok := Latest(tt.current, available)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("Latest(%v) = %v, %v，want %v, %v", tt.current, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestParseRequires(t *testing.T) {
	const gomod = `module github.com/xinqilin/GO-Playground/chapters/28-database

go 1.27

require modernc.org/sqlite v1.62.0

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect

	// 註解行要略過
	modernc.org/libc v1.80.1
)

replace (
	example.com/rice => ../rice
)

require example.com/soy v0.3.0 // indirect
`
	got, err := ParseRequires(gomod)
	want := []Require{
		{"modernc.org/sqlite", "v1.62.0", false},
		{"github.com/google/uuid", "v1.6.0", true},
		{"github.com/remyoudompheng/bigfft", "v0.0.0-20230129092748-24d4a6f8daec", true},
		{"modernc.org/libc", "v1.80.1", false},
		{"example.com/soy", "v0.3.0", true},
	}
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("ParseRequires =\n  %+v, %v\nwant\n  %+v", got, err, want)
	}

	for _, bad := range []string{
		"require example.com/soy\n",
		"require (\n\texample.com/soy v0.3.0 extra\n)\n",
	} {
		if got, err := ParseRequires(bad); err == nil {
			t.Errorf("ParseRequires(%q) = %+v，want 錯誤", bad, got)
		}
	}
}

func TestSelect(t *testing.T) {
	reqs := map[string][]string{
		"main":                     {"example.com/rice@v1.2.0", "example.com/sauce@v1.2.0"},
		"example.com/rice@v1.2.0":  {"example.com/pork@v1.3.0"},
		"example.com/sauce@v1.2.0": {"example.com/pork@v1.4.0", "example.com/salt@v1.0.0"},
		"example.com/pork@v1.3.0":  {"example.com/salt@v1.1.0"}, // pork 最後選 v1.4.0，但這一版的 require 一樣要算
		"example.com/salt@v1.1.0":  {"example.com/soy@v1.0.0"},
		"example.com/salt@v1.0.0":  {"example.com/oil@v1.0.0"},
	}
	got, err := Select(reqs)
	want := map[string]string{
		"example.com/rice":  "v1.2.0",
		"example.com/sauce": "v1.2.0",
		"example.com/pork":  "v1.4.0",
		"example.com/salt":  "v1.1.0",
		"example.com/soy":   "v1.0.0",
		"example.com/oil":   "v1.0.0",
	}
	if err != nil || !maps.Equal(got, want) {
		t.Errorf("Select =\n  %v, %v\nwant\n  %v", got, err, want)
	}

	// 數字的大小，不是字串的大小：v1.10.0 比 v1.9.0 新。
	got, err = Select(map[string][]string{
		"main":                    {"example.com/rice@v1.9.0", "example.com/egg@v1.0.0"},
		"example.com/egg@v1.0.0":  {"example.com/rice@v1.10.0"},
		"example.com/rice@v1.9.0": {},
	})
	if err != nil || got["example.com/rice"] != "v1.10.0" {
		t.Errorf("rice 選了 %q（%v），want v1.10.0", got["example.com/rice"], err)
	}

	if got, err := Select(map[string][]string{"main": {"example.com/rice@latest"}}); err == nil {
		t.Errorf("版本不合法時應該回傳錯誤，拿到 %v", got)
	}
}

func TestInternalAllowed(t *testing.T) {
	tests := []struct {
		importer, imported string
		want               bool
	}{
		{"example.com/shop", "example.com/shop/internal/price", true},
		{"example.com/shop/cmd/server", "example.com/shop/internal/price", true},
		{"example.com/shop/internal/tax", "example.com/shop/internal/price", true},
		{"example.com/other", "example.com/shop/internal/price", false},
		{"example.com/shopping", "example.com/shop/internal/price", false},
		{"example.com/shop/counter", "example.com/shop/kitchen/internal/recipe", false},
		{"example.com/shop/kitchen", "example.com/shop/kitchen/internal/recipe", true},
		{"example.com/shop/cmd", "example.com/shop/internal/a/internal/b", false},
		{"example.com/shop/internal/a/x", "example.com/shop/internal/a/internal/b", true},
		{"example.com/other", "example.com/shop/internals/price", true},
		{"example.com/other", "example.com/shop/price", true},
	}
	for _, tt := range tests {
		if got := InternalAllowed(tt.importer, tt.imported); got != tt.want {
			t.Errorf("InternalAllowed(%q, %q) = %v，want %v", tt.importer, tt.imported, got, tt.want)
		}
	}
}

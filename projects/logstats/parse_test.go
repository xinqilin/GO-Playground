package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		line    string
		want    string // Entry 的摘要；錯誤時是錯誤訊息的開頭
		wantErr bool
	}{
		{`192.0.2.10 - - [21/Nov/2026:09:30:01 +0800] "GET /menu HTTP/1.1" 200 1532 0.012`, "09:30:01 GET /menu 200 12ms", false},
		{`192.0.2.10 - - [21/Nov/2026:09:30:01 +0800] "GET /menu?lang=en HTTP/2.0" 304 - 0.003`, "09:30:01 GET /menu?lang=en 304 3ms", false},
		{`{"time":"2026-11-21T09:35:10.120+08:00","msg":"request","method":"POST","path":"/todos","status":201,"took":2350000}`, "09:35:10 POST /todos 201 2.35ms", false},
		{`192.0.2.16 - - [21/Nov/2026:09:32:0`, "格式不對", true},
		{`192.0.2.17 - - [21/Nov/2026:09:32:50 +0800] "GET /menu HTTP/1.1" 2OO 1532 0.012`, "格式不對", true},
		{`192.0.2.18 - - [31/Nov/2026:09:33:40 +0800] "GET /menu HTTP/1.1" 200 1532 0.012`, "時間不對", true},
		{`{"time":"2026-11-21T09:35:10+08:00","msg":"request","status":"200"}`, "JSON 不對", true},
	}
	for _, tt := range tests {
		e, err := ParseLine(tt.line)
		got := fmt.Sprintf("%s %s %s %d %v", e.Time.Format("15:04:05"), e.Method, e.Path, e.Status, e.Took)
		if err != nil {
			got = err.Error()
		}
		if (err != nil) != tt.wantErr || !strings.HasPrefix(got, tt.want) {
			t.Errorf("ParseLine(%q)\n  = %s\n想要 %s", tt.line, got, tt.want)
		}
	}
}

func TestParseLineSkip(t *testing.T) {
	for _, line := range []string{"", "   ", `{"time":"2026-11-21T09:35:00+08:00","level":"INFO","msg":"開始服務"}`} {
		if _, err := ParseLine(line); !errors.Is(err, errSkip) {
			t.Errorf("ParseLine(%q) 的錯誤 = %v，想要 errSkip", line, err)
		}
	}
}

// parseAccessCut 是不用 regexp、只用 strings.Cut 的版本，給 benchmark 比較用。
func parseAccessCut(line string) (Entry, error) {
	bad := errors.New("格式不對")
	_, rest, ok1 := strings.Cut(line, " [")
	ts, rest, ok2 := strings.Cut(rest, `] "`)
	req, rest, ok3 := strings.Cut(rest, `" `)
	method, req, ok4 := strings.Cut(req, " ")
	path, _, ok5 := strings.Cut(req, " ")
	fields := strings.Fields(rest) // 狀態碼、大小、處理時間
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || len(fields) != 3 {
		return Entry{}, bad
	}
	t, err := time.Parse(accessLayout, ts)
	if err != nil {
		return Entry{}, err
	}
	status, err := strconv.Atoi(fields[0])
	if err != nil {
		return Entry{}, bad
	}
	took, err := time.ParseDuration(fields[2] + "s")
	if err != nil {
		return Entry{}, bad
	}
	return Entry{Time: t, Method: method, Path: path, Status: status, Took: took}, nil
}

// TestCutAgrees 確認兩個版本對範例檔裡每一行合法的 log 結果都相同，benchmark 比的才是同一件事。
func TestCutAgrees(t *testing.T) {
	data, err := os.ReadFile("testdata/access.log")
	if err != nil {
		t.Fatal(err)
	}
	for i, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		want, err := parseAccess(line)
		if err != nil {
			continue
		}
		if got, err := parseAccessCut(line); err != nil || !got.Time.Equal(want.Time) ||
			got.Method != want.Method || got.Path != want.Path || got.Status != want.Status || got.Took != want.Took {
			t.Errorf("第 %d 行：regexp 得到 %+v，Cut 得到 %+v（%v）", i+1, want, got, err)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	const access = `198.51.100.4 - - [21/Nov/2026:09:30:12 +0800] "POST /orders HTTP/1.1" 201 88 0.045`
	const jsonLine = `{"time":"2026-11-21T09:35:10.120+08:00","level":"INFO","msg":"request","method":"POST","path":"/todos","status":201,"took":2350000}`
	b.Run("regexp", func(b *testing.B) {
		for b.Loop() {
			parseAccess(access)
		}
	})
	b.Run("cut", func(b *testing.B) {
		for b.Loop() {
			parseAccessCut(access)
		}
	})
	b.Run("json", func(b *testing.B) {
		for b.Loop() {
			parseJSON(jsonLine)
		}
	})
}

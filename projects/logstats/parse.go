package main

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Entry 是一個 HTTP 請求的紀錄。
type Entry struct {
	Time   time.Time
	Method string
	Path   string
	Status int
	Took   time.Duration
}

// errSkip 表示這一行不是請求的紀錄（空行、「開始服務」這類訊息），不算錯誤，直接略過。
var errSkip = errors.New("不是請求的紀錄")

// ParseLine 解析一行 log：開頭是 { 的當成 slog 的 JSON 格式，其他的當成 access log。
func ParseLine(line string) (Entry, error) {
	switch {
	case strings.TrimSpace(line) == "":
		return Entry{}, errSkip
	case strings.HasPrefix(line, "{"):
		return parseJSON(line)
	}
	return parseAccess(line)
}

// accessRE 對應這種格式（類似 nginx 的 access log，最後一欄是處理時間，單位是秒）：
//
//	192.0.2.10 - - [21/Nov/2026:09:30:01 +0800] "GET /menu HTTP/1.1" 200 1532 0.012
var accessRE = regexp.MustCompile(
	`^\S+ \S+ \S+ \[(?<time>[^\]]+)\] "(?<method>[A-Z]+) (?<path>\S+) [^"]*" (?<status>\d{3}) (?:\d+|-) (?<took>\d+\.\d+)$`)

// accessLayout 是 access log 裡的時間格式，用 Go 的參考時間寫成。
const accessLayout = "02/Jan/2006:15:04:05 -0700"

func parseAccess(line string) (Entry, error) {
	m := accessRE.FindStringSubmatch(line)
	if m == nil {
		return Entry{}, errors.New("格式不對")
	}
	get := func(name string) string { return m[accessRE.SubexpIndex(name)] }

	t, err := time.Parse(accessLayout, get("time"))
	if err != nil {
		return Entry{}, fmt.Errorf("時間不對：%w", err)
	}
	status, _ := strconv.Atoi(get("status"))         // regexp 已經保證是三位數字
	took, _ := time.ParseDuration(get("took") + "s") // 用十進位解析，"1.003s" 不會有浮點誤差
	return Entry{Time: t, Method: get("method"), Path: get("path"), Status: status, Took: took}, nil
}

// slogLine 是 slog.NewJSONHandler 寫出的一行（例如專題 4 的 API）。took 是奈秒的整數。
type slogLine struct {
	Time   time.Time `json:"time"`
	Msg    string    `json:"msg"`
	Method string    `json:"method"`
	Path   string    `json:"path"`
	Status int       `json:"status"`
	Took   int64     `json:"took"`
}

func parseJSON(line string) (Entry, error) {
	var l slogLine
	if err := json.Unmarshal([]byte(line), &l); err != nil {
		return Entry{}, fmt.Errorf("JSON 不對：%w", err)
	}
	if l.Msg != "request" {
		return Entry{}, errSkip
	}
	return Entry{Time: l.Time, Method: l.Method, Path: l.Path, Status: l.Status, Took: time.Duration(l.Took)}, nil
}

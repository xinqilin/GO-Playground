package main

import (
	"encoding/json"
	"testing"
)

func Example() {
	main()
	// Output:
	// == time：用 2006-01-02 當格式範本 ==
	// 2026-10-09 14:30:00
	// 2026-10-09T14:30:00Z
	// 16:00 1h30m0s
	// 75 <nil>
	// Friday <nil>
	// == encoding/json ==
	// {"name":"小明","age":30} <nil>
	// 小華 true [go java] <nil>
	// == os：讀寫檔案 ==
	// 第一行
	// 第二行
	// 2 行
	// == log/slog：結構化 log ==
	// level=INFO msg=使用者登入 user=bill attempt=2
	// level=WARN msg=額度快用完 remaining=3
	// == uuid（Go 1.27）==
	// 36
	// true
	// 017f22e2-79b0-7cc3-98c4-dc0c0c07398f
}

// JSON 來回轉換之後，資料要一模一樣（email 沒有匯出，不在比較範圍內）。
func TestUserJSONRoundTrip(t *testing.T) {
	in := User{Name: "阿德", Age: new(0), Tags: []string{"a"}}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var out User
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out.Name != in.Name || out.Age == nil || *out.Age != 0 || len(out.Tags) != 1 {
		t.Errorf("round trip: got %+v from %s", out, b)
	}
}

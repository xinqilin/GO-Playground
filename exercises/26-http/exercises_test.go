package ticket

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/synctest"
	"uuid"
)

// result 是一個請求的結果：狀態碼、標頭、去掉結尾換行的內容。
type result struct {
	code   int
	header http.Header
	body   string
}

func do(h http.Handler, req *http.Request) result {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	res := rec.Result()
	b, _ := io.ReadAll(res.Body)
	return result{res.StatusCode, res.Header, strings.TrimRight(string(b), "\n")}
}

func TestQueryInt(t *testing.T) {
	tests := []struct {
		target  string
		want    int
		wantErr bool
	}{
		{"/tickets", 20, false},
		{"/tickets?limit=", 20, false},
		{"/tickets?limit=5", 5, false},
		{"/tickets?limit=1", 1, false},
		{"/tickets?limit=100", 100, false},
		{"/tickets?limit=0", 0, true},
		{"/tickets?limit=101", 0, true},
		{"/tickets?limit=ten", 0, true},
		{"/tickets?page=3", 20, false},
	}
	for _, tt := range tests {
		got, err := QueryInt(httptest.NewRequest("GET", tt.target, nil), "limit", 20, 1, 100)
		if tt.wantErr {
			if err == nil || !strings.Contains(err.Error(), "limit") {
				t.Errorf("%s：err = %v，want 提到 limit 的錯誤", tt.target, err)
			}
			continue
		}
		if err != nil || got != tt.want {
			t.Errorf("%s = %d（%v），want %d", tt.target, got, err, tt.want)
		}
	}
}

func TestRequireKey(t *testing.T) {
	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "進來了") })
	tests := []struct {
		key, sent string
		want      int
	}{
		{"s3cret", "s3cret", 200},
		{"s3cret", "", 401},
		{"s3cret", "s3cre", 401},
		{"s3cret", "S3CRET", 401},
		{"", "", 401},
	}
	for _, tt := range tests {
		req := httptest.NewRequest("POST", "/next", nil)
		if tt.sent != "" {
			req.Header.Set("X-Staff-Key", tt.sent)
		}
		got := do(RequireKey(tt.key, ok), req)
		if got.code != tt.want {
			t.Errorf("key %q 送 %q：狀態碼 %d，want %d", tt.key, tt.sent, got.code, tt.want)
			continue
		}
		if tt.want == 401 && (got.body != `{"error":"需要員工金鑰"}` || got.header.Get("Content-Type") != "application/json") {
			t.Errorf("key %q 送 %q：401 的內容 %q、Content-Type %q", tt.key, tt.sent, got.body, got.header.Get("Content-Type"))
		}
		if tt.want == 200 && got.body != "進來了" {
			t.Errorf("通過時應該交給 next，內容是 %q", got.body)
		}
	}
}

func TestRequestID(t *testing.T) {
	var seen string
	h := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = FromContext(r.Context())
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Request-ID", "order-42")
	got := do(h, req)
	if seen != "order-42" || got.header.Get("X-Request-ID") != "order-42" {
		t.Errorf("合法的 id 要沿用：context %q，標頭 %q", seen, got.header.Get("X-Request-ID"))
	}

	for _, bad := range []string{"", "有中文", "a b", "<script>", strings.Repeat("x", 65)} {
		req := httptest.NewRequest("GET", "/", nil)
		if bad != "" {
			req.Header.Set("X-Request-ID", bad)
		}
		got := do(h, req)
		id := got.header.Get("X-Request-ID")
		if _, err := uuid.Parse(id); err != nil || seen != id {
			t.Errorf("送 %q：標頭 %q、context %q，want 同一個新的 UUID", bad, id, seen)
		}
	}

	if id := FromContext(httptest.NewRequest("GET", "/", nil).Context()); id != "" {
		t.Errorf("沒經過 RequestID 時 FromContext = %q，want 空字串", id)
	}
}

// TestMaxInFlight 用 synctest：synctest.Wait 等到其他 goroutine 都卡住，
// 就能確定「第 3 個請求是立刻被拒絕，還是在等」，不用靠 sleep 猜。
func TestMaxInFlight(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		release := make(chan struct{})
		slow := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			<-release
			io.WriteString(w, "做好了")
		})
		h := MaxInFlight(2, slow)
		results := make(chan result, 3)
		send := func() { results <- do(h, httptest.NewRequest("GET", "/", nil)) }

		go send()
		go send()
		synctest.Wait() // 兩個請求都在 slow 裡面等

		go send()
		synctest.Wait()
		pending := 2
		select {
		case busy := <-results:
			if busy.code != 503 || busy.header.Get("Retry-After") != "1" || busy.body != `{"error":"太忙了，請稍後再試"}` {
				t.Errorf("第 3 個請求：%d，Retry-After %q，內容 %q；want 503、1、太忙了", busy.code, busy.header.Get("Retry-After"), busy.body)
			}
		default:
			t.Error("第 3 個請求在等待，應該立刻回 503")
			pending = 3
		}

		close(release)
		for range pending {
			if r := <-results; r.code != 200 || r.body != "做好了" {
				t.Errorf("在裡面的請求：%d %q，want 200 做好了", r.code, r.body)
			}
		}
		if again := do(h, httptest.NewRequest("GET", "/", nil)); again.code != 200 {
			t.Errorf("前面的做完之後，新的請求應該能進去：%d", again.code)
		}
	})
}

func TestNewServer(t *testing.T) {
	srv := httptest.NewTestServer(t, NewServer(&Queue{}))
	client := srv.Client()
	call := func(method, path string) result {
		req, err := http.NewRequest(method, "http://queue.test"+path, nil)
		if err != nil {
			t.Fatal(err)
		}
		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		b, _ := io.ReadAll(res.Body)
		return result{res.StatusCode, res.Header, strings.TrimRight(string(b), "\n")}
	}

	steps := []struct {
		method, path string
		code         int
		body         string
		location     string
	}{
		{"POST", "/next", 204, "", ""},
		{"POST", "/tickets", 201, `{"number":1,"ahead":0}`, "/tickets/1"},
		{"POST", "/tickets", 201, `{"number":2,"ahead":1}`, "/tickets/2"},
		{"POST", "/tickets", 201, `{"number":3,"ahead":2}`, "/tickets/3"},
		{"GET", "/tickets/3", 200, `{"number":3,"ahead":2}`, ""},
		{"POST", "/next", 200, `{"number":1}`, ""},
		{"GET", "/tickets/3", 200, `{"number":3,"ahead":1}`, ""},
		{"GET", "/tickets/1", 404, "", ""},
		{"GET", "/tickets/0", 400, "", ""},
		{"GET", "/tickets/abc", 400, "", ""},
		{"DELETE", "/tickets/3", 405, "", ""},
		{"GET", "/next", 405, "", ""},
	}
	for _, s := range steps {
		got := call(s.method, s.path)
		if got.code != s.code {
			t.Errorf("%s %s：狀態碼 %d，want %d（內容 %q）", s.method, s.path, got.code, s.code, got.body)
			continue
		}
		if s.body != "" && got.body != s.body {
			t.Errorf("%s %s：內容 %q，want %q", s.method, s.path, got.body, s.body)
		}
		if got.body != "" && s.code < 300 && got.header.Get("Content-Type") != "application/json" {
			t.Errorf("%s %s：Content-Type %q，want application/json", s.method, s.path, got.header.Get("Content-Type"))
		}
		if (s.code == 400 || s.code == 404) && !strings.HasPrefix(got.body, `{"error":`) {
			t.Errorf("%s %s：錯誤的內容 %q，want {\"error\":...}", s.method, s.path, got.body)
		}
		if s.location != "" && got.header.Get("Location") != s.location {
			t.Errorf("%s %s：Location %q，want %q", s.method, s.path, got.header.Get("Location"), s.location)
		}
		if s.code == 204 && got.body != "" {
			t.Errorf("%s %s：204 不該有內容，拿到 %q", s.method, s.path, got.body)
		}
	}
}

func TestRename(t *testing.T) {
	s := NewSettings("阿嬤便當")

	got := do(http.HandlerFunc(s.Rename), httptest.NewRequest("PUT", "/settings/name", strings.NewReader(`{"name":""}`)))
	if got.code != 400 || got.body != "需要 name" || s.Name() != "阿嬤便當" {
		t.Errorf("空的 name：%d %q，店名變成 %q；want 400、需要 name、店名不變", got.code, got.body, s.Name())
	}

	got = do(http.HandlerFunc(s.Rename), httptest.NewRequest("PUT", "/settings/name", strings.NewReader(`不是 JSON`)))
	if got.code != 400 || got.body != "需要 name" || s.Name() != "阿嬤便當" {
		t.Errorf("解不開的內容：%d %q，店名變成 %q；want 400、需要 name、店名不變", got.code, got.body, s.Name())
	}

	got = do(http.HandlerFunc(s.Rename), httptest.NewRequest("PUT", "/settings/name", strings.NewReader(`{"name":"阿公便當"}`)))
	if got.code != 200 || got.body != `{"name":"阿公便當"}` || s.Name() != "阿公便當" {
		t.Errorf("合法的 name：%d %q，店名 %q", got.code, got.body, s.Name())
	}
}

package bus

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

// tracker 是測試用的 RoundTripper：不經過網路，把請求直接交給 h，
// 並記錄回應的 body 有沒有被關掉。
type tracker struct {
	h       http.Handler
	mu      sync.Mutex
	open    int // 目前還沒關的 body
	maxOpen int // 同時沒關的 body 最多有幾個
}

func (t *tracker) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		defer req.Body.Close()
	}
	rec := httptest.NewRecorder()
	t.h.ServeHTTP(rec, req)
	resp := rec.Result()
	t.mu.Lock()
	t.open++
	t.maxOpen = max(t.maxOpen, t.open)
	t.mu.Unlock()
	resp.Body = &trackedBody{ReadCloser: resp.Body, t: t}
	return resp, nil
}

func (t *tracker) stillOpen() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.open
}

type trackedBody struct {
	io.ReadCloser
	t    *tracker
	once sync.Once
}

func (b *trackedBody) Close() error {
	b.once.Do(func() {
		b.t.mu.Lock()
		b.t.open--
		b.t.mu.Unlock()
	})
	return b.ReadCloser.Close()
}

func TestBuildURL(t *testing.T) {
	tests := []struct {
		base, stop string
		routes     []string
		want       string
	}{
		{"https://bus.test/v1", "市政府", []string{"307", "藍5"},
			"https://bus.test/v1/stops/%E5%B8%82%E6%94%BF%E5%BA%9C/arrivals?route=307&route=%E8%97%8D5"},
		{"https://bus.test/v1/", "A/B 月台", nil,
			"https://bus.test/v1/stops/A%2FB%20%E6%9C%88%E5%8F%B0/arrivals"},
		{"http://bus.test", "站(1)", []string{"紅30"},
			"http://bus.test/stops/%E7%AB%99%281%29/arrivals?route=%E7%B4%8530"},
	}
	for _, tt := range tests {
		got, err := BuildURL(tt.base, tt.stop, tt.routes)
		if err != nil || got != tt.want {
			t.Errorf("BuildURL(%q, %q, %q) =\n  %q, %v\nwant\n  %q", tt.base, tt.stop, tt.routes, got, err, tt.want)
		}
	}
	for _, bad := range []string{"bus.test/v1", "http://%zz", "/v1"} {
		if got, err := BuildURL(bad, "市政府", nil); err == nil {
			t.Errorf("BuildURL(%q) = %q，want 錯誤", bad, got)
		}
	}
}

type ctxKey struct{}

func TestGetArrivals(t *testing.T) {
	var gotAccept string
	var gotCtx any
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ok", func(w http.ResponseWriter, r *http.Request) {
		gotAccept, gotCtx = r.Header.Get("Accept"), r.Context().Value(ctxKey{})
		io.WriteString(w, `[{"route":"307","minutes":3},{"route":"藍5","minutes":12}]`)
	})
	mux.HandleFunc("GET /missing", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, `{"message":"沒有這個站"}`)
	})
	mux.HandleFunc("GET /busy", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "busy", http.StatusServiceUnavailable)
	})
	mux.HandleFunc("GET /broken", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `[{"route":`)
	})
	mux.HandleFunc("GET /huge", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `[{"route":"`+strings.Repeat("9", maxBody)+`","minutes":1}]`)
	})
	tr := &tracker{h: mux}
	client := &http.Client{Transport: tr}
	ctx := context.WithValue(t.Context(), ctxKey{}, "我的 ctx")

	got, err := GetArrivals(ctx, client, "http://bus.test/ok")
	if want := []Arrival{{"307", 3}, {"藍5", 12}}; err != nil || !slices.Equal(got, want) {
		t.Errorf("200：%v, %v，want %v", got, err, want)
	}
	if gotAccept != "application/json" {
		t.Errorf("Accept 標頭是 %q，want application/json", gotAccept)
	}
	if gotCtx != "我的 ctx" {
		t.Error("請求沒有帶著 ctx（用 http.NewRequestWithContext 建立請求）")
	}

	for _, c := range []struct {
		path string
		want StatusError
	}{
		{"/missing", StatusError{404, "沒有這個站"}},
		{"/busy", StatusError{503, ""}},
	} {
		_, err := GetArrivals(ctx, client, "http://bus.test"+c.path)
		if se, ok := errors.AsType[*StatusError](err); !ok || *se != c.want {
			t.Errorf("%s：err = %v，want *StatusError %+v", c.path, err, c.want)
		}
	}

	for _, path := range []string{"/broken", "/huge"} {
		_, err := GetArrivals(ctx, client, "http://bus.test"+path)
		if _, isStatus := errors.AsType[*StatusError](err); err == nil || isStatus {
			t.Errorf("%s：err = %v，want JSON 解析失敗的錯誤（/huge 超過 maxBody，只讀一部分就解析不完整）", path, err)
		}
	}

	if n := tr.stillOpen(); n != 0 {
		t.Errorf("還有 %d 個回應的 body 沒關", n)
	}
}

func TestAuthTransport(t *testing.T) {
	var seen []string
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Header.Get("Authorization")+" | "+r.Header.Get("Accept"))
	})
	client := &http.Client{Transport: &AuthTransport{Token: "k-123", Next: &tracker{h: h}}}

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://bus.test/stops", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", "application/json")
	for range 2 { // 同一個請求送兩次，像重試那樣
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("client.Do：%v", err)
		}
		resp.Body.Close()
	}

	want := []string{"Bearer k-123 | application/json", "Bearer k-123 | application/json"}
	if !slices.Equal(seen, want) {
		t.Errorf("對方收到的標頭 %q，want %q", seen, want)
	}
	if a := req.Header.Get("Authorization"); a != "" {
		t.Errorf("原本的請求被改了：Authorization = %q（要先 Clone）", a)
	}
}

func TestRetryAfter(t *testing.T) {
	now := time.Date(2026, 11, 21, 8, 0, 0, 0, time.UTC)
	tests := []struct {
		value  string // 空字串表示沒有這個標頭
		want   time.Duration
		wantOK bool
	}{
		{"120", 2 * time.Minute, true},
		{"0", 0, true},
		{"Sat, 21 Nov 2026 08:01:30 GMT", 90 * time.Second, true},
		{"Saturday, 21-Nov-26 08:00:05 GMT", 5 * time.Second, true},
		{"Sat, 21 Nov 2026 07:59:00 GMT", 0, true},
		{"", 0, false},
		{"-5", 0, false},
		{"1.5", 0, false},
		{"等一下", 0, false},
	}
	for _, tt := range tests {
		h := http.Header{}
		if tt.value != "" {
			h.Set("Retry-After", tt.value)
		}
		got, ok := RetryAfter(h, now)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("Retry-After %q = %v, %v，want %v, %v", tt.value, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestCache(t *testing.T) {
	var mu sync.Mutex
	version, full, notModified := "v1", 0, 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if r.URL.Path == "/down" {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		etag := `"` + version + `"`
		if r.Header.Get("If-None-Match") == etag {
			notModified++
			w.WriteHeader(http.StatusNotModified)
			return
		}
		full++
		w.Header().Set("ETag", etag)
		io.WriteString(w, "時刻表 "+version)
	})
	tr := &tracker{h: h}
	client := &http.Client{Transport: tr}
	ctx := t.Context()
	var c Cache
	get := func(want string) []byte {
		t.Helper()
		b, err := c.Get(ctx, client, "http://bus.test/timetable")
		if err != nil || string(b) != want {
			t.Errorf("Get = %q, %v，want %q", b, err, want)
		}
		return b
	}

	get("時刻表 v1")
	get("時刻表 v1")
	mu.Lock()
	version = "v2"
	mu.Unlock()
	get("時刻表 v2")
	b := get("時刻表 v2")
	if full != 2 || notModified != 2 {
		t.Errorf("對方送了 %d 次完整內容、%d 次 304，want 2 和 2", full, notModified)
	}

	if len(b) > 0 {
		b[0] = 'X' // 呼叫的一方改了拿到的內容，不能影響快取
	}
	get("時刻表 v2")

	_, err := c.Get(ctx, client, "http://bus.test/down")
	if se, ok := errors.AsType[*StatusError](err); !ok || *se != (StatusError{Code: 502}) {
		t.Errorf("502：err = %v，want *StatusError{Code: 502}", err)
	}

	var wg sync.WaitGroup
	for range 20 {
		wg.Go(func() {
			b, err := c.Get(ctx, client, "http://bus.test/timetable")
			if err != nil || string(b) != "時刻表 v2" {
				t.Errorf("同時取：%q, %v", b, err)
			}
		})
	}
	wg.Wait()

	if n := tr.stillOpen(); n != 0 {
		t.Errorf("還有 %d 個回應的 body 沒關", n)
	}
}

func TestHealthy(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/down") {
			http.Error(w, "壞了", http.StatusServiceUnavailable)
			return
		}
		io.WriteString(w, "ok")
	})
	tr := &tracker{h: h}
	client := &http.Client{Transport: tr}
	urls := []string{
		"http://bus.test/a", "http://bus.test/down1", "://壞掉的網址",
		"http://bus.test/b", "http://bus.test/down2", "http://bus.test/c",
	}

	if n := Healthy(t.Context(), client, urls); n != 3 {
		t.Errorf("Healthy = %d，want 3", n)
	}
	if tr.maxOpen > 1 {
		t.Errorf("同時有 %d 個回應的 body 沒關，want 最多 1 個：每個回應用完就要關掉", tr.maxOpen)
	}
	if n := tr.stillOpen(); n != 0 {
		t.Errorf("Healthy 結束時還有 %d 個 body 沒關", n)
	}
}

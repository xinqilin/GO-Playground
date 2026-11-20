package main

import (
	"encoding/json/v2"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// fixedNow 是測試用的時鐘：台北時間 2026 年 11 月 21 日早上 9 點半。
func fixedNow() time.Time {
	return time.Date(2026, 11, 21, 9, 30, 0, 0, time.FixedZone("CST", 8*60*60))
}

// newTestHandler 組出完整的 API：時間固定、log 丟掉。path 是空字串時資料只放在記憶體。
func newTestHandler(path string) http.Handler {
	store, err := OpenStore(path)
	if err != nil {
		panic(err)
	}
	logger := slog.New(slog.DiscardHandler)
	return (&api{svc: &service{store: store, now: fixedNow}, logger: logger}).routes()
}

// send 用 httptest 送一個請求給 h，印出狀態碼、Location 標頭（有的話）和內容。不用開 port。
func send(h http.Handler, method, target, body string) {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	res := rec.Result()
	fmt.Printf("%s %s → %d", method, target, res.StatusCode)
	if loc := res.Header.Get("Location"); loc != "" {
		fmt.Printf("，Location: %s", loc)
	}
	fmt.Println()
	if b, _ := io.ReadAll(res.Body); len(b) > 0 {
		fmt.Println(string(b))
	}
}

func Example_api() {
	h := newTestHandler("")
	send(h, "POST", "/todos", `{"title": "  訂明天的便當  ", "due": "2026-11-22"}`)
	send(h, "POST", "/todos", `{"title": "繳電話費"}`)
	send(h, "PATCH", "/todos/1", `{"done": true}`)
	send(h, "GET", "/todos?done=false", "")
	send(h, "DELETE", "/todos/2", "")
	send(h, "GET", "/todos/2", "")
	// Output:
	// POST /todos → 201，Location: /todos/1
	// {"id":1,"title":"訂明天的便當","done":false,"due":"2026-11-22","created_at":"2026-11-21T09:30:00+08:00"}
	// POST /todos → 201，Location: /todos/2
	// {"id":2,"title":"繳電話費","done":false,"created_at":"2026-11-21T09:30:00+08:00"}
	// PATCH /todos/1 → 200
	// {"id":1,"title":"訂明天的便當","done":true,"due":"2026-11-22","created_at":"2026-11-21T09:30:00+08:00"}
	// GET /todos?done=false → 200
	// [{"id":2,"title":"繳電話費","done":false,"created_at":"2026-11-21T09:30:00+08:00"}]
	// DELETE /todos/2 → 204
	// GET /todos/2 → 404
	// {"error":"找不到這筆待辦"}
}

func Example_validation() {
	h := newTestHandler("")
	send(h, "POST", "/todos", `{"title": "   ", "due": "2026-11-20"}`)
	send(h, "POST", "/todos", `{"title": "繳電話費", "due": "11/30"}`)
	send(h, "POST", "/todos", `{"titel": "繳電話費"}`)
	send(h, "PATCH", "/todos/1", `{}`)
	send(h, "GET", "/todos?done=yes", "")
	// Output:
	// POST /todos → 422
	// {"error":"輸入不合法","fields":{"due":"不能早於今天","title":"不能是空的"}}
	// POST /todos → 422
	// {"error":"輸入不合法","fields":{"due":"格式要是 2006-01-02"}}
	// POST /todos → 400
	// {"error":"不認得的欄位 /titel"}
	// PATCH /todos/1 → 422
	// {"error":"輸入不合法","fields":{"body":"至少要改一個欄位"}}
	// GET /todos?done=yes → 400
	// {"error":"done 要是 true 或 false"}
}

// TestConcurrentCreate：50 個請求同時新增，編號不能重複，檔案裡也要剛好有 50 筆。
func TestConcurrentCreate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "todos.json")
	h := newTestHandler(path)
	var wg sync.WaitGroup
	for i := range 50 {
		wg.Go(func() {
			req := httptest.NewRequest("POST", "/todos", strings.NewReader(fmt.Sprintf(`{"title": "便當 %d"}`, i)))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != http.StatusCreated {
				t.Errorf("狀態碼 %d：%s", rec.Code, rec.Body)
			}
		})
	}
	wg.Wait()

	s, err := OpenStore(path) // 從檔案重新讀一次
	if err != nil {
		t.Fatal(err)
	}
	for i, todo := range s.List() {
		if todo.ID != i+1 {
			t.Fatalf("第 %d 筆的 ID 是 %d", i, todo.ID)
		}
	}
	if n := len(s.List()); n != 50 {
		t.Errorf("檔案裡有 %d 筆，想要 50 筆", n)
	}
}

// TestEndToEnd 用真的 http.Client 透過 httptest.NewTestServer（記憶體裡的網路，Go 1.27 起）呼叫整個服務。
func TestEndToEnd(t *testing.T) {
	srv := httptest.NewTestServer(t, newTestHandler(filepath.Join(t.TempDir(), "todos.json")))
	client := srv.Client()

	res, err := client.Post(srv.URL+"/todos", "application/json", strings.NewReader(`{"title": "訂便當"}`))
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("POST 狀態碼 %d", res.StatusCode)
	}

	res, err = client.Get(srv.URL + res.Header.Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var got Todo
	if err := json.UnmarshalRead(res.Body, &got); err != nil {
		t.Fatal(err)
	}
	if got.ID != 1 || got.Title != "訂便當" || !got.CreatedAt.Equal(fixedNow()) {
		t.Errorf("GET 拿到 %+v", got)
	}
}

func TestPanicBecomes500(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	h := logRequests(logger, recoverPanic(logger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("便當掉到地上了")
	})))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/todos", nil))
	if rec.Code != http.StatusInternalServerError || !strings.Contains(rec.Body.String(), "伺服器發生錯誤") {
		t.Errorf("狀態碼 %d，內容 %s", rec.Code, rec.Body)
	}
}

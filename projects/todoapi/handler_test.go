package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func newTestAPI() *api {
	return &api{store: NewStore(), logger: slog.New(slog.DiscardHandler)}
}

// send 直接呼叫 handler，不開 port。
func send(h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func Example() {
	h := newTestAPI().routes()
	steps := []struct{ method, path, body string }{
		{"POST", "/todos", `{"title":"寫 Go 手冊"}`},
		{"GET", "/todos", ""},
		{"PATCH", "/todos/1", `{"done":true}`},
		{"DELETE", "/todos/1", ""},
		{"GET", "/todos/1", ""},
	}
	for _, s := range steps {
		rec := send(h, s.method, s.path, s.body)
		line := fmt.Sprintf("%s %s → %d %s", s.method, s.path, rec.Code, rec.Body.String())
		fmt.Println(strings.TrimSpace(line)) // 204 沒有 body，去掉行尾空白
	}
	// Output:
	// POST /todos → 201 {"id":1,"title":"寫 Go 手冊","done":false}
	// GET /todos → 200 [{"id":1,"title":"寫 Go 手冊","done":false}]
	// PATCH /todos/1 → 200 {"id":1,"title":"寫 Go 手冊","done":true}
	// DELETE /todos/1 → 204
	// GET /todos/1 → 404 {"error":"找不到這筆待辦"}
}

func TestEmptyListIsArray(t *testing.T) {
	rec := send(newTestAPI().routes(), "GET", "/todos", "")
	if got := strings.TrimSpace(rec.Body.String()); got != "[]" {
		t.Errorf("空清單應該是 []，got %s", got)
	}
}

func TestValidation(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
	}{
		{"不是 JSON", "POST", "/todos", "hello", http.StatusBadRequest},
		{"多了不認得的欄位", "POST", "/todos", `{"title":"a","x":1}`, http.StatusBadRequest},
		{"標題是空白", "POST", "/todos", `{"title":"   "}`, http.StatusUnprocessableEntity},
		{"標題太長", "POST", "/todos", `{"title":"` + strings.Repeat("字", 101) + `"}`, http.StatusUnprocessableEntity},
		{"id 不是數字", "GET", "/todos/abc", "", http.StatusBadRequest},
		{"id 是 0", "GET", "/todos/0", "", http.StatusBadRequest},
		{"沒有 done 欄位", "PATCH", "/todos/1", `{}`, http.StatusBadRequest},
		{"更新不存在的", "PATCH", "/todos/99", `{"done":true}`, http.StatusNotFound},
		{"刪除不存在的", "DELETE", "/todos/99", "", http.StatusNotFound},
		{"方法不支援", "PUT", "/todos/1", "", http.StatusMethodNotAllowed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := send(newTestAPI().routes(), tt.method, tt.path, tt.body)
			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d（body %s）", rec.Code, tt.wantStatus, rec.Body.String())
			}
		})
	}
}

func TestRecoverPanic(t *testing.T) {
	a := newTestAPI()
	h := a.recoverPanic(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	}))
	rec := send(h, "GET", "/", "")
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
}

// 用 go test -race 執行：同時新增 100 筆，ID 不能重複。
func TestConcurrentCreate(t *testing.T) {
	s := NewStore()
	var wg sync.WaitGroup
	for range 100 {
		wg.Go(func() { s.Create("x") })
	}
	wg.Wait()
	seen := map[int]bool{}
	for _, todo := range s.List() {
		if seen[todo.ID] {
			t.Fatalf("重複的 ID %d", todo.ID)
		}
		seen[todo.ID] = true
	}
	if len(seen) != 100 {
		t.Errorf("got %d todos, want 100", len(seen))
	}
}

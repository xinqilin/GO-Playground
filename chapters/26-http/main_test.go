package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// httptest.NewRecorder 直接呼叫 handler，不用真的開 port，測試又快又穩定。
func Example_routes() {
	h := routes(newNoteStore(), slog.New(slog.DiscardHandler))
	do := func(method, path, body string) {
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		fmt.Print(rec.Code, " ", rec.Body.String())
	}
	do("GET", "/hello/Gopher", "")
	do("POST", "/notes", `{"text":"買牛奶"}`)
	do("GET", "/notes/1", "")
	do("GET", "/notes/2", "")
	do("DELETE", "/notes/1", "")
	// Output:
	// 200 哈囉，Gopher
	// 201 {"id":1,"text":"買牛奶"}
	// 200 {"id":1,"text":"買牛奶"}
	// 404 404 page not found
	// 405 Method Not Allowed
}

func TestCreateBadRequest(t *testing.T) {
	h := routes(newNoteStore(), slog.New(slog.DiscardHandler))
	for _, body := range []string{"", "not json", `{"text":""}`} {
		req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(body))
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("body %q: status = %d, want 400", body, rec.Code)
		}
	}
}

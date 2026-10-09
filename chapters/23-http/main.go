// 第 23 章：HTTP 服務。
//
// HTTP server 像一間餐廳：路由（mux）是帶位的服務生，handler 是廚師，middleware 是門口的登記簿。
// 執行：go run ./chapters/23-http，再開另一個終端機用 curl 測試。
package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// noteStore 把筆記放在記憶體裡。handler 會被多個 goroutine 同時呼叫，所以要加鎖。
type noteStore struct {
	mu    sync.Mutex
	notes map[int]string
	next  int
}

func newNoteStore() *noteStore {
	return &noteStore{notes: make(map[int]string), next: 1}
}

// routes 用 Go 1.22 起的路由寫法：「方法 + 路徑」，路徑裡可以有 {變數}。
func routes(s *noteStore, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello/{name}", hello)
	mux.HandleFunc("POST /notes", s.create)
	mux.HandleFunc("GET /notes/{id}", s.get)
	return logRequests(logger, mux)
}

// hello 用 r.PathValue 取出路徑裡的 {name}。
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "哈囉，%s\n", r.PathValue("name"))
}

// create 讀 JSON body，存一則筆記，回傳 201 和新筆記。
func (s *noteStore) create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Text == "" {
		http.Error(w, `需要 JSON：{"text": "..."}`, http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	id := s.next
	s.next++
	s.notes[id] = in.Text
	s.mu.Unlock()
	writeJSON(w, http.StatusCreated, map[string]any{"id": id, "text": in.Text})
}

// get 依 id 找筆記，找不到回 404。
func (s *noteStore) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "id 要是數字", http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	text, ok := s.notes[id]
	s.mu.Unlock()
	if !ok {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"id": id, "text": text})
}

// writeJSON 設好 Content-Type 和狀態碼，再把 v 編成 JSON 寫出去。
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("write json", "err", err)
	}
}

// logRequests 是 middleware：包住另一個 handler，在它前後加上自己的邏輯。
func logRequests(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Info("request", "method", r.Method, "path", r.URL.Path, "took", time.Since(start))
	})
}

func main() {
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           routes(newNoteStore(), slog.Default()),
		ReadHeaderTimeout: 5 * time.Second, // 一定要設逾時，否則慢速連線可以把 server 拖垮
	}
	slog.Info("listening", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

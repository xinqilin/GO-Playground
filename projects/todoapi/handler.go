package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// maxTitleLen 是標題最多幾個字（以 rune 計算，中文一個字算一個）。
const maxTitleLen = 100

// api 把 HTTP 層需要的依賴放在一起，handler 都是它的方法。
type api struct {
	store  *Store
	logger *slog.Logger
}

// routes 組出完整的 handler：最外層接住 panic，再記錄每個請求，最後交給路由。
func (a *api) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos", a.list)
	mux.HandleFunc("POST /todos", a.create)
	mux.HandleFunc("GET /todos/{id}", a.get)
	mux.HandleFunc("PATCH /todos/{id}", a.update)
	mux.HandleFunc("DELETE /todos/{id}", a.delete)
	return a.recoverPanic(a.logRequests(mux))
}

func (a *api) list(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.store.List())
}

func (a *api) create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Title string `json:"title"`
	}
	if err := decodeJSON(w, r, &in); err != nil {
		writeError(w, http.StatusBadRequest, `body 必須是 {"title": "..."}`)
		return
	}
	title := strings.TrimSpace(in.Title)
	if n := utf8.RuneCountInString(title); n == 0 || n > maxTitleLen {
		writeError(w, http.StatusUnprocessableEntity, "title 必須是 1 到 100 個字")
		return
	}
	writeJSON(w, http.StatusCreated, a.store.Create(title))
}

func (a *api) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	t, err := a.store.Get(id)
	if err != nil {
		a.writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (a *api) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var in struct {
		Done *bool `json:"done"` // 用指標才分得出「沒傳」和「傳 false」
	}
	if err := decodeJSON(w, r, &in); err != nil || in.Done == nil {
		writeError(w, http.StatusBadRequest, `body 必須是 {"done": true 或 false}`)
		return
	}
	t, err := a.store.SetDone(id, *in.Done)
	if err != nil {
		a.writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (a *api) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := a.store.Delete(id); err != nil {
		a.writeStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// writeStoreError 把 Store 的錯誤轉成 HTTP 狀態碼；不認得的錯誤記 log 並回 500。
func (a *api) writeStoreError(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "找不到這筆待辦")
		return
	}
	a.logger.Error("store", "err", err)
	writeError(w, http.StatusInternalServerError, "伺服器內部錯誤")
}

// pathID 取出路徑裡的 {id}；不合法時直接寫出 400，回傳 false。
func pathID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "id 必須是正整數")
		return 0, false
	}
	return id, true
}

// decodeJSON 限制 body 大小、拒絕不認得的欄位，避免使用者傳錯欄位名卻沒發現。
func decodeJSON(w http.ResponseWriter, r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("write json", "err", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// statusRecorder 記下 handler 寫出的狀態碼，給 log 用。
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

// logRequests 記錄每個請求的方法、路徑、狀態碼和花費時間。
func (a *api) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		a.logger.Info("request",
			"method", r.Method, "path", r.URL.Path,
			"status", rec.status, "took", time.Since(start))
	})
}

// recoverPanic 在最外層接住 handler 的 panic，回 500 而不是讓連線直接斷掉（第 11 章）。
func (a *api) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				a.logger.Error("panic", "value", v, "path", r.URL.Path)
				writeError(w, http.StatusInternalServerError, "伺服器內部錯誤")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

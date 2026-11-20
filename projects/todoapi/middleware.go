package main

import (
	"log/slog"
	"net/http"
	"time"
)

// statusWriter 包住 ResponseWriter，記下 handler 寫出的狀態碼。
type statusWriter struct {
	http.ResponseWriter
	status int
}

func (sw *statusWriter) WriteHeader(code int) {
	if sw.status == 0 {
		sw.status = code
	}
	sw.ResponseWriter.WriteHeader(code)
}

func (sw *statusWriter) Write(b []byte) (int, error) {
	if sw.status == 0 {
		sw.status = http.StatusOK
	}
	return sw.ResponseWriter.Write(b)
}

// Unwrap 讓 http.NewResponseController 找得到裡面原本的 ResponseWriter。
func (sw *statusWriter) Unwrap() http.ResponseWriter { return sw.ResponseWriter }

// logRequests 在每個請求結束後記一筆結構化的 log。
func logRequests(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w}
		next.ServeHTTP(sw, r)
		logger.InfoContext(r.Context(), "request",
			"method", r.Method, "path", r.URL.Path, "status", sw.status, "took", time.Since(start))
	})
}

// recoverPanic 把 handler 的 panic 變成 500 回應，並記下來（26 章）。
func recoverPanic(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			v := recover()
			if v == nil {
				return
			}
			if v == http.ErrAbortHandler {
				panic(v)
			}
			logger.ErrorContext(r.Context(), "panic", "path", r.URL.Path, "value", v)
			writeJSON(w, http.StatusInternalServerError, errorBody{Error: "伺服器發生錯誤"})
		}()
		next.ServeHTTP(w, r)
	})
}

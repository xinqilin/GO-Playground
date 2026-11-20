package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

// handler 是整個服務對外的 Handler：路由外面包兩層 middleware。
// 外層先拿到請求、最後拿到結果，所以 logRequests 放最外面，才記得到 recoverPanic 改成的 500。
func (a *api) handler() http.Handler {
	return logRequests(a.logger, recoverPanic(a.logger, a.routes()))
}

// statusWriter 包住 ResponseWriter，記下 handler 寫出的狀態碼和位元組數。
type statusWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (sw *statusWriter) WriteHeader(code int) {
	if sw.status == 0 {
		sw.status = code
	}
	sw.ResponseWriter.WriteHeader(code)
}

func (sw *statusWriter) Write(b []byte) (int, error) {
	if sw.status == 0 {
		sw.status = http.StatusOK // 沒有呼叫 WriteHeader 就直接 Write，等於 200
	}
	n, err := sw.ResponseWriter.Write(b)
	sw.bytes += n
	return n, err
}

// Unwrap 讓 http.NewResponseController 找得到裡面原本的 ResponseWriter，
// Flush、SetWriteDeadline 這些功能才不會因為包了一層就失效。
func (sw *statusWriter) Unwrap() http.ResponseWriter {
	return sw.ResponseWriter
}

// logRequests 每個請求結束後記一筆：方法、路徑、對到的 pattern、狀態碼、大小、花了多久。
func logRequests(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w}
		next.ServeHTTP(sw, r)
		logger.InfoContext(r.Context(), "request",
			"method", r.Method, "path", r.URL.Path, "pattern", r.Pattern,
			"status", sw.status, "bytes", sw.bytes, "took", time.Since(start))
	})
}

// recoverPanic 把 handler 的 panic 變成 500，並記下來。沒有它的話，net/http 會記下 panic 並直接切斷連線，
// 客戶端只會看到連線被關掉，拿不到任何回應。
func recoverPanic(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			v := recover()
			if v == nil {
				return
			}
			if v == http.ErrAbortHandler {
				panic(v) // 這是 handler 故意中止連線的訊號，交給 net/http 處理
			}
			logger.ErrorContext(r.Context(), "panic", "path", r.URL.Path, "value", v)
			writeJSON(w, http.StatusInternalServerError, errorBody{Error: "伺服器發生錯誤"})
		}()
		next.ServeHTTP(w, r)
	})
}

// dropVolatile 拿掉時間和耗時，範例的輸出才會每次一樣。
func dropVolatile(groups []string, a slog.Attr) slog.Attr {
	if len(groups) == 0 && (a.Key == slog.TimeKey || a.Key == "took") {
		return slog.Attr{}
	}
	return a
}

// middlewareDemo 多註冊一個會 panic 的路由，看兩層 middleware 怎麼處理。
func middlewareDemo() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: dropVolatile}))
	a := &api{shop: newShop(), logger: logger}
	mux := a.routes()
	mux.HandleFunc("GET /oops", func(w http.ResponseWriter, r *http.Request) {
		var m map[string]int
		m["排骨飯"]++ // 寫入 nil map：panic
	})
	h := logRequests(logger, recoverPanic(logger, mux))

	send(h, "GET", "/healthz", "")
	send(h, "GET", "/orders/7", "")
	send(h, "GET", "/oops", "")
}

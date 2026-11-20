package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"time"
)

// handlerTransport 是一個 RoundTripper：不經過網路，直接把請求交給 h，再用 httptest 的記錄器收下回應。
// 本章的示範用它代替真的網路，輸出才會固定。
type handlerTransport struct {
	h http.Handler
}

func (t handlerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		defer req.Body.Close() // RoundTripper 一定要關掉請求的 body，不管成功或失敗
	}
	rec := httptest.NewRecorder()
	t.h.ServeHTTP(rec, req)
	return rec.Result(), nil
}

// logTransport 包住另一個 RoundTripper，記下每個請求的方法、網址、結果和花了多久。
type logTransport struct {
	next   http.RoundTripper
	logger *slog.Logger
}

func (t logTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := t.next.RoundTrip(req)
	if err != nil {
		t.logger.WarnContext(req.Context(), "request failed", "method", req.Method, "url", req.URL.Redacted(), "err", err)
		return nil, err
	}
	t.logger.InfoContext(req.Context(), "request", "method", req.Method, "url", req.URL.Redacted(),
		"status", resp.StatusCode, "took", time.Since(start))
	return resp, nil
}

// dropVolatile 拿掉時間和耗時，範例的輸出才會每次一樣。
func dropVolatile(groups []string, a slog.Attr) slog.Attr {
	if len(groups) == 0 && (a.Key == slog.TimeKey || a.Key == "took") {
		return slog.Attr{}
	}
	return a
}

// transportDemo 換掉 Transport，所有經過這個 client 的請求都會被記錄，呼叫的程式碼不用改。
func transportDemo() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: dropVolatile}))
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: logTransport{next: handlerTransport{fakeShop()}, logger: logger},
	}
	for _, u := range []string{"http://shop.test/menu", "http://amei:s3cret@shop.test/orders/7"} {
		resp, err := client.Get(u)
		if err != nil {
			continue
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

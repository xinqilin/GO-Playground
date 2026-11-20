package main

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// Fetcher 送出一次請求，回傳 HTTP 狀態碼。正式執行用 httpFetcher，測試換成假的。
type Fetcher interface {
	Fetch(ctx context.Context, rawURL string) (int, error)
}

// httpFetcher 先用 HEAD 檢查，不下載內容；伺服器不接受 HEAD（405）時改用 GET。
type httpFetcher struct{ client *http.Client }

func (f httpFetcher) Fetch(ctx context.Context, rawURL string) (int, error) {
	status, err := f.do(ctx, http.MethodHead, rawURL)
	if err == nil && status == http.StatusMethodNotAllowed {
		return f.do(ctx, http.MethodGet, rawURL)
	}
	return status, err
}

func (f httpFetcher) do(ctx context.Context, method, rawURL string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, nil)
	if err != nil {
		return 0, err
	}
	resp, err := f.client.Do(req)
	if err != nil {
		return 0, err
	}
	resp.Body.Close() // 只需要狀態碼，GET 的內容不讀，直接關掉
	return resp.StatusCode, nil
}

var errScheme = errors.New("只檢查 http 和 https 網址")

// validateURL 在送出請求之前，先擋掉明顯不對的網址。
func validateURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errScheme
	}
	if u.Host == "" {
		return errors.New("網址裡沒有主機名稱")
	}
	return nil
}

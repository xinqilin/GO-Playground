package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"testing"
	"time"
)

// roundTripFunc 讓一個函式當成 http.RoundTripper：請求不會真的送到網路上。
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// countingBody 記錄 Body 被關了幾次。
type countingBody struct {
	io.Reader
	closed *int
}

func (b countingBody) Close() error {
	*b.closed++
	return nil
}

func TestHTTPFetcher(t *testing.T) {
	tests := []struct {
		name        string
		head, get   int // 假的伺服器對 HEAD 和 GET 回的狀態碼
		wantStatus  int
		wantMethods []string
	}{
		{"HEAD 就夠了", 200, 200, 200, []string{"HEAD"}},
		{"404 不改用 GET", 404, 200, 404, []string{"HEAD"}},
		{"不接受 HEAD 就改用 GET", 405, 200, 200, []string{"HEAD", "GET"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var methods []string
			closed := 0
			client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				methods = append(methods, r.Method)
				status := tt.head
				if r.Method == http.MethodGet {
					status = tt.get
				}
				body := countingBody{strings.NewReader("今日菜單"), &closed}
				return &http.Response{StatusCode: status, Body: body, Request: r}, nil
			})}

			status, err := httpFetcher{client}.Fetch(t.Context(), "https://shop.test/menu")
			if err != nil || status != tt.wantStatus {
				t.Errorf("Fetch = %d, %v；想要 %d", status, err, tt.wantStatus)
			}
			if !slices.Equal(methods, tt.wantMethods) {
				t.Errorf("送出的請求 = %v，想要 %v", methods, tt.wantMethods)
			}
			if closed != len(methods) {
				t.Errorf("關了 %d 個 Body，想要 %d 個", closed, len(methods))
			}
		})
	}
}

func TestHTTPFetcherTimeout(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		<-r.Context().Done() // 永遠不回應，直到請求被取消
		return nil, r.Context().Err()
	})}
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
	defer cancel()

	_, err := httpFetcher{client}.Fetch(ctx, "https://slow.test")
	if !errors.Is(err, context.DeadlineExceeded) || describe(err) != "逾時" {
		t.Errorf("Fetch 的錯誤 = %v（%s），想要包著 DeadlineExceeded", err, describe(err))
	}
}

func Example_validateURL() {
	for _, s := range []string{"https://go.dev", "go.dev", "ftp://files.test", "https://", "https://shop.test/%zz"} {
		fmt.Printf("%-24q %v\n", s, validateURL(s))
	}
	// Output:
	// "https://go.dev"         <nil>
	// "go.dev"                 只檢查 http 和 https 網址
	// "ftp://files.test"       只檢查 http 和 https 網址
	// "https://"               網址裡沒有主機名稱
	// "https://shop.test/%zz"  parse "https://shop.test/%zz": invalid URL escape "%zz"
}

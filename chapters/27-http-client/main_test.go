package main

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

func Example_getDemo() {
	getDemo()
	// Output:
	// 200 application/json
	// [{"dish":"滷肉飯","price":50},{"dish":"排骨飯","price":90}]
	// err 是 <nil>，狀態是 404 Not Found
}

func Example_urlDemo() {
	urlDemo()
	// Output:
	// http://shop.test/search?dish=%E6%8E%92%E9%AA%A8%E9%A3%AF+%E5%A4%A7&note=%E4%B8%8D%E8%A6%81%E8%BE%A3%26%E5%A4%9A%E9%A3%AF&tag=%E4%BE%BF%E7%95%B6&tag=%E5%A4%96%E5%B8%B6
	// http://shop.test/api/orders/A/7 <nil>
	// http://shop.test/orders/A%2F7
	// shop.test /search 排骨飯 2
}

func Example_jsonDemo() {
	jsonDemo()
	// Output:
	// [{滷肉飯 50} {排骨飯 90}] <nil>
	// {ID:1 Dish:排骨飯 Qty:2 Total:180} <nil>
	// HTTP 409：排骨飯 賣完了
	// 賣完了，換一道菜
	// HTTP 404：找不到這筆訂單
}

func Example_retryDemo() {
	retryDemo()
	// Output:
	// 店家收到第 1 次請求
	// 店家收到第 2 次請求
	// 店家收到第 3 次請求
	// 200 今天的便當：排骨飯
}

func Example_transportDemo() {
	transportDemo()
	// Output:
	// level=INFO msg=request method=GET url=http://shop.test/menu status=200
	// level=INFO msg=request method=GET url=http://amei:xxxxx@shop.test/orders/7 status=404
}

// TestClientTimeout：店家要 5 秒才做好，客人只等 2 秒。
// synctest 裡的時間是假的，時間到了就跳過去，測試不用真的等。
func TestClientTimeout(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		start := time.Now()
		gaveUp := make(chan time.Duration, 2)
		srv := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case <-time.After(5 * time.Second):
				io.WriteString(w, "便當好了")
			case <-r.Context().Done(): // 客人走了，店家這邊也收到通知
				gaveUp <- time.Since(start)
			}
		}))
		client := srv.Client()
		client.Timeout = 2 * time.Second

		_, err := client.Get("http://shop.test/slow")
		if !errors.Is(err, context.DeadlineExceeded) || time.Since(start) != 2*time.Second {
			t.Fatalf("Client.Timeout：%v 之後拿到 %v，want 2s 之後 DeadlineExceeded", time.Since(start), err)
		}
		if got := <-gaveUp; got != 2*time.Second {
			t.Errorf("店家在 %v 才知道客人走了，want 2s", got)
		}

		// 單一請求的期限用 context：1 秒，比 Client.Timeout 短，先到的那個生效。
		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://shop.test/slow", nil)
		if err != nil {
			t.Fatal(err)
		}
		before := time.Now()
		if _, err := client.Do(req); !errors.Is(err, context.DeadlineExceeded) || time.Since(before) != time.Second {
			t.Errorf("context：%v 之後拿到 %v，want 1s 之後 DeadlineExceeded", time.Since(before), err)
		}
		<-gaveUp
	})
}

// countConns 回傳一個記憶體裡的 server，以及它到目前為止接受了幾條連線。
func countConns(t *testing.T) (*httptest.Server, *atomic.Int64) {
	var conns atomic.Int64
	srv := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "排骨飯 90")
	}))
	srv.Config.ConnState = func(_ net.Conn, s http.ConnState) {
		if s == http.StateNew {
			conns.Add(1)
		}
	}
	return srv, &conns
}

// TestConnReuse：body 讀完並關掉，連線才會回到連線池給下一個請求用。
func TestConnReuse(t *testing.T) {
	t.Run("讀完並關掉", func(t *testing.T) {
		srv, conns := countConns(t)
		for range 5 {
			resp, err := srv.Client().Get("http://shop.test/menu")
			if err != nil {
				t.Fatal(err)
			}
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
		if n := conns.Load(); n != 1 {
			t.Errorf("5 個請求用了 %d 條連線，want 1", n)
		}
	})
	t.Run("沒有關掉", func(t *testing.T) {
		srv, conns := countConns(t)
		var open []io.Closer
		for range 5 {
			resp, err := srv.Client().Get("http://shop.test/menu")
			if err != nil {
				t.Fatal(err)
			}
			open = append(open, resp.Body) // 沒讀也沒關：這條連線一直被佔著
		}
		if n := conns.Load(); n != 5 {
			t.Errorf("5 個請求用了 %d 條連線，want 5（每次都要開新的）", n)
		}
		for _, b := range open {
			b.Close()
		}
	})
}

// TestGetWithRetry 用 synctest 檢查重試的次數和等待的時間。
func TestGetWithRetry(t *testing.T) {
	tests := []struct {
		name             string
		status           int
		deadline         time.Duration // 0 表示 ctx 沒有期限
		wantCalls        int64
		wantErr          error
		minWait, maxWait time.Duration
	}{
		// 等 1 + 2 + 4 = 7 秒，每次再加上最多一半的抖動，最多 10.5 秒
		{"一直 503：試 4 次", 503, 0, 4, nil, 7 * time.Second, 10500 * time.Millisecond},
		{"404 不重試", 404, 0, 1, nil, 0, 0},
		{"等待時 ctx 到期", 503, 2 * time.Second, 2, context.DeadlineExceeded, 2 * time.Second, 2 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var calls atomic.Int64
				srv := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					http.Error(w, http.StatusText(tt.status), tt.status)
				}))
				ctx := t.Context()
				if tt.deadline > 0 {
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(ctx, tt.deadline)
					defer cancel()
				}

				start := time.Now()
				resp, err := getWithRetry(ctx, srv.Client(), "http://shop.test/today", 4, time.Second)
				waited := time.Since(start)
				if resp != nil {
					resp.Body.Close()
				}

				if !errors.Is(err, tt.wantErr) || (tt.wantErr == nil && resp.StatusCode != tt.status) {
					t.Errorf("err = %v，want %v", err, tt.wantErr)
				}
				if calls.Load() != tt.wantCalls {
					t.Errorf("店家收到 %d 次請求，want %d", calls.Load(), tt.wantCalls)
				}
				if waited < tt.minWait || waited > tt.maxWait {
					t.Errorf("總共等了 %v，want %v 到 %v", waited, tt.minWait, tt.maxWait)
				}
			})
		})
	}
}

// TestDo 檢查 do 處理不是 JSON 的錯誤內容，以及太大的回應。
func TestDo(t *testing.T) {
	h := http.NewServeMux()
	h.HandleFunc("GET /gateway", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "<html>upstream down</html>", http.StatusBadGateway)
	})
	h.HandleFunc("GET /huge", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `["`+strings.Repeat("滷", maxResponse)+`"]`)
	})
	client := &http.Client{Transport: handlerTransport{h}}
	ctx := t.Context()

	_, err := getJSON[[]string](ctx, client, "http://shop.test/gateway")
	if apiErr, ok := errors.AsType[*APIError](err); !ok || apiErr.Status != 502 || apiErr.Message != "Bad Gateway" {
		t.Errorf("不是 JSON 的錯誤內容：err = %v，want HTTP 502：Bad Gateway", err)
	}

	// 只讀前 1 MB，JSON 在中間被截斷，所以解析失敗，而不是把整個回應讀進記憶體。
	if _, err := getJSON[[]string](ctx, client, "http://shop.test/huge"); err == nil {
		t.Errorf("超過 %d bytes 的回應應該解析失敗", maxResponse)
	}
}

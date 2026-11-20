//go:build solution

package bus

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// BuildURL 用 url.PathEscape 跳脫站名：它會把 / 變成 %2F，站名才不會被拆成兩段路徑。
func BuildURL(base, stop string, routes []string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	if u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("base 要是完整的網址：%q", base)
	}
	s := strings.TrimSuffix(u.String(), "/") + "/stops/" + url.PathEscape(stop) + "/arrivals"
	if len(routes) > 0 {
		s += "?" + url.Values{"route": routes}.Encode()
	}
	return s, nil
}

// closeBody 把 body 讀完再關掉，連線才能回到連線池。
func closeBody(resp *http.Response) {
	io.Copy(io.Discard, io.LimitReader(resp.Body, maxBody))
	resp.Body.Close()
}

func GetArrivals(ctx context.Context, client *http.Client, u string) ([]Arrival, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeBody(resp)
	body := io.LimitReader(resp.Body, maxBody)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var e struct {
			Message string `json:"message"`
		}
		json.UnmarshalRead(body, &e) // 解不開就讓 Message 留空
		return nil, &StatusError{Code: resp.StatusCode, Message: e.Message}
	}
	var arrivals []Arrival
	if err := json.UnmarshalRead(body, &arrivals); err != nil {
		return nil, err
	}
	return arrivals, nil
}

// RoundTrip 先用 Clone 複製一份請求，改的是複本。
func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r := req.Clone(req.Context())
	r.Header.Set("Authorization", "Bearer "+t.Token)
	return t.Next.RoundTrip(r)
}

func RetryAfter(h http.Header, now time.Time) (time.Duration, bool) {
	v := h.Get("Retry-After")
	if v == "" {
		return 0, false
	}
	if secs, err := strconv.Atoi(v); err == nil {
		if secs < 0 {
			return 0, false
		}
		return time.Duration(secs) * time.Second, true
	}
	at, err := http.ParseTime(v)
	if err != nil {
		return 0, false
	}
	return max(at.Sub(now), 0), true
}

// Get 只在讀寫 entries 時拿鎖，送請求的時候不拿：一個慢的網址不會卡住其他的。
func (c *Cache) Get(ctx context.Context, client *http.Client, u string) ([]byte, error) {
	c.mu.Lock()
	cached, have := c.entries[u]
	c.mu.Unlock()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	if have {
		req.Header.Set("If-None-Match", cached.etag)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeBody(resp)

	switch {
	case resp.StatusCode == http.StatusNotModified && have:
		return bytes.Clone(cached.body), nil
	case resp.StatusCode == http.StatusOK:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if etag := resp.Header.Get("ETag"); etag != "" {
			c.mu.Lock()
			if c.entries == nil {
				c.entries = make(map[string]entry)
			}
			c.entries[u] = entry{etag: etag, body: body}
			c.mu.Unlock()
		}
		return bytes.Clone(body), nil
	}
	return nil, &StatusError{Code: resp.StatusCode}
}

// Healthy 把 defer 放在只處理一個網址的函式裡，每一圈結束就關掉那一個 body。
func Healthy(ctx context.Context, client *http.Client, urls []string) int {
	n := 0
	for _, u := range urls {
		if ok(ctx, client, u) {
			n++
		}
	}
	return n
}

func ok(ctx context.Context, client *http.Client, u string) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer closeBody(resp)
	return resp.StatusCode == http.StatusOK
}

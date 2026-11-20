//go:build !solution

package bus

import (
	"context"
	"net/http"
	"time"
)

// BuildURL 組出查詢某一站到站時間的網址：base 後面接 /stops/{站名}/arrivals，
// 每個 route 都是一個 route 查詢參數。例如 BuildURL("https://bus.test/v1", "市政府", []string{"307", "藍5"}) 是
//
//	https://bus.test/v1/stops/%E5%B8%82%E6%94%BF%E5%BA%9C/arrivals?route=307&route=%E8%97%8D5
//
// 站名可能有 /、空白和括號，要整個當成一段路徑來跳脫。base 結尾有沒有 / 都可以；
// routes 是空的時候不要有 ?。base 解析失敗，或缺少 scheme、主機時，回傳錯誤。
func BuildURL(base, stop string, routes []string) (string, error) {
	// TODO: 實作
	return "", nil
}

// GetArrivals 對 u 送出 GET，回傳到站資訊：
//   - 請求要帶著 ctx，標頭要有 Accept: application/json。
//   - 2xx：內容是 JSON 陣列，解成 []Arrival。
//   - 其他狀態碼：回傳 *StatusError；內容是 {"message":"..."} 時 Message 是那個字串，否則留空。
//   - 不管結果如何都要關掉 body，而且最多只讀 maxBody 個 byte。
func GetArrivals(ctx context.Context, client *http.Client, u string) ([]Arrival, error) {
	// TODO: 實作
	return nil, nil
}

// RoundTrip 幫請求加上 Authorization: Bearer <Token> 標頭，交給 t.Next 送出。
// RoundTripper 不能修改傳進來的請求：呼叫的一方可能還要用它，例如重試時再送一次。
func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// TODO: 實作
	return nil, nil
}

// RetryAfter 解讀回應的 Retry-After 標頭，回傳要等多久。
// 標頭可以是秒數（"120"），也可以是 HTTP 日期（"Sat, 21 Nov 2026 08:00:00 GMT"），日期要減掉 now。
// 沒有這個標頭、格式不對、秒數是負的時候 ok 是 false；日期已經過了，回傳 0 和 true。
func RetryAfter(h http.Header, now time.Time) (wait time.Duration, ok bool) {
	// TODO: 實作
	return 0, false
}

// Get 對 u 送出 GET，用 ETag 避免重複下載沒變的內容：
//   - 第一次正常取得；回應有 ETag 標頭時，把內容和 ETag 記下來。
//   - 之後再取同一個網址，帶上 If-None-Match: <記下的 ETag>。
//     對方回 304 Not Modified 時回傳記下的內容；回 200 時回傳新的內容，並換掉記下的內容和 ETag。
//   - 其他狀態碼回傳 *StatusError（Message 留空），記下的內容不變。
//
// 可以被很多 goroutine 同時呼叫。
func (c *Cache) Get(ctx context.Context, client *http.Client, u string) ([]byte, error) {
	// TODO: 實作
	return nil, nil
}

// Healthy 依序對每個網址送出 GET，回傳回應是 200 的數量；送不出去的網址不算。
//
// 找錯：網址一多，同時開著的連線就越來越多，要等到函式結束才一起放掉。
func Healthy(ctx context.Context, client *http.Client, urls []string) int {
	n := 0
	for _, u := range urls {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
		if err != nil {
			continue
		}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			n++
		}
	}
	return n
}

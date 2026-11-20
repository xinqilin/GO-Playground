// Package bus 是第 27 章「HTTP 用戶端」的練習題：查詢公車到站時間的 client。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/27-http-client/
//
// 只想跑某一題時加上 -run：
//
//	go test ./exercises/27-http-client/ -run TestGetArrivals
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/27-http-client/
//
// 測試都不經過網路：請求交給測試裡寫好的假站牌 API。
package bus

import (
	"fmt"
	"net/http"
	"sync"
)

// Arrival 是某一路公車還有幾分鐘到站。
type Arrival struct {
	Route   string `json:"route"`
	Minutes int    `json:"minutes"`
}

// StatusError 是對方回了 2xx 以外的狀態碼。
type StatusError struct {
	Code    int
	Message string // 對方 JSON 內容裡的 message；沒有的話是空字串
}

func (e *StatusError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("HTTP %d", e.Code)
	}
	return fmt.Sprintf("HTTP %d：%s", e.Code, e.Message)
}

// maxBody 是練習 2 最多願意讀的回應大小。
const maxBody = 1 << 20

// AuthTransport 是練習 3 的 RoundTripper：幫每個請求加上 API 金鑰，再交給 Next。
type AuthTransport struct {
	Token string
	Next  http.RoundTripper // 不能是 nil
}

// Cache 是練習 5 的快取：依網址記下內容和 ETag。零值就能用。
type Cache struct {
	mu      sync.Mutex
	entries map[string]entry
}

type entry struct {
	etag string
	body []byte
}

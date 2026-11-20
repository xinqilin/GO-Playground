//go:build !solution

package ticket

import (
	"context"
	"encoding/json/v2"
	"net/http"
)

// QueryInt 讀網址查詢字串裡的整數參數，例如 /tickets?limit=20 的 limit。
// 沒有這個參數（或是空字串）時回傳 def；不是整數，或不在 min 到 max 之間（含）時回傳錯誤，
// 錯誤訊息要提到參數的名稱。
func QueryInt(r *http.Request, name string, def, min, max int) (int, error) {
	// TODO: 實作
	return 0, nil
}

// RequireKey 是 middleware：請求的 X-Staff-Key 標頭要等於 key，才交給 next；
// 否則回 401，內容是 JSON {"error":"需要員工金鑰"}，Content-Type 是 application/json。
// key 是空字串時一律拒絕（設定忘了填，不能變成誰都能進來）。
// 比對金鑰要用固定時間的比較，不能讓回應時間透露金鑰猜對了幾個字。
func RequireKey(key string, next http.Handler) http.Handler {
	// TODO: 實作
	return next
}

// RequestID 是 middleware：每個請求都有一個 request id。
// 請求的 X-Request-ID 標頭格式合法（1 到 64 個英文字母、數字或 -）就沿用，否則產生一個新的 UUID。
// id 寫進回應的 X-Request-ID 標頭，並放進請求的 context，next 用 FromContext 取得。
func RequestID(next http.Handler) http.Handler {
	// TODO: 實作
	return next
}

// FromContext 回傳 RequestID 放進 ctx 的 request id；沒有時回傳空字串。
func FromContext(ctx context.Context) string {
	// TODO: 實作
	return ""
}

// MaxInFlight 是 middleware：同時最多 n 個請求在 next 裡執行。
// 已經有 n 個時，新的請求不等待，直接回 503，Retry-After 標頭是 "1"，
// 內容是 JSON {"error":"太忙了，請稍後再試"}。
func MaxInFlight(n int, next http.Handler) http.Handler {
	// TODO: 實作
	return next
}

// NewServer 回傳叫號機的 API：
//   - POST /tickets：取號。回 201，Location 是 /tickets/{號碼}，內容是 {"number":號碼,"ahead":前面幾人}。
//   - GET /tickets/{number}：查詢。回 200 和 {"number":號碼,"ahead":前面幾人}；
//     number 不是正整數回 400，不在等候中回 404，內容都是 {"error":"..."}。
//   - POST /next：叫號。回 200 和 {"number":號碼}；沒有人在等時回 204，沒有內容。
//
// 有內容的回應，Content-Type 都是 application/json。方法不對時回 405（ServeMux 會處理）。
func NewServer(q *Queue) http.Handler {
	// TODO: 實作
	return http.NotFoundHandler()
}

// Rename 把店名改成請求內容 {"name":"..."} 裡的 name，回 200 和 {"name":"新店名"}。
// 內容解不開或 name 是空的時候，回 400 和 "需要 name"。
//
// 找錯：送了空的 name，回應是 400，但店名還是被改成空的，回應的內容也怪怪的。
func (s *Settings) Rename(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name string `json:"name"`
	}
	if err := json.UnmarshalRead(r.Body, &in); err != nil || in.Name == "" {
		http.Error(w, "需要 name", http.StatusBadRequest)
	}
	s.mu.Lock()
	s.name = in.Name
	s.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.MarshalWrite(w, map[string]string{"name": in.Name})
}

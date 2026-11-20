package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
)

// send 用 httptest 送一個請求給 h，印出狀態碼、指定的標頭，以及內容。不用開 port，也不用網路。
// body 不是空的時候，Content-Type 是 application/json。
func send(h http.Handler, method, target, body string, headers ...string) {
	contentType := ""
	if body != "" {
		contentType = "application/json"
	}
	sendAs(h, contentType, method, target, body, headers...)
}

// sendAs 和 send 一樣，但可以指定請求的 Content-Type。
func sendAs(h http.Handler, contentType, method, target, body string, headers ...string) {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	res := rec.Result() // 標頭是第一次寫出內容時的樣子，和真的 server 送出去的一樣
	fmt.Printf("%s %s → %d", method, target, res.StatusCode)
	for _, k := range headers {
		if v := res.Header.Get(k); v != "" {
			fmt.Printf("，%s: %s", k, v)
		}
	}
	fmt.Println()
	if b, _ := io.ReadAll(res.Body); len(b) > 0 {
		fmt.Println(strings.TrimRight(string(b), "\n"))
	}
}

// hello 是最簡單的 handler：簽名符合的函式，用 http.HandlerFunc 轉成 http.Handler。
func hello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "客人"
	}
	fmt.Fprintf(w, "歡迎光臨，%s\n", name)
}

// counter 自己實作 ServeHTTP。每個請求在自己的 goroutine 裡執行，共用的欄位要用 atomic 或 Mutex 保護。
type counter struct {
	n atomic.Int64
}

func (c *counter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "您是今天第 %d 位客人\n", c.n.Add(1))
}

func handlerDemo() {
	send(http.HandlerFunc(hello), "GET", "/hello?name=Amei", "", "Content-Type")
	send(http.HandlerFunc(hello), "GET", "/hello", "")

	c := &counter{}
	send(c, "GET", "/", "")
	send(c, "GET", "/", "")
}

// lateHeader 先寫內容，才設標頭和狀態碼。第一次 Write 時標頭和狀態碼就送出去了，之後的設定都沒有作用。
func lateHeader(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, `{"ok":true}`)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// rightOrder 照順序：標頭、狀態碼、內容。
func rightOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, `{"ok":true}`)
}

func writeOrderDemo() {
	send(http.HandlerFunc(lateHeader), "POST", "/late", "", "Content-Type")
	send(http.HandlerFunc(rightOrder), "POST", "/right", "", "Content-Type")
}

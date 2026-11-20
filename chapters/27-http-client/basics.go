package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// newClient 回傳正式用的 client：整個程式共用一個，它的 Transport 管著連線池。
// 從 DefaultTransport 複製一份再調整，保留它的連線逾時、Proxy、HTTP/2 等預設值。
func newClient() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConnsPerHost = 10 // 預設是 2：同一個主機最多留 2 條閒置的連線
	return &http.Client{
		Timeout:   10 * time.Second, // 從送出到讀完 body 的總時間
		Transport: t,
	}
}

// newDemoClient 回傳示範用的 client：設定和正式的一樣，只是請求交給 fakeShop，不經過網路。
func newDemoClient() *http.Client {
	c := newClient()
	c.Transport = handlerTransport{fakeShop()}
	return c
}

// getDemo 送出 GET，讀狀態碼、標頭和內容。注意 404 不是 error。
func getDemo() {
	client := newDemoClient()

	resp, err := client.Get("http://shop.test/menu")
	if err != nil {
		fmt.Println("送不出去：", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("讀內容失敗：", err)
		return
	}
	fmt.Println(resp.StatusCode, resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	resp2, err := client.Get("http://shop.test/orders/99")
	if err != nil {
		fmt.Println("送不出去：", err)
		return
	}
	defer resp2.Body.Close()
	fmt.Printf("err 是 %v，狀態是 %s\n", err, resp2.Status)
}

package main

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// APIError 是對方回了 4xx、5xx 的錯誤。Message 是對方 JSON 裡的 error，沒有的話是狀態的文字。
type APIError struct {
	Status  int
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("HTTP %d：%s", e.Status, e.Message)
}

// maxResponse 是願意讀的回應大小上限。對方再怎麼出錯，也不能讓我們的記憶體被吃光。
const maxResponse = 1 << 20

// do 送出請求；2xx 時把內容解到 out（可以是 nil），其他狀態碼回傳 *APIError。
// 不管哪一種結果，都把 body 讀完並關掉，連線才能還回去重用。
func do(client *http.Client, req *http.Request, out any) error {
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		io.Copy(io.Discard, io.LimitReader(resp.Body, maxResponse))
		resp.Body.Close()
	}()
	body := io.LimitReader(resp.Body, maxResponse)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		apiErr := &APIError{Status: resp.StatusCode, Message: http.StatusText(resp.StatusCode)}
		var e struct {
			Error string `json:"error"`
		}
		if json.UnmarshalRead(body, &e) == nil && e.Error != "" {
			apiErr.Message = e.Error
		}
		return apiErr
	}
	if out == nil {
		return nil
	}
	return json.UnmarshalRead(body, out)
}

// getJSON 送出 GET，把 JSON 回應解成 T。
func getJSON[T any](ctx context.Context, client *http.Client, u string) (T, error) {
	var out T
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return out, err
	}
	req.Header.Set("Accept", "application/json")
	err = do(client, req, &out)
	return out, err
}

// postJSON 把 in 編成 JSON 送出 POST，把 JSON 回應解成 T。
func postJSON[T any](ctx context.Context, client *http.Client, u string, in any) (T, error) {
	var out T
	b, err := json.Marshal(in)
	if err != nil {
		return out, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(b))
	if err != nil {
		return out, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	err = do(client, req, &out)
	return out, err
}

func jsonDemo() {
	ctx := context.Background()
	client := newDemoClient()

	menu, err := getJSON[[]MenuItem](ctx, client, "http://shop.test/menu")
	fmt.Println(menu, err)

	type orderIn struct {
		Dish string `json:"dish"`
		Qty  int    `json:"qty"`
	}
	o, err := postJSON[Order](ctx, client, "http://shop.test/orders", orderIn{"排骨飯", 2})
	fmt.Printf("%+v %v\n", o, err)

	_, err = postJSON[Order](ctx, client, "http://shop.test/orders", orderIn{"排骨飯", 2})
	fmt.Println(err)
	if apiErr, ok := errors.AsType[*APIError](err); ok && apiErr.Status == http.StatusConflict {
		fmt.Println("賣完了，換一道菜")
	}

	_, err = getJSON[Order](ctx, client, "http://shop.test/orders/99")
	fmt.Println(err)
}

// urlDemo 網址用 net/url 組，不要自己拼字串：中文、空白、& 這些字元都要跳脫。
func urlDemo() {
	q := url.Values{}
	q.Set("dish", "排骨飯 大")
	q.Set("note", "不要辣&多飯")
	q.Add("tag", "便當")
	q.Add("tag", "外帶")
	fmt.Println("http://shop.test/search?" + q.Encode())

	u, err := url.JoinPath("http://shop.test/api/", "orders", "A/7")
	fmt.Println(u, err)
	fmt.Println("http://shop.test/orders/" + url.PathEscape("A/7"))

	parsed, err := url.Parse("http://shop.test/search?dish=%E6%8E%92%E9%AA%A8%E9%A3%AF&page=2")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(parsed.Host, parsed.Path, parsed.Query().Get("dish"), parsed.Query().Get("page"))
}

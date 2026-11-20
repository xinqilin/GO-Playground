package main

import (
	"context"
	"encoding/json/v2"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

func Example_handlerDemo() {
	handlerDemo()
	// Output:
	// GET /hello?name=Amei → 200，Content-Type: text/plain; charset=utf-8
	// 歡迎光臨，Amei
	// GET /hello → 200
	// 歡迎光臨，客人
	// GET / → 200
	// 您是今天第 1 位客人
	// GET / → 200
	// 您是今天第 2 位客人
}

func Example_writeOrderDemo() {
	writeOrderDemo()
	// Output:
	// POST /late → 200，Content-Type: text/plain; charset=utf-8
	// {"ok":true}
	// POST /right → 201，Content-Type: application/json
	// {"ok":true}
}

func Example_routeDemo() {
	routeDemo()
	// Output:
	// GET /orders/42 → 200
	// pattern="GET /orders/{id}" id="42" path=""
	// GET /orders/latest → 200
	// pattern="GET /orders/latest" id="" path=""
	// HEAD /orders/42 → 200
	// pattern="GET /orders/{id}" id="42" path=""
	// DELETE /orders/42 → 405，Allow: GET, HEAD
	// Method Not Allowed
	// GET /menu/today → 200
	// pattern="GET /menu/" id="" path=""
	// GET /menu → 307，Location: /menu/
	// <a href="/menu/">Temporary Redirect</a>.
	// GET /files/2026/11/menu.pdf → 200
	// pattern="GET /files/{path...}" id="" path="2026/11/menu.pdf"
	// GET / → 200
	// pattern="GET /{$}" id="" path=""
	// GET /about → 404
	// 404 page not found
}

func Example_apiDemo() {
	apiDemo()
	// Output:
	// GET /menu → 200
	// [{"dish":"滷肉飯","price":50,"left":20},{"dish":"排骨飯","price":90,"left":5},{"dish":"雞腿飯","price":100,"left":0}]
	// POST /orders → 201，Location: /orders/1
	// {"id":1,"dish":"排骨飯","qty":2,"total":180}
	// GET /orders/1 → 200
	// {"id":1,"dish":"排骨飯","qty":2,"total":180}
	// POST /orders → 409
	// {"error":"排骨飯 賣完了，剩 3 份"}
	// POST /orders → 409
	// {"error":"雞腿飯 賣完了，剩 0 份"}
	// GET /orders/99 → 404
	// {"error":"找不到這筆訂單"}
	// GET /orders/abc → 400
	// {"error":"訂單編號要是數字"}
}

func Example_decodeDemo() {
	decodeDemo()
	// Output:
	// POST /orders → 400
	// {"error":"qty 要在 1 到 20 之間"}
	// POST /orders → 400
	// {"error":"欄位 /qty 的型別不對"}
	// POST /orders → 400
	// {"error":"不認得的欄位 /spicy"}
	// POST /orders → 400
	// {"error":"JSON 格式不對"}
	// POST /orders → 413
	// {"error":"內容太大"}
	// POST /orders → 415
	// {"error":"Content-Type 要是 application/json"}
}

func Example_middlewareDemo() {
	middlewareDemo()
	// Output:
	// level=INFO msg=request method=GET path=/healthz pattern="GET /healthz" status=200 bytes=3
	// GET /healthz → 200
	// ok
	// level=INFO msg=request method=GET path=/orders/7 pattern="GET /orders/{id}" status=404 bytes=33
	// GET /orders/7 → 404
	// {"error":"找不到這筆訂單"}
	// level=ERROR msg=panic path=/oops value="assignment to entry in nil map"
	// level=INFO msg=request method=GET path=/oops pattern="GET /oops" status=500 bytes=33
	// GET /oops → 500
	// {"error":"伺服器發生錯誤"}
}

// TestEndToEnd 用 Go 1.27 的 httptest.NewTestServer：真的 http.Client、真的 HTTP，但網路在記憶體裡，不用開 port。
func TestEndToEnd(t *testing.T) {
	a := &api{shop: newShop(), logger: slog.New(slog.DiscardHandler)}
	srv := httptest.NewTestServer(t, a.handler())
	client := srv.Client() // 一定要用這個 client，它才會把請求送到記憶體裡的 server

	res, err := client.Post("http://shop.test/orders", "application/json", strings.NewReader(`{"dish":"滷肉飯","qty":3}`))
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	loc := res.Header.Get("Location")
	if res.StatusCode != http.StatusCreated || loc != "/orders/1" {
		t.Fatalf("POST /orders = %d，Location %q；want 201、/orders/1", res.StatusCode, loc)
	}

	res, err = client.Get("http://shop.test" + loc)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var got Order
	if err := json.UnmarshalRead(res.Body, &got); err != nil {
		t.Fatal(err)
	}
	if want := (Order{ID: 1, Dish: "滷肉飯", Qty: 3, Total: 150}); got != want {
		t.Errorf("GET %s = %+v，want %+v", loc, got, want)
	}
}

// TestConcurrentOrders 20 個人同時搶 5 份排骨飯：剛好 5 個成功，其他人拿到 409。用 -race 跑。
func TestConcurrentOrders(t *testing.T) {
	h := (&api{shop: newShop(), logger: slog.New(slog.DiscardHandler)}).routes()
	var mu sync.Mutex
	codes := map[int]int{}
	var wg sync.WaitGroup
	for range 20 {
		wg.Go(func() {
			req := httptest.NewRequest("POST", "/orders", strings.NewReader(`{"dish":"排骨飯","qty":1}`))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			mu.Lock()
			codes[rec.Code]++
			mu.Unlock()
		})
	}
	wg.Wait()
	if codes[http.StatusCreated] != 5 || codes[http.StatusConflict] != 15 {
		t.Errorf("狀態碼統計 %v，want 201 五次、409 十五次", codes)
	}
}

// TestShutdownWaits Shutdown 會等進行中的請求做完；synctest 的假時間讓「做 3 秒」不用真的等。
func TestShutdownWaits(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		started := make(chan struct{})
		slow := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			close(started)
			time.Sleep(3 * time.Second) // 正在做便當
			io.WriteString(w, "便當好了")
		})
		srv := httptest.NewTestServer(t, slow)
		client := srv.Client()

		got := make(chan string)
		go func() {
			res, err := client.Get("http://shop.test/")
			if err != nil {
				got <- "錯誤：" + err.Error()
				return
			}
			defer res.Body.Close()
			b, _ := io.ReadAll(res.Body)
			got <- string(b)
		}()

		<-started
		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Config.Shutdown(ctx); err != nil {
			t.Fatalf("Shutdown：%v", err)
		}
		if waited := time.Since(start); waited < 3*time.Second {
			t.Errorf("Shutdown 只等了 %v，進行中的請求還沒做完", waited)
		}
		if body := <-got; body != "便當好了" {
			t.Errorf("客戶端拿到 %q，want 便當好了", body)
		}
	})
}

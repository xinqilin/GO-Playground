package main

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"strconv"
	"strings"
)

// api 是訂單 API。需要的東西都放在欄位裡，handler 是它的方法。
type api struct {
	shop   *shop
	logger *slog.Logger
}

func (a *api) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /menu", a.handle(a.listMenu))
	mux.HandleFunc("POST /orders", a.handle(a.createOrder))
	mux.HandleFunc("GET /orders/{id}", a.handle(a.getOrder))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "ok\n")
	})
	return mux
}

// apiFunc 是會回傳 error 的 handler。錯誤怎麼變成 HTTP 回應，集中在 handle 一個地方決定。
type apiFunc func(w http.ResponseWriter, r *http.Request) error

func (a *api) handle(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err == nil {
			return
		}
		status, msg := statusOf(err)
		if status == http.StatusInternalServerError {
			// 內部錯誤的細節只記在 log，不告訴客戶端。
			a.logger.ErrorContext(r.Context(), "內部錯誤", "path", r.URL.Path, "err", err)
		}
		writeJSON(w, status, errorBody{Error: msg})
	}
}

// errorBody 是所有錯誤回應共同的格式。
type errorBody struct {
	Error string `json:"error"`
}

// httpError 是可以原樣告訴客戶端的錯誤，例如輸入不合法。
type httpError struct {
	status int
	msg    string
}

func (e *httpError) Error() string { return e.msg }

// statusOf 決定錯誤對應的狀態碼和要給客戶端看的訊息。
func statusOf(err error) (int, string) {
	if he, ok := errors.AsType[*httpError](err); ok {
		return he.status, he.msg
	}
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, ErrSoldOut):
		return http.StatusConflict, err.Error()
	case errors.Is(err, ErrUnknownDish):
		return http.StatusBadRequest, err.Error()
	}
	return http.StatusInternalServerError, "伺服器發生錯誤"
}

// writeJSON 先編碼，成功了才寫標頭和狀態碼；編碼失敗時什麼都還沒寫，呼叫的一方還能回 500。
func writeJSON(w http.ResponseWriter, status int, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b) // 寫失敗通常是客戶端已經斷線，沒有什麼能做的
	return nil
}

// maxBody 是請求內容的上限。沒有上限的話，一個很大的請求就能吃光記憶體。
const maxBody = 1 << 20

// decodeJSON 檢查 Content-Type、限制大小，嚴格地解碼，並把各種失敗轉成給客戶端看的訊息。
func decodeJSON(w http.ResponseWriter, r *http.Request, v any) error {
	if mt, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type")); mt != "application/json" {
		return &httpError{http.StatusUnsupportedMediaType, "Content-Type 要是 application/json"}
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxBody)
	err := json.UnmarshalRead(r.Body, v, json.RejectUnknownMembers(true))
	if err == nil {
		return nil
	}
	if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
		return &httpError{http.StatusRequestEntityTooLarge, "內容太大"}
	}
	if errors.Is(err, json.ErrUnknownName) {
		se, _ := errors.AsType[*json.SemanticError](err)
		return &httpError{http.StatusBadRequest, "不認得的欄位 " + string(se.JSONPointer)}
	}
	if se, ok := errors.AsType[*json.SemanticError](err); ok && se.JSONPointer != "" {
		return &httpError{http.StatusBadRequest, "欄位 " + string(se.JSONPointer) + " 的型別不對"}
	}
	if _, ok := errors.AsType[*jsontext.SyntacticError](err); ok || errors.Is(err, io.ErrUnexpectedEOF) {
		return &httpError{http.StatusBadRequest, "JSON 格式不對"}
	}
	return err // 讀取 body 時的其他錯誤，例如連線中斷
}

func (a *api) listMenu(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, http.StatusOK, a.shop.menu())
}

func (a *api) createOrder(w http.ResponseWriter, r *http.Request) error {
	var in struct {
		Dish string `json:"dish"`
		Qty  int    `json:"qty"`
	}
	if err := decodeJSON(w, r, &in); err != nil {
		return err
	}
	if in.Qty < 1 || in.Qty > 20 {
		return &httpError{http.StatusBadRequest, "qty 要在 1 到 20 之間"}
	}
	o, err := a.shop.place(in.Dish, in.Qty)
	if err != nil {
		return err
	}
	w.Header().Set("Location", "/orders/"+strconv.Itoa(o.ID))
	return writeJSON(w, http.StatusCreated, o)
}

func (a *api) getOrder(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return &httpError{http.StatusBadRequest, "訂單編號要是數字"}
	}
	o, err := a.shop.find(id)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, o)
}

// apiDemo 依序送出幾個請求：看菜單、點餐、查訂單，以及幾種店家這邊的錯誤。
func apiDemo() {
	h := (&api{shop: newShop(), logger: slog.New(slog.DiscardHandler)}).routes()
	send(h, "GET", "/menu", "")
	send(h, "POST", "/orders", `{"dish":"排骨飯","qty":2}`, "Location")
	send(h, "GET", "/orders/1", "")
	send(h, "POST", "/orders", `{"dish":"排骨飯","qty":4}`)
	send(h, "POST", "/orders", `{"dish":"雞腿飯","qty":1}`)
	send(h, "GET", "/orders/99", "")
	send(h, "GET", "/orders/abc", "")
}

// decodeDemo 客戶端送來的內容有各種錯法，每一種都要回 4xx 和看得懂的訊息。
func decodeDemo() {
	h := (&api{shop: newShop(), logger: slog.New(slog.DiscardHandler)}).routes()
	send(h, "POST", "/orders", `{"dish":"滷肉飯","qty":0}`)
	send(h, "POST", "/orders", `{"dish":"滷肉飯","qty":"二"}`)
	send(h, "POST", "/orders", `{"dish":"滷肉飯","qty":2,"spicy":true}`)
	send(h, "POST", "/orders", `{"dish":"滷肉飯",`)
	send(h, "POST", "/orders", `{"dish":"`+strings.Repeat("滷", maxBody/3)+`"}`)
	sendAs(h, "text/plain", "POST", "/orders", `{"dish":"滷肉飯","qty":2}`)
}

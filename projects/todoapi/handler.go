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
)

// api 是 HTTP 這一層：把請求轉成 service 的呼叫，再把結果和錯誤轉成 HTTP 回應。
type api struct {
	svc    *service
	logger *slog.Logger
}

func (a *api) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos", a.handle(a.list))
	mux.HandleFunc("POST /todos", a.handle(a.create))
	mux.HandleFunc("GET /todos/{id}", a.handle(a.get))
	mux.HandleFunc("PATCH /todos/{id}", a.handle(a.update))
	mux.HandleFunc("DELETE /todos/{id}", a.handle(a.delete))
	// logRequests 在最外層，才記得到 recoverPanic 把 panic 改成的 500。
	return logRequests(a.logger, recoverPanic(a.logger, mux))
}

// apiFunc 是會回傳 error 的 handler。錯誤怎麼變成回應，集中在 handle 決定。
type apiFunc func(w http.ResponseWriter, r *http.Request) error

func (a *api) handle(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err == nil {
			return
		}
		status, body := errorResponse(err)
		if status == http.StatusInternalServerError {
			// 內部錯誤的細節只記在 log，不告訴客戶端。
			a.logger.ErrorContext(r.Context(), "內部錯誤", "path", r.URL.Path, "err", err)
		}
		writeJSON(w, status, body)
	}
}

// httpError 是可以原樣告訴客戶端的錯誤，例如網址裡的 id 不是數字。
type httpError struct {
	status int
	msg    string
}

func (e *httpError) Error() string { return e.msg }

// errorBody 是所有錯誤回應共同的格式；fields 只在輸入不合法時出現。
type errorBody struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitzero"`
}

// errorResponse 決定錯誤對應的狀態碼和回應內容。
func errorResponse(err error) (int, errorBody) {
	if he, ok := errors.AsType[*httpError](err); ok {
		return he.status, errorBody{Error: he.msg}
	}
	if ve, ok := errors.AsType[ValidationError](err); ok {
		return http.StatusUnprocessableEntity, errorBody{Error: "輸入不合法", Fields: ve}
	}
	if errors.Is(err, ErrNotFound) {
		return http.StatusNotFound, errorBody{Error: err.Error()}
	}
	return http.StatusInternalServerError, errorBody{Error: "伺服器發生錯誤"}
}

// writeJSON 先編碼，成功了才寫狀態碼。Deterministic 讓 map 的 key 依序輸出（json/v2 預設不排序）。
func writeJSON(w http.ResponseWriter, status int, v any) error {
	b, err := json.Marshal(v, json.Deterministic(true))
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)
	return nil
}

const maxBody = 1 << 20

// decodeJSON 檢查 Content-Type、限制大小、拒絕不認得的欄位，並把失敗轉成給客戶端看的訊息（26 章）。
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
	if se, ok := errors.AsType[*json.SemanticError](err); ok && se.JSONPointer != "" {
		if errors.Is(err, json.ErrUnknownName) {
			return &httpError{http.StatusBadRequest, "不認得的欄位 " + string(se.JSONPointer)}
		}
		return &httpError{http.StatusBadRequest, "欄位 " + string(se.JSONPointer) + " 的型別不對"}
	}
	if _, ok := errors.AsType[*jsontext.SyntacticError](err); ok || errors.Is(err, io.ErrUnexpectedEOF) {
		return &httpError{http.StatusBadRequest, "JSON 格式不對"}
	}
	return err
}

// pathID 取出網址裡的 {id}。
func pathID(r *http.Request) (int, error) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		return 0, &httpError{http.StatusBadRequest, "id 要是正整數"}
	}
	return id, nil
}

func (a *api) list(w http.ResponseWriter, r *http.Request) error {
	var done *bool
	if v := r.URL.Query().Get("done"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return &httpError{http.StatusBadRequest, "done 要是 true 或 false"}
		}
		done = &b
	}
	return writeJSON(w, http.StatusOK, a.svc.List(done))
}

func (a *api) create(w http.ResponseWriter, r *http.Request) error {
	var in NewTodo
	if err := decodeJSON(w, r, &in); err != nil {
		return err
	}
	t, err := a.svc.Create(in)
	if err != nil {
		return err
	}
	w.Header().Set("Location", "/todos/"+strconv.Itoa(t.ID))
	return writeJSON(w, http.StatusCreated, t)
}

func (a *api) get(w http.ResponseWriter, r *http.Request) error {
	id, err := pathID(r)
	if err != nil {
		return err
	}
	t, err := a.svc.Get(id)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, t)
}

func (a *api) update(w http.ResponseWriter, r *http.Request) error {
	id, err := pathID(r)
	if err != nil {
		return err
	}
	var p Patch
	if err := decodeJSON(w, r, &p); err != nil {
		return err
	}
	t, err := a.svc.Update(id, p)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, t)
}

func (a *api) delete(w http.ResponseWriter, r *http.Request) error {
	id, err := pathID(r)
	if err != nil {
		return err
	}
	if err := a.svc.Delete(id); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

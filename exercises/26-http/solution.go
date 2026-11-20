//go:build solution

package ticket

import (
	"context"
	"crypto/subtle"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"uuid"
)

func QueryInt(r *http.Request, name string, def, min, max int) (int, error) {
	s := r.URL.Query().Get(name)
	if s == "" {
		return def, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < min || n > max {
		return 0, fmt.Errorf("%s 要是 %d 到 %d 之間的整數", name, min, max)
	}
	return n, nil
}

// writeJSON 設好標頭和狀態碼，再寫出 v。
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.MarshalWrite(w, v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// RequireKey 用 subtle.ConstantTimeCompare：不管哪一個位元組不同，比較花的時間都一樣。
func RequireKey(key string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get("X-Staff-Key")
		if key == "" || subtle.ConstantTimeCompare([]byte(got), []byte(key)) != 1 {
			writeError(w, http.StatusUnauthorized, "需要員工金鑰")
			return
		}
		next.ServeHTTP(w, r)
	})
}

type requestIDKey struct{}

var validID = regexp.MustCompile(`^[A-Za-z0-9-]{1,64}$`)

// RequestID 不相信客戶端給的值，格式不對就換掉：它會被寫進 log 和回應標頭。
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if !validID.MatchString(id) {
			id = uuid.NewV7().String()
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDKey{}, id)))
	})
}

func FromContext(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

// MaxInFlight 用有 n 個空位的 channel 當號碼牌（19 章的 semaphore）；select 的 default 讓它不等待。
func MaxInFlight(n int, next http.Handler) http.Handler {
	slots := make(chan struct{}, n)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case slots <- struct{}{}:
			defer func() { <-slots }()
			next.ServeHTTP(w, r)
		default:
			w.Header().Set("Retry-After", "1")
			writeError(w, http.StatusServiceUnavailable, "太忙了，請稍後再試")
		}
	})
}

type ticketBody struct {
	Number int `json:"number"`
	Ahead  int `json:"ahead"`
}

func NewServer(q *Queue) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /tickets", func(w http.ResponseWriter, r *http.Request) {
		n, ahead := q.Take()
		w.Header().Set("Location", "/tickets/"+strconv.Itoa(n))
		writeJSON(w, http.StatusCreated, ticketBody{n, ahead})
	})
	mux.HandleFunc("GET /tickets/{number}", func(w http.ResponseWriter, r *http.Request) {
		n, err := strconv.Atoi(r.PathValue("number"))
		if err != nil || n < 1 {
			writeError(w, http.StatusBadRequest, "號碼要是正整數")
			return
		}
		ahead, ok := q.Position(n)
		if !ok {
			writeError(w, http.StatusNotFound, "這個號碼不在等候中")
			return
		}
		writeJSON(w, http.StatusOK, ticketBody{n, ahead})
	})
	mux.HandleFunc("POST /next", func(w http.ResponseWriter, r *http.Request) {
		n, ok := q.Call()
		if !ok {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		writeJSON(w, http.StatusOK, map[string]int{"number": n})
	})
	return mux
}

// Rename 回了錯誤之後要 return，否則下面的程式照樣執行。
func (s *Settings) Rename(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name string `json:"name"`
	}
	if err := json.UnmarshalRead(r.Body, &in); err != nil || in.Name == "" {
		http.Error(w, "需要 name", http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	s.name = in.Name
	s.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.MarshalWrite(w, map[string]string{"name": in.Name})
}

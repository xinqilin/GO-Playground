package main

import (
	"encoding/json/v2"
	"net/http"
	"strconv"
	"sync"
)

// MenuItem 是菜單上的一項。
type MenuItem struct {
	Dish  string `json:"dish"`
	Price int    `json:"price"`
}

// Order 是一筆訂單。
type Order struct {
	ID    int    `json:"id"`
	Dish  string `json:"dish"`
	Qty   int    `json:"qty"`
	Total int    `json:"total"`
}

// fakeShop 是示範用的店家 API，和「HTTP 服務」那一章的很像，但小得多。
func fakeShop() http.Handler {
	var mu sync.Mutex
	stock := map[string]int{"排骨飯": 3, "滷肉飯": 20}
	prices := map[string]int{"排骨飯": 90, "滷肉飯": 50}
	nextID := 1

	reply := func(w http.ResponseWriter, status int, v any) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.MarshalWrite(w, v)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /menu", func(w http.ResponseWriter, r *http.Request) {
		reply(w, http.StatusOK, []MenuItem{{"滷肉飯", 50}, {"排骨飯", 90}})
	})
	mux.HandleFunc("POST /orders", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Dish string `json:"dish"`
			Qty  int    `json:"qty"`
		}
		if err := json.UnmarshalRead(r.Body, &in); err != nil || in.Qty < 1 {
			reply(w, http.StatusBadRequest, map[string]string{"error": "訂單格式不對"})
			return
		}
		mu.Lock()
		defer mu.Unlock()
		if stock[in.Dish] < in.Qty {
			reply(w, http.StatusConflict, map[string]string{"error": in.Dish + " 賣完了"})
			return
		}
		stock[in.Dish] -= in.Qty
		o := Order{ID: nextID, Dish: in.Dish, Qty: in.Qty, Total: prices[in.Dish] * in.Qty}
		nextID++
		w.Header().Set("Location", "/orders/"+strconv.Itoa(o.ID))
		reply(w, http.StatusCreated, o)
	})
	mux.HandleFunc("GET /orders/{id}", func(w http.ResponseWriter, r *http.Request) {
		reply(w, http.StatusNotFound, map[string]string{"error": "找不到這筆訂單"})
	})
	return mux
}

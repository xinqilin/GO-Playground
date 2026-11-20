package main

import (
	"fmt"
	"net/http"
)

// whichRoute 印出這個請求對到哪個 pattern，以及路徑裡的變數。
func whichRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pattern=%q id=%q path=%q", r.Pattern, r.PathValue("id"), r.PathValue("path"))
}

// routeDemo 一個請求同時符合好幾個 pattern 時，最具體的那一個贏，和註冊的順序無關。
func routeDemo() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /orders/{id}", whichRoute)
	mux.HandleFunc("GET /orders/latest", whichRoute)
	mux.HandleFunc("POST /orders", whichRoute)
	mux.HandleFunc("GET /menu/", whichRoute)
	mux.HandleFunc("GET /files/{path...}", whichRoute)
	mux.HandleFunc("GET /{$}", whichRoute)

	for _, req := range []struct{ method, target string }{
		{"GET", "/orders/42"},
		{"GET", "/orders/latest"},
		{"HEAD", "/orders/42"},
		{"DELETE", "/orders/42"},
		{"GET", "/menu/today"},
		{"GET", "/menu"},
		{"GET", "/files/2026/11/menu.pdf"},
		{"GET", "/"},
		{"GET", "/about"},
	} {
		send(mux, req.method, req.target, "", "Allow", "Location")
	}
}

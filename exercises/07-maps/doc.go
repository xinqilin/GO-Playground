// Package mapping 是第 07 章「map」的練習題。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/07-maps/
//
// 一開始全部是紅燈，寫完一題就會綠一題。只想跑某一題時加上 -run：
//
//	go test ./exercises/07-maps/ -run TestTally
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/07-maps/
package mapping

// Item 是一項商品的價格和庫存。
type Item struct {
	Price int
	Qty   int
}

// Station 是一個車站和它所在的路線。轉乘站在每條路線上各出現一次。
type Station struct {
	Name string
	Line string
}

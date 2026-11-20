// Package kiosk 是第 21 章「命令列工具」的練習題。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/21-cli/
//
// 只想跑某一題時加上 -run：
//
//	go test ./exercises/21-cli/ -run TestLoadConfig
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/21-cli/
//
// 每一題都不碰 os.Args、os.Stdout、os.Exit：參數、輸入輸出和環境變數都從參數傳進來，
// 測試才能直接呼叫，就像本章的 run。
package kiosk

// Booking 是一筆訂位（練習 1）。
type Booking struct {
	Name    string // -name，一定要有
	People  int    // -people，預設 2，1 到 12 人
	Time    string // -time，預設 "18:00"，格式是 15:04
	Outdoor bool   // -outdoor，要不要坐戶外
	Note    string // 旗標後面剩下的參數，用一個空格接起來
}

// Config 是點餐機的設定（練習 2）。
type Config struct {
	Port  int    // -port，環境變數 KIOSK_PORT，預設 8080
	Shop  string // -shop，環境變數 KIOSK_SHOP，預設 "總店"
	Debug bool   // -debug，環境變數 KIOSK_DEBUG（strconv.ParseBool 認得的值），預設 false
}

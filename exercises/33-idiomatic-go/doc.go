// Package idiom 是第 33 章「寫出道地的 Go」的練習題：命名、doc comment、零值可用、
// functional options、在邊界上複製。有幾題是小小的檢查工具，code review 時用得到。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/33-idiomatic-go/
//
// 只想跑某一題時加上 -run：
//
//	go test ./exercises/33-idiomatic-go/ -run TestMixedCaps
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/33-idiomatic-go/
package idiom

// DocProblems 回傳的問題，依這個順序列出。
const (
	ProblemMissing = "缺少說明"
	ProblemName    = "沒有以名稱開頭"
	ProblemPeriod  = "沒有以句號結尾"
)

// Initialisms 是名稱裡要整個大寫（或整個小寫）的縮寫。
var Initialisms = map[string]bool{
	"api": true, "html": true, "http": true, "id": true, "ip": true,
	"json": true, "sql": true, "tcp": true, "url": true, "uuid": true,
}

// Tally 統計每道菜賣了幾份。零值就可以使用。
type Tally struct {
	counts map[string]int
}

// Option 是 New 的選項。
type Option func(*Machine)

// Machine 發出「A001」這種格式的號碼牌。用 New 建立。
type Machine struct {
	prefix string
	next   int
	width  int
}

// Queue 是排隊的號碼，依排隊的順序。零值是空的隊伍。
type Queue struct {
	waiting []int
}

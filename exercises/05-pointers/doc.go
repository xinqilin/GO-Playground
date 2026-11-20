// Package pointers 是第 05 章「指標」的練習題。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/05-pointers/
//
// 一開始全部是紅燈，寫完一題就會綠一題。只想跑某一題時加上 -run：
//
//	go test ./exercises/05-pointers/ -run TestSwap
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/05-pointers/
package pointers

// Profile 是會員資料。
type Profile struct {
	Name   string
	Phone  string
	Points int
}

// ProfilePatch 描述要更新哪些欄位：nil 的欄位不更新，非 nil 的欄位改成它指向的值。
// 用指標才分得出「不更新」和「改成零值」，例如把點數歸零。
type ProfilePatch struct {
	Name   *string
	Phone  *string
	Points *int
}

// Item 是一項商品。
type Item struct {
	Name  string
	Price int
	Stock int
}

// Node 是排隊名單裡的一個人。Next 指向排在他後面的人，最後一個人的 Next 是 nil。
type Node struct {
	Name string
	Next *Node
}

// Order 是一張訂單，Note 是可有可無的備註。
type Order struct {
	ID   int
	Note *string
}

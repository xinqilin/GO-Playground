//go:build !solution

package pointers

// Swap 交換 a 和 b 指向的兩個值。
func Swap(a, b *int) {
	// TODO: 實作
}

// Clamp 把 p 指向的值限制在 lo 到 hi 之間：小於 lo 就改成 lo，大於 hi 就改成 hi。
// p 是 nil 時什麼都不做。
func Clamp(p *int, lo, hi int) {
	// TODO: 實作
}

// ApplyPatch 依照 patch 更新 p：patch 裡非 nil 的欄位才更新。
func ApplyPatch(p *Profile, patch ProfilePatch) {
	// TODO: 實作
}

// Cheapest 回傳指向 items 裡最便宜那一項的指標，讓呼叫端可以直接修改那一項，例如改庫存。
// 價格相同時取前面那一項。items 是空的時回傳 nil。
//
// 找錯：透過回傳的指標改庫存，items 裡的資料卻沒有變。
func Cheapest(items []Item) *Item {
	var best *Item
	for _, it := range items {
		if best == nil || it.Price < best.Price {
			best = &it
		}
	}
	return best
}

// PushBack 把 name 排到隊伍最後面，回傳隊伍的第一個人。head 是 nil 代表隊伍是空的。
func PushBack(head *Node, name string) *Node {
	// TODO: 實作
	return head
}

// Names 從第一個人開始，依序回傳隊伍裡每個人的名字。隊伍是空的時回傳 nil。
func Names(head *Node) []string {
	// TODO: 實作
	return nil
}

// Remove 把隊伍裡第一個叫 name 的人移出隊伍，回傳新的第一個人。找不到時隊伍不變。
func Remove(head *Node, name string) *Node {
	// TODO: 實作
	return head
}

// Duplicate 複製一張訂單，換上新的編號。複製出來的訂單要和原本的完全獨立：
// 修改其中一張的備註，另一張不能跟著變。
//
// 找錯：改了複本的備註，原本那張也跟著變了。
func Duplicate(o Order, newID int) Order {
	dup := o
	dup.ID = newID
	return dup
}

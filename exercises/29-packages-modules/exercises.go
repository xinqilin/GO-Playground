//go:build !solution

package modver

import "strings"

// Parse 解析 "v1.2.3" 或 "v1.2.3-rc.1" 這種版本字串。規則：
//   - 開頭是 v，後面是三個用 . 分開的非負整數；除了 0 本身，數字不能以 0 開頭（v1.02.0 不合法）。
//   - 可以接著 - 和預發行標記，標記不能是空的，只能有英文字母、數字、. 和 -。
//
// 不合法時回傳錯誤。
func Parse(s string) (Version, error) {
	// TODO: 實作
	return Version{}, nil
}

// Compare 比較兩個版本，a 比較舊回傳 -1，一樣回傳 0，比較新回傳 1。
//   - 先依序比較 Major、Minor、Patch 的數字大小（v1.10.0 比 v1.9.0 新）。
//   - 數字都一樣時，預發行版比正式版舊（v1.0.0-rc.1 比 v1.0.0 舊）。
//   - 兩個都是預發行版時，把標記用 . 切開逐段比較：都是數字就比數字大小，
//     都不是數字就比字串，數字比非數字舊；前面都一樣時，段數少的比較舊。
func Compare(a, b Version) int {
	// TODO: 實作
	return 0
}

// Latest 回傳 go get -u 會把 current 升到哪一版：在 available 裡，主版本號和 current 一樣、
// 不是預發行版、而且比 current 新的版本中最新的那一個。找不到時回傳 current 和 false。
func Latest(current Version, available []Version) (Version, bool) {
	// TODO: 實作
	return Version{}, false
}

// ParseRequires 讀出 go.mod 內容裡所有的 require，依照出現的順序。
// require 可以寫成一行（require path version），也可以寫成 require ( ... ) 區塊；
// 行尾的 // indirect 表示間接相依。其他指令（module、go、replace 區塊等）略過。
// require 那一行不是剛好「路徑 版本」兩個欄位（不算註解）時回傳錯誤。
func ParseRequires(gomod string) ([]Require, error) {
	// TODO: 實作
	return nil, nil
}

// Select 用最小版本選擇（MVS）算出建置時每個 module 用哪一個版本。
// reqs["main"] 是主 module 的 require，reqs["path@version"] 是那個版本自己的 require，
// 每一筆都寫成 "path@version"；沒有出現在 reqs 裡的就是沒有 require。
// 從 main 出發，走過所有找得到的 path@version（包括最後沒被選上的版本），
// 每個 path 選走過的版本裡最新的一個。回傳 path → 版本字串，不包含 main。
// 有版本解析不了時回傳錯誤。
func Select(reqs map[string][]string) (map[string]string, error) {
	// TODO: 實作
	return nil, nil
}

// InternalAllowed 回傳 importer 這個 package 能不能 import imported。
// 規則：imported 的路徑裡有一段剛好叫 internal 時（有好幾段的話看最後一段），
// 只有在 internal 的上一層目錄底下（包括那一層本身）的 package 才能 import 它。
//
// 找錯：example.com/shopping 竟然可以 import example.com/shop/internal/price。
func InternalAllowed(importer, imported string) bool {
	before, _, ok := strings.Cut(imported, "/internal")
	if !ok {
		return true
	}
	return strings.HasPrefix(importer, before)
}

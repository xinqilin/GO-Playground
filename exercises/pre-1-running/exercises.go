//go:build !solution

package running

// Motto 回傳一句話：我會寫 Go 了
// 只要把 return 後面的 "" 換成那句話（記得保留雙引號）。
func Motto() string {
	// TODO: 把 "" 換成答案
	return ""
}

// Receipt 回傳三行的收據，行和行之間用 \n 分開：
//
//	排骨便當 85
//	雞腿便當 90
//	合計 175
//
// 最後一行後面沒有 \n。
func Receipt() string {
	// TODO: 把 "" 換成答案
	return ""
}

// Total 回傳買 4 個雞腿便當（一個 90 元）的總價。請寫成算式讓電腦算，不要自己先算好。
func Total() int {
	// TODO: 把 0 換成算式
	return 0
}

// Average 回傳 72、85、91 三次考試的平均分數（整數）。
//
// 找錯：程式可以執行，但算出來的平均是 187，比滿分還高。
func Average() int {
	return 72 + 85 + 91/3
}

// Steps 依照正確的順序回傳泡茶的四個步驟：
// 燒開水、放茶葉、倒熱水、等三分鐘。
//
// 找錯：步驟都在，但順序錯了。
func Steps() []string {
	return []string{
		"放茶葉",
		"倒熱水",
		"燒開水",
		"等三分鐘",
	}
}

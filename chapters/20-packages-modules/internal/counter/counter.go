// Package counter 放在 internal 底下：只有 chapters/20-packages-modules 這棵目錄樹裡的程式碼能 import 它。
package counter

// n 是套件層級的變數，套件載入時初始化成零值。
// 示範用；多個 goroutine 同時呼叫 Next 的話要加鎖（第 17 章）。
var n int

// Next 回傳下一個號碼。
func Next() int {
	n++
	return n
}

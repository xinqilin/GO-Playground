// Package steps 記錄初始化的順序，給 initOrderDemo 看。
package steps

// Log 依序記下每一個被記錄的步驟。
var Log []string

// Record 記下一個步驟，並原樣回傳 v，方便放在變數的初始值裡。
func Record(name string, v int) int {
	Log = append(Log, name)
	return v
}

// 被 import 的 package 一定先初始化完，才輪到 import 它的 package。
func init() {
	Record("steps 的 init", 0)
}

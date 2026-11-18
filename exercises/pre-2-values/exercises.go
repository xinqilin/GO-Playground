//go:build !solution

package values

// Elapsed 回傳同一天裡，從 startH 點 startM 分到 endH 點 endM 分經過了幾分鐘。
// 時間用 24 小時制，結束時間不會早於開始時間。例如 9 點 50 分到 11 點 5 分是 75 分鐘。
func Elapsed(startH, startM, endH, endM int) int {
	// TODO: 實作
	return 0
}

// WeeksAndDays 把天數換成幾週又幾天，例如 17 天是 2 週又 3 天。
func WeeksAndDays(days int) (weeks, rest int) {
	// TODO: 實作
	return 0, 0
}

// IsHalfFare 回傳這個年齡能不能買半票：未滿 12 歲，或 65 歲以上。
func IsHalfFare(age int) bool {
	// TODO: 實作
	return false
}

// Predict 先不要執行，用追蹤表算出下面這段程式結束時 a、b、c 的值，
// 把三個答案填進 return：
//
//	a := 17
//	b := a / 5
//	c := a % 5
//	a = b*5 + c
//	b = b + c
//	c = a - b*b
func Predict() (a, b, c int) {
	// TODO: 把 0, 0, 0 換成你追蹤出來的答案
	return 0, 0, 0
}

// Coins 回傳找 amount 元時，50 元、10 元、5 元、1 元硬幣各要幾枚，總枚數要最少。
// 例如 78 元是 1 枚 50 元、2 枚 10 元、1 枚 5 元、3 枚 1 元。amount 不會是負數。
func Coins(amount int) (c50, c10, c5, c1 int) {
	// TODO: 實作
	return 0, 0, 0, 0
}

// Rotate 讓三個座位上的人往前換一格：原本依序坐著 a、b、c，
// 換完之後第一個位置坐 b、第二個坐 c、第三個坐 a。
//
// 找錯：Rotate("小明", "小華", "小美") 應該是 小華、小美、小明，
// 結果卻是 小華、小美、小華，小明不見了。
func Rotate(a, b, c string) (string, string, string) {
	a = b
	b = c
	c = a
	return a, b, c
}

// Attendance 回傳出席率（百分比）：出席 part 次、總共 total 次，要保留小數。
// total 一定大於 0。例如 20 次出席 18 次是 90，8 次出席 3 次是 37.5。
//
// 找錯：除了全勤是 100 以外，算出來全部是 0。
func Attendance(part, total int) float64 {
	return float64(part / total * 100)
}

//go:build !solution

package values

// Split 把 total 元平分給 people 個人，回傳每人拿到幾元、剩下幾元。
// people 一定大於 0。例如 Split(100, 3) 回傳 33, 1。
func Split(total, people int) (each, left int) {
	// TODO: 實作
	return 0, 0
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

// AverageScore 回傳兩次分數的平均，要保留小數，例如 80 和 85 的平均是 82.5。
//
// 找錯：80 和 85 算出來是 82。
func AverageScore(a, b int) float64 {
	return float64((a + b) / 2)
}

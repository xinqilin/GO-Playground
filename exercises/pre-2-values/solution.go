//go:build solution

package values

// Elapsed 先把兩個時間都換成從 0 點開始算的第幾分鐘，再相減。
func Elapsed(startH, startM, endH, endM int) int {
	return (endH*60 + endM) - (startH*60 + startM)
}

func WeeksAndDays(days int) (weeks, rest int) {
	return days / 7, days % 7
}

func IsHalfFare(age int) bool {
	return age < 12 || age >= 65
}

// Predict 的追蹤表：a=17, b=3, c=2；a=3*5+2=17；b=3+2=5；c=17-5*5=-8。
func Predict() (a, b, c int) {
	return 17, 5, -8
}

// Coins 從面額最大的硬幣開始：先算用幾枚，剩下的錢再交給下一種硬幣。
func Coins(amount int) (c50, c10, c5, c1 int) {
	c50 = amount / 50
	amount = amount % 50
	c10 = amount / 10
	amount = amount % 10
	c5 = amount / 5
	c1 = amount % 5
	return c50, c10, c5, c1
}

// Rotate 一次把三個值放到新位置，右邊的值都先拿好，不會被覆蓋。
func Rotate(a, b, c string) (string, string, string) {
	a, b, c = b, c, a
	return a, b, c
}

// Attendance 要先轉成 float64 再除；兩個整數先除，18 / 20 就已經是 0 了。
func Attendance(part, total int) float64 {
	return float64(part) * 100 / float64(total)
}

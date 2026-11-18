//go:build solution

package decisions

func ShippingFee(total int, member bool) int {
	if member || total >= 499 {
		return 0
	}
	return 60
}

func SumOdd(n int) int {
	total := 0
	for i := 1; i <= n; i++ {
		if i%2 == 1 {
			total += i
		}
	}
	return total
}

// CountDigits 每除以 10 就少一位數，數一數要除幾次才會變成 0。
// 先算一位、再用 for 處理剩下的位數，0 就自然是 1 位數。
func CountDigits(n int) int {
	count := 1
	for n >= 10 {
		n = n / 10
		count++
	}
	return count
}

// ParkingFee 的「不足 30 分鐘也算 30 分鐘」是無條件進位：先加 29 再除以 30。
func ParkingFee(minutes int) int {
	if minutes <= 30 {
		return 0
	}
	units := (minutes - 30 + 29) / 30
	fee := units * 20
	if fee > 200 {
		fee = 200
	}
	return fee
}

// RightTriangle 第 row 列先印 n-row 個空白，再印 row 顆星。
func RightTriangle(n int) string {
	s := ""
	for row := 1; row <= n; row++ {
		if row > 1 {
			s += "\n"
		}
		for i := 1; i <= n-row; i++ {
			s += " "
		}
		for i := 1; i <= row; i++ {
			s += "*"
		}
	}
	return s
}

// SumRange 的條件要用 <=，才會把 to 本身也加進去。
func SumRange(from, to int) int {
	total := 0
	for i := from; i <= to; i++ {
		total += i
	}
	return total
}

// DaysToSave 的條件是「還沒達到目標就繼續存」，也就是 total < target；
// 用 <= 的話，剛好達到目標時還會多存一天。
func DaysToSave(target, daily int) int {
	days := 0
	total := 0
	for total < target {
		days++
		total += daily
	}
	return days
}

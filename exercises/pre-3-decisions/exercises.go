//go:build !solution

package decisions

// ShippingFee 回傳網購的運費：會員一律免運；不是會員時，消費滿 499 元免運，否則運費 60 元。
func ShippingFee(total int, member bool) int {
	// TODO: 實作
	return 0
}

// SumOdd 回傳 1 到 n（包含 n）之間所有奇數的和。n 小於 1 時回傳 0。
// 例如 SumOdd(7) 是 1+3+5+7 = 16。
func SumOdd(n int) int {
	// TODO: 實作
	return 0
}

// CountDigits 回傳 n 是幾位數。n 不會是負數，0 算 1 位數。
// 例如 CountDigits(2026) 是 4。
func CountDigits(n int) int {
	// TODO: 實作
	return 0
}

// ParkingFee 回傳停車 minutes 分鐘的停車費：
//   - 前 30 分鐘免費；
//   - 超過的部分每 30 分鐘 20 元，不足 30 分鐘也算 30 分鐘；
//   - 一次最多收 200 元。
//
// 例如停 31 分鐘是 20 元，停 61 分鐘是 40 元，停 10 小時是 200 元。
func ParkingFee(minutes int) int {
	// TODO: 實作
	return 0
}

// RightTriangle 回傳 n 列、靠右對齊的星號三角形，每列之間用 \n 分開，最後一列後面沒有 \n。
// 第 1 列有 1 顆星，最後一列有 n 顆星，星號左邊用空白補齊。n 小於 1 時回傳 ""。
// 例如 n 是 3 時：
//
//	  *
//	 **
//	***
func RightTriangle(n int) string {
	// TODO: 實作
	return ""
}

// SumRange 回傳從 from 加到 to（兩端都包含）的總和；from 大於 to 時回傳 0。
//
// 找錯：SumRange(1, 10) 應該是 55，結果是 45；SumRange(5, 5) 應該是 5，結果是 0。
func SumRange(from, to int) int {
	total := 0
	for i := from; i < to; i++ {
		total += i
	}
	return total
}

// DaysToSave 回傳每天存 daily 元時，第幾天存款第一次達到 target 元（含）以上。
// target 和 daily 都大於 0。例如目標 100 元、每天存 30 元，第 4 天有 120 元，回傳 4。
//
// 找錯：目標 100 元、每天存 50 元，第 2 天剛好存到 100 元，結果卻回傳 3。
func DaysToSave(target, daily int) int {
	days := 0
	total := 0
	for total <= target {
		days++
		total += daily
	}
	return days
}

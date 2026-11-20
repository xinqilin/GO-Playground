//go:build !solution

package variables

// SplitBill 把 total 元平分給 people 個人：each 是每個人付的整數元，rest 是分不完的零頭。
// 例如 1000 元分給 3 個人，每人 333 元，剩 1 元。people 小於 1 時回傳 0, 0。
func SplitBill(total, people int) (each, rest int) {
	// TODO: 實作
	return 0, 0
}

// Tables 回傳 guests 位客人、每桌最多坐 perTable 人時，至少需要幾張桌子。
// 例如 10 位客人、每桌 4 人，需要 3 張桌子。guests 或 perTable 小於 1 時回傳 0。
func Tables(guests, perTable int) int {
	// TODO: 實作
	return 0
}

// Average 回傳三個分數的平均，例如 Average(80, 85, 91) 是 85.333…。
//
// 找錯：這個版本算出來的結果永遠沒有小數。
func Average(a, b, c int) float64 {
	return float64((a + b + c) / 3)
}

// 容量單位：KB 是 1024 bytes，MB 是 1024 KB，GB 是 1024 MB。
const (
	// TODO: 用 iota 改寫這三個常數，不要直接寫出數字
	KB = 0
	MB = 0
	GB = 0
)

// FormatSize 把 byte 數寫成好讀的字串，KB、MB、GB 都取到小數第一位：
//
//	512        -> "512 B"
//	1536       -> "1.5 KB"
//	3145728    -> "3.0 MB"
//	1073741824 -> "1.0 GB"
//
// 小於 1 KB 用 B，小於 1 MB 用 KB，小於 1 GB 用 MB，其他都用 GB。負數當成 0。
func FormatSize(bytes int64) string {
	// TODO: 實作
	return ""
}

// Grant 回傳 p 加上 q 裡所有權限之後的結果。
func Grant(p, q Perm) Perm {
	// TODO: 實作
	return 0
}

// Revoke 回傳 p 拿掉 q 裡所有權限之後的結果。
func Revoke(p, q Perm) Perm {
	// TODO: 實作
	return 0
}

// Has 回傳 p 是否擁有 q 裡的「每一個」權限。q 是 0 時回傳 true。
func Has(p, q Perm) bool {
	// TODO: 實作
	return false
}

// AddStock 把庫存 stock 加上 n。庫存計數器是 uint8，最多只能記到 255：
// 相加超過 255 時回傳 255, false，不能讓它繞回小的數字；沒超過時回傳相加的結果和 true。
func AddStock(stock, n uint8) (uint8, bool) {
	// TODO: 實作
	return 0, false
}

// Ticket 回傳號碼牌 current 往後數 delta 號的號碼。號碼從 1 到 999 循環：
// 999 的下一號是 1，1 的前一號是 999。delta 可以是負數（往回數），也可以超過 999。
// current 一定在 1 到 999 之間。
func Ticket(current, delta int) int {
	// TODO: 實作
	return 0
}

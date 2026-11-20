// Package failing 故意放了一個有 bug 的函式，用來看 go test 失敗時長什麼樣子。
package failing

// MemberPrice 回傳會員價：打 9 折，四捨五入到整數。
func MemberPrice(price int) int {
	return price / 10 * 9 // bug：先除再乘，個位數直接不見了
}

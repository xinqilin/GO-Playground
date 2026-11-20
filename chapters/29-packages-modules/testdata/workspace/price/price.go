// Package price 是另一個 module：放在同一個 workspace 裡，改了馬上就能在 shop 看到效果。
package price

// Of 回傳一道菜的價格，不認得的菜回傳 0。
func Of(dish string) int {
	switch dish {
	case "排骨飯":
		return 90
	case "滷肉飯":
		return 50
	}
	return 0
}

// Package recipe 是廚房內部的食譜，只有 kitchen 底下的程式碼可以 import。
package recipe

// Steps 回傳一道菜的步驟數。
func Steps(dish string) int { return len(dish) }

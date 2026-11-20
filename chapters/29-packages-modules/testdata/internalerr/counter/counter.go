// Package counter 在 kitchen 外面，不能 import kitchen/internal 底下的東西。
package counter

import "example.com/shop/kitchen/internal/recipe"

// Peek 想偷看廚房的食譜。
func Peek(dish string) int { return recipe.Steps(dish) }

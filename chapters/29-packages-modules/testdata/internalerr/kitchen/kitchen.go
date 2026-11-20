// Package kitchen 可以 import 自己底下的 internal/recipe。
package kitchen

import "example.com/shop/kitchen/internal/recipe"

// Cook 回傳做這道菜要幾個步驟。
func Cook(dish string) int { return recipe.Steps(dish) }

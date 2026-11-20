// Package order 下單時通知廚房。
package order

import "example.com/cycle/kitchen"

// Place 下單。
func Place(dish string) { kitchen.Cook(dish) }

// Package kitchen 做好之後想回頭更新訂單：兩個 package 互相 import。
package kitchen

import "example.com/cycle/order"

// Cook 做菜，做完再下一張單（只是為了製造循環）。
func Cook(dish string) { order.Place(dish) }

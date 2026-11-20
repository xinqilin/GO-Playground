// Package receipt 印收據。它放在 internal 底下：
// 只有 chapters/29-packages-modules 這棵目錄樹裡的程式碼可以 import 它，
// 所以之後怎麼改它的 API，都不會弄壞別人的程式。
package receipt

import (
	"fmt"
	"strings"
)

// Line 是收據上的一行。
type Line struct {
	Dish  string
	Qty   int
	Price int
}

// Format 把每一行和總計排成收據。
func Format(lines []Line) string {
	var b strings.Builder
	total := 0
	for _, l := range lines {
		fmt.Fprintf(&b, "%s x%d %d\n", l.Dish, l.Qty, l.Qty*l.Price)
		total += l.Qty * l.Price
	}
	fmt.Fprintf(&b, "合計 %d\n", total)
	return b.String()
}

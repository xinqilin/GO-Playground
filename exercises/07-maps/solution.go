//go:build solution

package mapping

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

func Tally(orders []string) (counts map[string]int, top string) {
	counts = make(map[string]int)
	for _, o := range orders {
		counts[o]++
	}
	for name, n := range counts {
		if top == "" || n > counts[top] || (n == counts[top] && name < top) {
			top = name
		}
	}
	return counts, top
}

func Invert(codes map[string]string) (map[string]string, bool) {
	names := make(map[string]string, len(codes))
	for code, name := range codes {
		if _, dup := names[name]; dup {
			return nil, false
		}
		names[name] = code
	}
	return names, true
}

func GroupByLine(stations []Station) map[string][]string {
	lines := make(map[string][]string)
	for _, s := range stations {
		lines[s.Line] = append(lines[s.Line], s.Name)
	}
	for _, names := range lines {
		slices.Sort(names) // 原地排序 map 裡的 slice，不需要再存回去
	}
	return lines
}

// Restock 從 map 拿出來的 struct 是複本，改完要存回去。
func Restock(stock map[string]Item, name string, n int) int {
	item := stock[name]
	item.Qty += n
	stock[name] = item
	return item.Qty
}

// Summary 先把 key 排序，再依序組字串。
func Summary(menu map[string]int) string {
	var parts []string
	for _, name := range slices.Sorted(maps.Keys(menu)) {
		parts = append(parts, fmt.Sprintf("%s=%d", name, menu[name]))
	}
	return strings.Join(parts, ", ")
}

func Common(a, b []string) []string {
	inA := make(map[string]bool, len(a))
	for _, s := range a {
		inA[s] = true
	}
	both := make(map[string]bool)
	for _, s := range b {
		if inA[s] {
			both[s] = true
		}
	}
	if len(both) == 0 {
		return nil
	}
	return slices.Sorted(maps.Keys(both))
}

// Hops 用廣度優先搜尋：先走完所有距離 1 的站，再走距離 2 的站，第一次碰到 to 時就是最短距離。
// dist 同時記錄「已經排進佇列」，避免繞圈。
func Hops(links map[string][]string, from, to string) int {
	dist := map[string]int{from: 0}
	queue := []string{from}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == to {
			return dist[cur]
		}
		for _, next := range links[cur] {
			if _, seen := dist[next]; !seen {
				dist[next] = dist[cur] + 1
				queue = append(queue, next)
			}
		}
	}
	return -1
}

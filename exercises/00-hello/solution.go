//go:build solution

package hello

import (
	"fmt"
	"strconv"
	"strings"
)

func Greeting(name string) string {
	if name == "" {
		name = "陌生人"
	}
	return "哈囉，" + name + "！"
}

func BentoTotal(qty int) int {
	if qty < 1 {
		return 0
	}
	price := 85
	if qty >= 10 {
		price = 80
	}
	return qty * price
}

func Clock(seconds int) string {
	if seconds < 0 {
		seconds = 0
	}
	h := seconds / 3600
	m := seconds % 3600 / 60
	s := seconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func Box(text string) string {
	border := "+" + strings.Repeat("-", len(text)+2) + "+"
	return border + "\n| " + text + " |\n" + border
}

func GoVersion(line string) (major, minor int, ok bool) {
	version, found := strings.CutPrefix(line, "go ")
	if !found {
		return 0, 0, false
	}
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return 0, 0, false
	}
	major, err1 := strconv.Atoi(parts[0])
	minor, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || major < 0 || minor < 0 {
		return 0, 0, false
	}
	return major, minor, true
}

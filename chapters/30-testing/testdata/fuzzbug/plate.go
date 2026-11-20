// Package fuzzbug 是一個有 bug 的車牌整理函式，用來看 fuzz 怎麼找到人沒想到的輸入。
package fuzzbug

import (
	"errors"
	"regexp"
	"strings"
)

var (
	errBadPlate = errors.New("車牌格式不對")
	plateRE     = regexp.MustCompile(`^[A-Z]{2,3}-[0-9]{4}$`)
)

// NormalizePlate 把車牌整理成 ABC-1234 的寫法。
func NormalizePlate(s string) (string, error) {
	s = strings.ToUpper(strings.ReplaceAll(s, " ", ""))
	if !strings.Contains(s, "-") {
		// 使用者沒打 -：在最後 4 碼前面補一個。
		s = s[:len(s)-4] + "-" + s[len(s)-4:]
	}
	if !plateRE.MatchString(s) {
		return "", errBadPlate
	}
	return s, nil
}

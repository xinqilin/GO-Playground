package parking

import "errors"

// ErrBadPlate 表示車牌的格式不對。
var ErrBadPlate = errors.New("車牌格式不對")

// NormalizePlate 把使用者輸入的車牌整理成標準寫法：
// 去掉空白和 -、英文字母轉大寫，剩下的必須是「2 或 3 個英文字母加 4 個數字」
// 或「4 個數字加 2 個英文字母」，輸出時在字母和數字之間加上 -。
// 例如 " abc 1234" 變成 "ABC-1234"，"1234ab" 變成 "1234-AB"。
func NormalizePlate(s string) (string, error) {
	var b, shape []byte
	for _, r := range s {
		switch {
		case r == ' ' || r == '-':
			continue
		case 'a' <= r && r <= 'z':
			r -= 'a' - 'A'
			fallthrough
		case 'A' <= r && r <= 'Z':
			shape = append(shape, 'A')
		case '0' <= r && r <= '9':
			shape = append(shape, '9')
		default:
			return "", ErrBadPlate
		}
		b = append(b, byte(r))
	}
	var cut int
	switch string(shape) {
	case "AA9999", "AAA9999":
		cut = len(b) - 4
	case "9999AA":
		cut = 4
	default:
		return "", ErrBadPlate
	}
	return string(b[:cut]) + "-" + string(b[cut:]), nil
}

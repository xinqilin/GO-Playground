// Package display 決定號碼在看板上長什麼樣子。
package display

import "fmt"

// Number 把號碼排成看板上的格式，例如 A007。
func Number(prefix string, n int) string {
	return fmt.Sprintf("%s%03d", prefix, n)
}

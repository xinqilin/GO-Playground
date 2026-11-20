package main

import (
	"fmt"
	"strings"
)

// joinLabels 收任何「元素有 String 方法」的 slice。[]NTD 不能當成 []fmt.Stringer 傳，泛型可以。
func joinLabels[T fmt.Stringer](s []T) string {
	parts := make([]string, len(s))
	for i, v := range s {
		parts[i] = v.String()
	}
	return strings.Join(parts, "、")
}

func whenDemo() {
	prices := []NTD{85, 95}
	fmt.Println(joinLabels(prices))
}

// kindOf 想依 T 的實際型別分流時，先轉成 any 再做 type switch。
func kindOf[T any](v T) string {
	switch x := any(v).(type) {
	case int:
		return fmt.Sprint("整數 ", x)
	case string:
		return "字串 " + x
	default:
		return fmt.Sprintf("其他 %T", x)
	}
}

func kindDemo() {
	fmt.Println(kindOf(85), kindOf("排骨飯"), kindOf(NTD(30)))
}

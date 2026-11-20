package main

import (
	"fmt"
	"strconv"
)

// parsePeople 把使用者輸入的人數轉成 int。失敗時回傳 0 和一個說明原因的 error。
func parsePeople(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if n < 1 || n > 12 {
		return 0, fmt.Errorf("人數 %d 不在 1 到 12 之間", n)
	}
	return n, nil
}

// basicsDemo：每次呼叫完都先檢查 err，有錯就處理並跳過，沒錯才使用結果。
func basicsDemo() {
	for _, s := range []string{"4", "四", "20"} {
		n, err := parsePeople(s)
		if err != nil {
			fmt.Printf("錯誤：%v\n", err)
			continue
		}
		fmt.Println("人數", n)
	}
}

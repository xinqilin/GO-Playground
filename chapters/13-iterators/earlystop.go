package main

import (
	"fmt"
	"iter"
)

// receipts 模擬一個需要收尾的資料來源：開始時打開收據匣，結束時一定要關上。
func receipts() iter.Seq[string] {
	return func(yield func(string) bool) {
		fmt.Println("打開收據匣")
		defer fmt.Println("關上收據匣")
		for _, r := range []string{"A-001", "A-002", "A-003"} {
			if !yield(r) {
				return
			}
		}
	}
}

// earlyStopDemo：break 讓 yield 回傳 false，疊代器 return，它的 defer 在迴圈結束前執行。
func earlyStopDemo() {
	for r := range receipts() {
		fmt.Println("處理", r)
		if r == "A-002" {
			break
		}
	}
	fmt.Println("迴圈結束")
}

// rudeNumbers 不理會 yield 的回傳值：迴圈已經 break 了，它還繼續交出下一個。
func rudeNumbers() iter.Seq[int] {
	return func(yield func(int) bool) {
		for n := 1; n <= 3; n++ {
			yield(n)
		}
	}
}

func ignoreYieldDemo() {
	defer func() {
		fmt.Println("panic:", recover())
	}()
	for n := range rudeNumbers() {
		fmt.Println(n)
		break
	}
}

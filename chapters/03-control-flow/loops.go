package main

import "fmt"

// forClassic 是三段式寫法：初始化；條件；每一圈結束後要做的事。
func forClassic() {
	for n := 1; n <= 3; n++ {
		fmt.Printf("請 %d 號到櫃台\n", n)
	}
}

// forCondition 只寫條件，就是其他語言的 while。
func forCondition() {
	balance, months := 0, 0
	for balance < 10000 {
		balance += 3000
		months++
	}
	fmt.Println(months, "個月後存到", balance, "元")
}

// forInfinite 什麼條件都不寫，就會一直跑，靠 break 或 return 離開。
func forInfinite() {
	queue := []string{"A01", "A02", "A03"}
	for {
		if len(queue) == 0 {
			fmt.Println("沒有人在排隊了")
			break
		}
		next := queue[0]
		queue = queue[1:]
		fmt.Println("叫號", next)
	}
}

// forRangeInt 從 Go 1.22 起可以直接 range 一個整數，i 會是 0、1、2。
func forRangeInt() {
	for i := range 3 {
		fmt.Println("第", i, "圈")
	}
}

// forContinue 用 continue 跳過這一圈剩下的部分，直接進入下一圈。
func forContinue() {
	for n := range 10 {
		if n%3 != 0 {
			continue
		}
		fmt.Println(n, "是 3 的倍數")
	}
}

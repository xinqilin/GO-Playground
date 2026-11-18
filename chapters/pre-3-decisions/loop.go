package main

import "fmt"

// countdown 的三個要素：起點 i := 3，條件 i >= 1，每圈結束時 i--。
func countdown() {
	for i := 3; i >= 1; i-- {
		fmt.Println(i)
	}
	fmt.Println("發射")
}

// laps 跑 5 圈操場。
func laps() {
	for lap := 1; lap <= 5; lap++ {
		fmt.Println("第", lap, "圈")
	}
}

// sumTo100 用累加器 total 把 1 到 100 加起來。
func sumTo100() {
	total := 0
	for i := 1; i <= 100; i++ {
		total += i
	}
	fmt.Println(total)
}

// countMultiples 用計數器 count 數出 1 到 30 裡有幾個 3 的倍數。
func countMultiples() {
	count := 0
	for i := 1; i <= 30; i++ {
		if i%3 == 0 {
			count++
		}
	}
	fmt.Println(count)
}

// savings 第 1 天存 10 元，之後每天比前一天多存 10 元，10 天一共存多少。
func savings() {
	daily := 10
	total := 0
	for day := 1; day <= 10; day++ {
		total += daily
		daily += 10
	}
	fmt.Println(total)
}

// bounce 不知道要跑幾圈，只知道什麼時候該停：球彈起的高度不足 5 公分就停。
func bounce() {
	height := 100
	times := 0
	for height >= 5 {
		height = height / 2
		times++
		fmt.Println("第", times, "次彈起", height, "公分")
	}
}

// firstMultiple 找出大於 50 的第一個 7 的倍數，找到就用 break 離開迴圈。
func firstMultiple() {
	for n := 51; n <= 100; n++ {
		if n%7 == 0 {
			fmt.Println(n)
			break
		}
	}
}

// openDays 一週 7 天，第 3 天公休：continue 跳過這一圈剩下的部分，直接進入下一圈。
func openDays() {
	for day := 1; day <= 7; day++ {
		if day == 3 {
			continue
		}
		fmt.Println("第", day, "天營業")
	}
}

// triangle 外層迴圈決定第幾列，內層迴圈決定這一列印幾顆星。
func triangle() {
	for row := 1; row <= 4; row++ {
		for col := 1; col <= row; col++ {
			fmt.Print("*")
		}
		fmt.Println()
	}
}

// seats 印出 3 排、每排 4 個座位的編號。
func seats() {
	for row := 1; row <= 3; row++ {
		for seat := 1; seat <= 4; seat++ {
			if seat > 1 {
				fmt.Print(" ")
			}
			fmt.Print(row, "-", seat)
		}
		fmt.Println()
	}
}

// neverTen 的條件是 n != 10，但 n 是 1、3、5…，永遠不會剛好等於 10。
// 為了示範，跑到第 8 圈就用 break 強制離開；少了這個保護，它會一直跑下去。
func neverTen() {
	rounds := 0
	for n := 1; n != 10; n += 2 {
		fmt.Print(n, " ")
		rounds++
		if rounds == 8 {
			fmt.Println("… 已經跑了 8 圈，n 永遠不會剛好等於 10")
			break
		}
	}
}

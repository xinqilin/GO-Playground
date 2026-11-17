package main

import (
	"bufio"
	"fmt"
	"strings"
)

// scanLines 示範用 bufio.Scanner 一行一行讀。strings.NewReader 把字串當成輸入來源。
func scanLines() {
	input := "排骨飯 85\n雞腿飯 90\n\n滷蛋 15\n"
	sc := bufio.NewScanner(strings.NewReader(input))
	n := 0
	for sc.Scan() {
		n++
		fmt.Printf("第 %d 行：%q\n", n, sc.Text())
	}
	if err := sc.Err(); err != nil {
		fmt.Println("讀取失敗：", err)
	}
}

// scanWords 示範換掉切割方式，改成一個字一個字讀：連續的空白和換行都會被跳過。
func scanWords() {
	sc := bufio.NewScanner(strings.NewReader("  排骨飯   2\n滷蛋 3 "))
	sc.Split(bufio.ScanWords)
	var words []string
	for sc.Scan() {
		words = append(words, sc.Text())
	}
	if err := sc.Err(); err != nil {
		fmt.Println("讀取失敗：", err)
	}
	fmt.Println(len(words), words)
}

// scanTooLong 示範一行超過 Scanner 的上限（預設 64 KB）時，Scan 回傳 false，原因要從 Err 拿。
func scanTooLong() {
	long := strings.Repeat("x", 70_000)
	sc := bufio.NewScanner(strings.NewReader("短的一行\n" + long + "\n最後一行\n"))
	for sc.Scan() {
		fmt.Println("讀到", len(sc.Text()), "bytes")
	}
	fmt.Println("Err:", sc.Err())
}

// scanBigBuffer 示範用 Buffer 放寬一行的上限。
func scanBigBuffer() {
	long := strings.Repeat("x", 70_000)
	sc := bufio.NewScanner(strings.NewReader(long + "\n最後一行\n"))
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024) // 一行最多 1 MB
	for sc.Scan() {
		fmt.Println("讀到", len(sc.Text()), "bytes")
	}
	fmt.Println("Err:", sc.Err())
}

// sscanDemo 示範 fmt.Sscan 用空白切開字串，依序填進變數。& 讓函式能把值寫進變數，第 05 章介紹。
func sscanDemo() {
	var name string
	var qty, price int
	n, err := fmt.Sscan("排骨飯 2 85", &name, &qty, &price)
	fmt.Println(n, err, name, qty*price)

	n, err = fmt.Sscan("滷蛋 三 15", &name, &qty, &price)
	fmt.Println(n, err)
}

package main

import "fmt"

// arrayBasics 示範陣列的宣告：長度寫在型別裡，沒有給值的格子是零值。
func arrayBasics() {
	var sales [7]int // 一週七天的營業額，一開始全部是 0
	sales[0] = 3200
	sales[6] = 4100
	fmt.Println(sales, len(sales))

	weekdays := [...]string{"一", "二", "三", "四", "五"} // ... 讓編譯器自己數長度
	fmt.Println(len(weekdays), weekdays[len(weekdays)-1])

	booked := [5]bool{1: true, 3: true} // 指定索引給值，其他格是 false
	fmt.Println(booked)
}

// doubleAll 收到的是整個陣列的複本。
func doubleAll(prices [3]int) {
	for i := range prices {
		prices[i] *= 2
	}
}

// arrayIsValue 示範陣列是值：指定和傳參數都會複製整個陣列；同型別的陣列可以用 == 比較。
func arrayIsValue() {
	a := [3]int{85, 90, 15}
	b := a // 複製三格
	b[0] = 100
	fmt.Println(a, b)

	doubleAll(a)
	fmt.Println(a)

	fmt.Println(a == [3]int{85, 90, 15}, a == b)
}

// indexOutOfRange 示範執行時的索引超出範圍會 panic。這裡用 recover 接住才能印出訊息。
func indexOutOfRange() {
	defer func() {
		fmt.Println("panic:", recover())
	}()
	prices := [3]int{85, 90, 15}
	i := 3
	fmt.Println(prices[i])
}

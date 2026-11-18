package main

import "fmt"

// listBasics 建立一份價格清單，用索引讀取、修改，再用 append 加在最後面。
func listBasics() {
	prices := []int{85, 90, 75}
	fmt.Println(prices)
	fmt.Println(prices[0], prices[2])
	fmt.Println(len(prices))
	prices[1] = 95
	fmt.Println(prices)
	prices = append(prices, 100)
	fmt.Println(prices, len(prices))
}

// stringList 清單裡也可以放文字。最後一個元素的索引是 len-1。
func stringList() {
	menu := []string{"排骨飯", "雞腿飯", "素食餐"}
	fmt.Println(menu[0])
	fmt.Println(menu[len(menu)-1])
}

// emptyList 從空的清單開始，一筆一筆加進去。
func emptyList() {
	orders := []int{}
	fmt.Println(len(orders))
	orders = append(orders, 120)
	orders = append(orders, 85)
	fmt.Println(orders, len(orders))
}

// loopIndex 用索引走過清單：i 從 0 走到 len-1。
func loopIndex() {
	menu := []string{"排骨飯", "雞腿飯", "素食餐"}
	for i := 0; i < len(menu); i++ {
		fmt.Println(i, menu[i])
	}
}

// loopRange 用 range：每一圈同時拿到索引和值。
func loopRange() {
	prices := []int{85, 90, 75}
	for i, p := range prices {
		fmt.Println("第", i+1, "個：", p, "元")
	}
}

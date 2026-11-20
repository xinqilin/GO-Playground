package main

import "fmt"

// scopeLevels 示範作用域：內層看得到外層的名字，外層看不到內層的。
func scopeLevels() {
	fmt.Println("套件層級：", shopName)
	qty := 3 // 函式層級
	{
		price := 85 // 區塊層級：只在這對大括號裡有效
		fmt.Println("區塊裡：", qty*price)
	}
	// 這裡已經看不到 price
	fmt.Println("函式裡：", qty)
}

// shadowPackage 示範函式裡的同名變數會遮蔽套件層級的變數。
func shadowPackage() {
	shopName := "阿姨便當 二號店"
	fmt.Println(shopName)
	printShopName()
}

// printShopName 看到的永遠是套件層級的 shopName。
func printShopName() {
	fmt.Println(shopName)
}

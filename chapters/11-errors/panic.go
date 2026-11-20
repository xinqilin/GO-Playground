package main

import "fmt"

// mustPrice 給寫死在程式裡的品項用：查不到代表程式本身寫錯了，所以直接 panic。
func mustPrice(name string) int {
	p, err := price(name)
	if err != nil {
		panic(fmt.Sprintf("mustPrice(%q)：%v", name, err))
	}
	return p
}

// comboPrice 在程式啟動時就算好，品項名稱打錯會立刻當掉，不會等到客人點餐才發現。
var comboPrice = mustPrice("排骨飯") + 25

func mustDemo() {
	fmt.Println(comboPrice)

	defer func() {
		fmt.Println("panic:", recover())
	}()
	fmt.Println(mustPrice("珍奶"))
}

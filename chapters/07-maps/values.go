package main

import "fmt"

type item struct {
	price int
	stock int
}

// structValue 示範 map 的值是 struct 時，不能直接改欄位：要拿出來、改完、再放回去。
func structValue() {
	menu := map[string]item{"排骨": {price: 85, stock: 10}}
	it := menu["排骨"] // 拿到的是複本
	it.stock--
	fmt.Println(menu["排骨"].stock, it.stock)
	menu["排骨"] = it // 放回去
	fmt.Println(menu["排骨"].stock)
}

// pointerValue 示範值改成指標，就能透過指標直接修改；但 key 不存在時拿到的是 nil 指標。
func pointerValue() {
	menu := map[string]*item{"排骨": {price: 85, stock: 10}}
	menu["排骨"].stock--
	fmt.Println(menu["排骨"].stock)

	if it, ok := menu["魚排"]; ok {
		it.stock--
	} else {
		fmt.Println("沒有魚排", it == nil)
	}
}

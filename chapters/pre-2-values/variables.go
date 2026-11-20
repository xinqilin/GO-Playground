package main

import "fmt"

// variables 把值放進有名字的盒子，之後用名字拿出來用。
func variables() {
	price := 85
	qty := 3
	fmt.Println(price * qty)

	var name string = "小明"
	fmt.Println(name, "買了", qty, "個便當")
}

// reassign 讓同一個變數裝進新的值。右邊先算好，再放進左邊的盒子。
func reassign() {
	count := 0
	fmt.Println(count)
	count = count + 1
	fmt.Println(count)
	count = count + 1
	fmt.Println(count)
}

// shortcuts 是「拿自己來算，再放回自己」的簡寫。
func shortcuts() {
	points := 10
	points += 5
	fmt.Println(points)
	points -= 3
	fmt.Println(points)
	points *= 2
	fmt.Println(points)
	points++
	fmt.Println(points)
}

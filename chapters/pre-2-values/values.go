package main

import "fmt"

// values 印出四種最常見的值：整數、小數、文字、真假。
func values() {
	fmt.Println(85)
	fmt.Println(3.5)
	fmt.Println("排骨便當")
	fmt.Println(true)
}

// valueTypes 用 %T 印出每個值的型別。
func valueTypes() {
	fmt.Printf("%v 是 %T\n", 85, 85)
	fmt.Printf("%v 是 %T\n", 3.5, 3.5)
	fmt.Printf("%v 是 %T\n", "排骨便當", "排骨便當")
	fmt.Printf("%v 是 %T\n", true, true)
}

// textVsNumber 說明同樣是 +，數字是相加，文字是接起來。
func textVsNumber() {
	fmt.Println(85 + 3)
	fmt.Println("85" + "3")
	fmt.Println("排骨" + "便當")
}

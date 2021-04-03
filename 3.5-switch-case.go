package main

import "fmt"

func main() {
	ans := "a"
	switch ans {
	case "a":
		fmt.Println("ans is a")
		fallthrough
	case "b":
		fmt.Println("ans is b")
	case "c":
		fmt.Println("ans is c")
	default:
		fmt.Println("no answer")
	}
}

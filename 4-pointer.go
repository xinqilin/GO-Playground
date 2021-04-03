package main

import "fmt"

func main() {
	var x int = 5
	fmt.Println(x)

	fmt.Println(&x)

	var pointer *int = &x
	fmt.Println(pointer)

	fmt.Println(*pointer)
}

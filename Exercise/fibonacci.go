package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(20)
	fmt.Println("random number: ", n)
	ans := fibonacci(n)
	fmt.Println("ans: ", ans)
}

func fibonacci(n int) int {
	if n == 1 || n == 2 {
		return n
	} else {
		return fibonacci(n-1) + fibonacci(n-2)
	}

}

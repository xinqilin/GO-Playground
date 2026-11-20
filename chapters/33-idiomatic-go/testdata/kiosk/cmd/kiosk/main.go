// kiosk 是門口的取號機。
package main

import (
	"fmt"

	"example.com/kiosk/internal/display"
	"example.com/kiosk/internal/queue"
)

func main() {
	var q queue.Queue
	for range 3 {
		fmt.Printf("您的號碼：%s\n", display.Number("A", q.Take()))
	}
}

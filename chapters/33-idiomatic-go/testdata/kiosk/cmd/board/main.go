// board 是櫃台上方的叫號看板。
package main

import (
	"fmt"

	"example.com/kiosk/internal/display"
	"example.com/kiosk/internal/queue"
)

func main() {
	var q queue.Queue
	q.Take()
	q.Take()
	for {
		n, ok := q.Call()
		if !ok {
			fmt.Println("目前沒有人在等")
			return
		}
		fmt.Println("請", display.Number("A", n), "到櫃台")
	}
}

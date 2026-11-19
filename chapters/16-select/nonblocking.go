package main

import "fmt"

// tryReceive 用 default 做不等待的收：有就收，沒有就馬上離開。
func tryReceive(ch <-chan int) string {
	select {
	case v := <-ch:
		return fmt.Sprint("收到 ", v)
	default:
		return "現在沒有"
	}
}

func tryReceiveDemo() {
	ch := make(chan int, 1)
	fmt.Println(tryReceive(ch))
	ch <- 17
	fmt.Println(tryReceive(ch))
}

// dropDemo：用 default 做不等待的送。紀錄的 channel 滿了就丟掉，點餐流程不會因為寫紀錄而卡住。
func dropDemo() {
	logs := make(chan string, 2)
	dropped := 0
	for _, o := range []string{"排骨飯", "雞腿飯", "魚排飯", "滷肉飯", "控肉飯"} {
		select {
		case logs <- o:
		default:
			dropped++
		}
	}
	fmt.Println(len(logs), dropped)
}

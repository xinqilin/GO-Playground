package main

import "fmt"

// waitFlag 是錯誤的寫法：用一個普通的 bool 當「準備好了」的訊號。
// 兩個 goroutine 之間沒有同步，編譯器和 CPU 都可以重新安排讀寫：
// 迴圈可能永遠看不到 ready 變成 true，也可能看到 true，卻讀到還沒寫好的 menu。
func waitFlag() []string {
	var menu []string
	ready := false
	go func() {
		menu = []string{"排骨飯", "雞腿飯"}
		ready = true
	}()
	for !ready {
	}
	return menu
}

// waitDone 用 channel 當訊號：close(done) 保證發生在 <-done 收到之前，
// 所以 close 之前寫好的 menu，收到之後一定看得到。
func waitDone() []string {
	var menu []string
	done := make(chan struct{})
	go func() {
		menu = []string{"排骨飯", "雞腿飯"}
		close(done)
	}()
	<-done
	return menu
}

func orderDemo() {
	fmt.Println(waitDone())
}

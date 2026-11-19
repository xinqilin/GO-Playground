package main

import (
	"errors"
	"fmt"
	"sync"
)

// onceDemo：10 個 goroutine 同時要菜單，讀取菜單的函式只會執行一次，
// 其他人等它執行完，拿到同一份結果。
func onceDemo() {
	loads := 0
	menu := sync.OnceValue(func() []string {
		loads++
		fmt.Println("讀取菜單")
		return []string{"排骨飯", "雞腿飯"}
	})

	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() { menu() })
	}
	wg.Wait()
	fmt.Println(menu(), "讀取次數", loads)
}

// onceErrDemo：OnceValues 連錯誤也只算一次，之後每次呼叫都拿到同一個錯誤，不會重試。
func onceErrDemo() {
	tries := 0
	load := sync.OnceValues(func() (int, error) {
		tries++
		return 0, errors.New("價目表讀取失敗")
	})

	for range 3 {
		_, err := load()
		fmt.Println(err)
	}
	fmt.Println("實際嘗試", tries, "次")
}

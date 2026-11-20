package bento_test

import (
	"fmt"

	"github.com/xinqilin/GO-Playground/chapters/29-packages-modules/bento"
)

// ExampleMenu 會出現在 go doc 和 pkg.go.dev 的文件裡，go test 也會檢查它的輸出。
// 測試檔用 bento_test 這個 package，就只能用匯出的名稱，和真正的使用者一樣。
func ExampleMenu() {
	m := bento.New()
	m.Add("排骨飯", 90)
	price, ok := m.Price("排骨飯")
	fmt.Println(price, ok)
	// Output: 90 true
}

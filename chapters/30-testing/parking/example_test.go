package parking_test

import (
	"fmt"
	"time"

	"github.com/xinqilin/GO-Playground/chapters/30-testing/parking"
)

// 外部測試 package（parking_test）只能用匯出的名稱，和真正的使用者一樣。
// Example 會顯示在 pkg.go.dev 的文件裡（go doc -ex 也會列出來），go test 會比對它的輸出。

func ExampleFee() {
	tw := time.FixedZone("UTC+8", 8*60*60)
	in := time.Date(2026, 11, 20, 23, 50, 0, 0, tw)
	fee, _ := parking.Fee(in, in.Add(30*time.Minute))
	fmt.Println(fee, "元：跨過午夜，兩天各算一個單位")
	// Output: 40 元：跨過午夜，兩天各算一個單位
}

func ExampleNormalizePlate() {
	for _, s := range []string{" abc 1234", "1234ab", "ab-12"} {
		p, err := parking.NormalizePlate(s)
		fmt.Printf("%q %v\n", p, err)
	}
	// Output:
	// "ABC-1234" <nil>
	// "1234-AB" <nil>
	// "" 車牌格式不對
}

// 走訪 map 的順序每次不同，用 Unordered output 只比對有哪些行。
func ExampleNormalizePlate_unordered() {
	inputs := map[string]string{"三碼": "abc1234", "兩碼": "AB 1234", "數字在前": "1234-ab"}
	for kind, s := range inputs {
		p, _ := parking.NormalizePlate(s)
		fmt.Println(kind, p)
	}
	// Unordered output:
	// 數字在前 1234-AB
	// 三碼 ABC-1234
	// 兩碼 AB-1234
}

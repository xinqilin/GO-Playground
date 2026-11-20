package greet

import (
	"fmt"
	"testing"
)

// 測試檔和被測的程式碼同一個 package，所以也能測沒有匯出的 normalize。
func TestNormalize(t *testing.T) {
	tests := []struct{ in, want string }{
		{"Gopher", "Gopher"},
		{"  Go  ", "Go"},
		{"", "陌生人"},
	}
	for _, tt := range tests {
		if got := normalize(tt.in); got != tt.want {
			t.Errorf("normalize(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

// ExampleHello 會出現在 go doc 的文件裡，go test 也會檢查它的輸出。
func ExampleHello() {
	fmt.Println(Hello("Gopher"))
	// Output: 哈囉，Gopher
}

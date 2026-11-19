package textutil

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

// table-driven test：把案例寫成表格，一個迴圈跑完；t.Run 讓每個案例有自己的名字。
func TestReverse(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"空字串", "", ""},
		{"英文", "Go", "oG"},
		{"中文", "你好世界", "界世好你"},
		{"混合", "Go語言", "言語oG"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reverse(tt.in); got != tt.want {
				t.Errorf("Reverse(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// assertEqual 是測試輔助函式。t.Helper() 讓失敗訊息指向呼叫它的那一行，而不是這裡。
func assertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestIsPalindrome(t *testing.T) {
	assertEqual(t, IsPalindrome("Never odd or even"), true)
	assertEqual(t, IsPalindrome("上海自來水來自海上"), true)
	assertEqual(t, IsPalindrome("Go"), false)
}

// Example 開頭的函式是可執行的文件：go test 會比對 Output，go doc 會把它顯示出來。
func ExampleReverse() {
	fmt.Println(Reverse("Go語言"))
	// Output: 言語oG
}

// Benchmark：go test -bench=. 才會執行。Go 1.24 起用 b.Loop()，不用再自己寫 b.N 迴圈。
func BenchmarkReverse(b *testing.B) {
	s := strings.Repeat("Go語言", 100)
	for b.Loop() {
		Reverse(s)
	}
}

// Fuzz test：平常 go test 只跑 f.Add 的種子；加上 -fuzz 才會自動產生大量隨機輸入。
func FuzzReverse(f *testing.F) {
	for _, seed := range []string{"", "Go", "你好", "!12345"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		if !utf8.ValidString(s) {
			t.Skip("只測合法的 UTF-8")
		}
		if got := Reverse(Reverse(s)); got != s {
			t.Errorf("反轉兩次應該回到原字串：%q 變成 %q", s, got)
		}
	})
}

package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/xinqilin/GO-Playground/chapters/33-idiomatic-go/ticket"
)

func Example_zeroDemo() {
	zeroDemo()
	// Output:
	// 叫號機：2 人在等
	// strings.Builder：排骨飯 x2
	// sync.Mutex：零值是沒有上鎖的鎖
	// nil slice 也能 append：[滷肉飯]
	// nil map 可以讀：0，寫入會 panic
}

func Example_interfaceDemo() {
	interfaceDemo()
	// Output:
	// 取餐櫃台：1 人在等
	// 結帳櫃台：3 人在等
	// 最忙的是：結帳
}

func Example_flowDemo() {
	flowDemo()
	// Output:
	// 請 1 號到結帳櫃台，後面還有 3 位
	// 請 2 號到結帳櫃台
	// 錯誤：沒有指定櫃台
	// 請 3 號到取餐櫃台
	// 請 4 號到取餐櫃台
	// 沒有人在等：取餐櫃台：ticket: 沒有人在等
}

func Example_calledDemo() {
	calledDemo()
	// Output:
	// 改過的複本：[99 2]
	// 叫號機裡的：[1 2]
}

// fakeWaiter 是測試用的櫃台：要幾個人在等就是幾個人，不用真的取號。
type fakeWaiter int

func (f fakeWaiter) Waiting() int { return int(f) }

func TestBusiest(t *testing.T) {
	tests := []struct {
		name     string
		counters map[string]waiter
		want     string
	}{
		{"最多的", map[string]waiter{"結帳": fakeWaiter(2), "取餐": fakeWaiter(5)}, "取餐"},
		{"一樣多取名稱在前的", map[string]waiter{"結帳": fakeWaiter(3), "取餐": fakeWaiter(3)}, "取餐"},
		{"都沒有人", map[string]waiter{"結帳": fakeWaiter(0)}, "結帳"},
		{"沒有櫃台", nil, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := busiest(tt.counters); got != tt.want {
				t.Errorf("busiest = %q，想要 %q", got, tt.want)
			}
		})
	}
}

// failWriter 每次寫入都失敗，用來確認 printBoard 會回傳寫入的錯誤。
type failWriter struct{}

var errDiskFull = errors.New("磁碟滿了")

func (failWriter) Write([]byte) (int, error) { return 0, errDiskFull }

func TestPrintBoard(t *testing.T) {
	var b strings.Builder
	counters := map[string]waiter{"結帳": fakeWaiter(2), "取餐": fakeWaiter(0)}
	if err := printBoard(&b, counters); err != nil {
		t.Fatal(err)
	}
	if want := "取餐櫃台：0 人在等\n結帳櫃台：2 人在等\n"; b.String() != want {
		t.Errorf("printBoard 寫出\n%s想要\n%s", b.String(), want)
	}
	if err := printBoard(failWriter{}, counters); !errors.Is(err, errDiskFull) {
		t.Errorf("寫入失敗時的錯誤 = %v，想要 errDiskFull", err)
	}
}

// TestServeSame 確認兩種寫法在每一種情況下結果都相同。
func TestServeSame(t *testing.T) {
	for _, taken := range []int{0, 1, 2, 5} {
		for _, counter := range []string{"", "結帳"} {
			var a, b ticket.Machine
			for range taken {
				a.Take()
				b.Take()
			}
			for range taken + 1 {
				msgA, errA := serveNested(&a, counter)
				msgB, errB := serve(&b, counter)
				if msgA != msgB || errString(errA) != errString(errB) {
					t.Errorf("取了 %d 張、櫃台 %q：serveNested = %q, %v；serve = %q, %v", taken, counter, msgA, errA, msgB, errB)
				}
			}
		}
	}
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

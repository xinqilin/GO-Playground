package crew

import (
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

// meet 回傳一個函式，給 n 個 goroutine 各呼叫一次：每個呼叫的人都會等到 n 個都到齊才回傳 true。
// 如果工作是一個接一個做的，永遠等不到其他人，等了 timeout 就回傳 false。
func meet(n int, timeout time.Duration) func() bool {
	var mu sync.Mutex
	arrived := 0
	all := make(chan struct{})
	return func() bool {
		mu.Lock()
		arrived++
		if arrived == n {
			close(all)
		}
		mu.Unlock()
		select {
		case <-all:
			return true
		case <-time.After(timeout):
			return false
		}
	}
}

const wait = 500 * time.Millisecond

func TestEach(t *testing.T) {
	const n = 4
	seen := make([]bool, n)
	together := make([]bool, n)
	arrive := meet(n, wait)
	Each(n, func(i int) {
		seen[i] = true
		together[i] = arrive()
	})
	for i := range n {
		if !seen[i] {
			t.Fatalf("f(%d) 沒有被呼叫", i)
		}
		if !together[i] {
			t.Fatalf("f(%d) 等不到其他人：工作沒有同時執行", i)
		}
	}
	Each(0, func(int) { t.Error("n 是 0 時不應該呼叫 f") })
	Each(-1, func(int) { t.Error("n 是負的時不應該呼叫 f") })
}

func TestMapParallel(t *testing.T) {
	arrive := meet(3, wait)
	got := MapParallel([]int{85, 95, 30}, func(p int) string {
		if !arrive() {
			return "沒有同時執行"
		}
		return fmt.Sprintf("NT$%d", p)
	})
	if want := []string{"NT$85", "NT$95", "NT$30"}; !slices.Equal(got, want) {
		t.Errorf("MapParallel = %q, want %q", got, want)
	}
	if got := MapParallel([]int{}, func(int) int { return 0 }); len(got) != 0 {
		t.Errorf("空的輸入 = %v", got)
	}
}

var errSoldOut = errors.New("賣完了")

func TestCheckAll(t *testing.T) {
	arrive := meet(4, wait)
	err := CheckAll([]string{"中山店", "雙連店", "圓山店", "民權店"}, func(s string) error {
		if !arrive() {
			return errors.New("沒有同時執行")
		}
		if s == "雙連店" || s == "民權店" {
			return fmt.Errorf("%s：%w", s, errSoldOut)
		}
		return nil
	})
	if err == nil || err.Error() != "雙連店：賣完了\n民權店：賣完了" {
		t.Errorf("CheckAll = %q, want 雙連店和民權店的錯誤，依照分店的順序", err)
	}
	if !errors.Is(err, errSoldOut) {
		t.Errorf("errors.Is(err, errSoldOut) = false")
	}
	if err := CheckAll([]string{"中山店"}, func(string) error { return nil }); err != nil {
		t.Errorf("全部沒問題時 CheckAll = %v, want nil", err)
	}
}

func TestCountWords(t *testing.T) {
	got := CountWords([]string{"排骨飯 雞腿飯", "排骨飯 紅茶 排骨飯", "", "紅茶"})
	want := map[string]int{"排骨飯": 3, "雞腿飯": 1, "紅茶": 2}
	if !maps.Equal(got, want) {
		t.Errorf("CountWords = %v, want %v", got, want)
	}
	if got := CountWords(nil); len(got) != 0 {
		t.Errorf("CountWords(nil) = %v", got)
	}
}

func TestTry(t *testing.T) {
	ran := false
	if err := Try(func() { ran = true }); err != nil || !ran {
		t.Errorf("正常結束：Try = %v，f 有執行 %v", err, ran)
	}

	err := Try(func() { panic("鍋子破了") })
	if err == nil || err.Error() != "panic: 鍋子破了" {
		t.Errorf("panic 字串：Try = %v, want panic: 鍋子破了", err)
	}

	err = Try(func() { panic(io.ErrUnexpectedEOF) })
	if !errors.Is(err, io.ErrUnexpectedEOF) || !strings.HasPrefix(fmt.Sprint(err), "panic: ") {
		t.Errorf("panic error：Try = %v，要能用 errors.Is 找到原本的錯誤", err)
	}

	err = Try(func() {
		var m map[string]int
		m["排骨飯"]++
	})
	if err == nil || !strings.Contains(err.Error(), "nil map") {
		t.Errorf("執行期錯誤：Try = %v", err)
	}
}

func TestTotal(t *testing.T) {
	tree := &Branch{Name: "總店", Sales: 100, Children: []*Branch{
		{Name: "中山店", Sales: 40, Children: []*Branch{{Name: "雙連店", Sales: 15}}},
		{Name: "圓山店", Sales: 25},
		{Name: "民權店"},
	}}
	if got := Total(tree); got != 180 {
		t.Errorf("Total = %d, want 180", got)
	}
	if got := Total(&Branch{Sales: 7}); got != 7 {
		t.Errorf("沒有子分店時 Total = %d, want 7", got)
	}
	if got := Total(nil); got != 0 {
		t.Errorf("Total(nil) = %d, want 0", got)
	}
}

func TestRunAll(t *testing.T) {
	const n = 3
	arrive := meet(n, wait)
	together := make([]bool, n)
	tasks := make([]func(), n)
	for i := range n {
		tasks[i] = func() { together[i] = arrive() }
	}
	RunAll(tasks)
	for i, ok := range together {
		if !ok {
			t.Fatalf("第 %d 個工作等不到其他人：工作沒有同時執行", i+1)
		}
	}
}

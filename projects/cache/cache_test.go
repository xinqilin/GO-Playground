package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
	"strconv"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

func Example_demo() {
	h := &headOffice{} // 測試裡不用真的等
	demo(context.Background(), os.Stdout, New[string, int](3, 10*time.Minute), h.price)
	// Output:
	// 排骨飯 100 元　從總店　［排骨飯］
	// 雞腿飯 110 元　從總店　［雞腿飯 排骨飯］
	// 排骨飯 100 元　從快取　［排骨飯 雞腿飯］
	// 魚排飯 105 元　從總店　［魚排飯 排骨飯 雞腿飯］
	// 控肉飯 95 元　從總店　［控肉飯 魚排飯 排骨飯］
	// 雞腿飯 110 元　從總店　［雞腿飯 控肉飯 魚排飯］
	// 控肉飯 95 元　從快取　［控肉飯 雞腿飯 魚排飯］
	// 炒泡麵：查不到（炒泡麵：沒有這個品項）
	// 找到 2 次、沒找到 6 次、問總店 6 次、擠掉 2 筆、過期 0 筆
}

func TestLRUOrder(t *testing.T) {
	c := New[string, int](3, 0)
	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)
	c.Get("a")     // a 變成最近用過的
	c.Set("b", 20) // 更新也算用過
	c.Set("d", 4)  // 滿了：丟掉最久沒用的 c
	if got, want := c.Keys(), []string{"d", "b", "a"}; !slices.Equal(got, want) {
		t.Errorf("Keys() = %v，想要 %v", got, want)
	}
	if v, ok := c.Get("b"); !ok || v != 20 {
		t.Errorf("Get(b) = %d, %v，想要 20, true", v, ok)
	}
	if _, ok := c.Get("c"); ok {
		t.Error("c 應該被丟掉了")
	}
	if got := c.Stats().Evictions; got != 1 {
		t.Errorf("Evictions = %d，想要 1", got)
	}
}

// TestTTL：存活時間從 Set 起算，Get 不會延長；時間是假的，十分鐘一瞬間就過了。
func TestTTL(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := New[string, int](10, 10*time.Minute)
		c.Set("排骨飯", 100)

		time.Sleep(5 * time.Minute)
		if _, ok := c.Get("排骨飯"); !ok {
			t.Fatal("5 分鐘時應該還在")
		}
		time.Sleep(5*time.Minute - time.Nanosecond)
		if _, ok := c.Get("排骨飯"); !ok {
			t.Fatal("差 1 奈秒滿 10 分鐘時應該還在")
		}
		time.Sleep(time.Nanosecond)
		if _, ok := c.Get("排骨飯"); ok {
			t.Fatal("滿 10 分鐘時應該過期了")
		}

		c.Set("排骨飯", 105) // 重新存入，重新計時
		time.Sleep(9 * time.Minute)
		if v, ok := c.Get("排骨飯"); !ok || v != 105 {
			t.Errorf("重新存入 9 分鐘後 = %d, %v，想要 105, true", v, ok)
		}
		if s := c.Stats(); s.Expirations != 1 || s.Hits != 3 || s.Misses != 1 {
			t.Errorf("Stats() = %+v", s)
		}
	})
}

func TestNoTTL(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := New[string, int](1, 0)
		c.Set("排骨飯", 100)
		time.Sleep(100 * 365 * 24 * time.Hour)
		if _, ok := c.Get("排骨飯"); !ok {
			t.Error("ttl 是 0 時不應該過期")
		}
	})
}

// TestGetOrLoadOnce：十個人同時查同一種便當，只問總店一次，大家剛好等 300ms。
func TestGetOrLoadOnce(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		h := &headOffice{delay: 300 * time.Millisecond}
		c := New[string, int](3, time.Minute)
		start := time.Now()
		var wg sync.WaitGroup
		for range 10 {
			wg.Go(func() {
				p, err := c.GetOrLoad(context.Background(), "雞腿飯", h.price)
				if p != 110 || err != nil {
					t.Errorf("GetOrLoad = %d, %v", p, err)
				}
			})
		}
		wg.Wait()
		if h.calls != 1 {
			t.Errorf("問了總店 %d 次，想要 1 次", h.calls)
		}
		if d := time.Since(start); d != 300*time.Millisecond {
			t.Errorf("花了 %v，想要 300ms", d)
		}
	})
}

// TestGetOrLoadError：錯誤不存進快取，下一次會再問一次。
func TestGetOrLoadError(t *testing.T) {
	h := &headOffice{}
	c := New[string, int](3, time.Minute)
	for range 2 {
		if _, err := c.GetOrLoad(context.Background(), "炒泡麵", h.price); !errors.Is(err, errNoItem) {
			t.Errorf("err = %v，想要 errNoItem", err)
		}
	}
	if h.calls != 2 || c.Stats().Loads != 2 {
		t.Errorf("問了總店 %d 次、Loads = %d，想要都是 2", h.calls, c.Stats().Loads)
	}
}

// TestGetOrLoadPanic：load panic 時，panic 傳給第一個人，等待的人拿到 errPanicked，之後還能再載入。
func TestGetOrLoadPanic(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := New[string, int](3, time.Minute)
		boom := func(context.Context, string) (int, error) {
			time.Sleep(time.Second)
			panic("總店的系統壞了")
		}

		waiter := make(chan error)
		go func() {
			defer func() { recover() }() // 第一個人：panic 會傳到這裡
			c.GetOrLoad(context.Background(), "排骨飯", boom)
		}()
		synctest.Wait() // 確定第一個人已經開始載入
		go func() {
			_, err := c.GetOrLoad(context.Background(), "排骨飯", boom)
			waiter <- err
		}()
		if err := <-waiter; !errors.Is(err, errPanicked) {
			t.Errorf("等待的人拿到 %v，想要 errPanicked", err)
		}

		h := &headOffice{}
		if p, err := c.GetOrLoad(context.Background(), "排骨飯", h.price); p != 100 || err != nil {
			t.Errorf("再載入一次 = %d, %v，想要 100, nil", p, err)
		}
	})
}

// TestWaiterGivesUp：等待的人可以用自己的 ctx 放棄，載入的人不受影響，結果照樣存起來。
func TestWaiterGivesUp(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		h := &headOffice{delay: 300 * time.Millisecond}
		c := New[string, int](3, time.Minute)
		go c.GetOrLoad(context.Background(), "魚排飯", h.price)
		synctest.Wait()

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if _, err := c.GetOrLoad(ctx, "魚排飯", h.price); !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("err = %v，想要 context.DeadlineExceeded", err)
		}
		time.Sleep(200 * time.Millisecond)
		synctest.Wait()
		if p, ok := c.Get("魚排飯"); !ok || p != 105 {
			t.Errorf("載入完成後 Get = %d, %v，想要 105, true", p, ok)
		}
	})
}

// TestConcurrent 讓很多 goroutine 同時讀寫，用 -race 檢查有沒有漏掉的鎖，最後檢查 map 和串列一致。
func TestConcurrent(t *testing.T) {
	c := New[int, string](50, time.Minute)
	var wg sync.WaitGroup
	for g := range 8 {
		wg.Go(func() {
			r := rand.New(rand.NewPCG(uint64(g), 0))
			for range 2000 {
				k := r.IntN(100)
				switch r.IntN(4) {
				case 0:
					c.Set(k, strconv.Itoa(k))
				case 1:
					c.Delete(k)
				default:
					if v, ok := c.Get(k); ok && v != strconv.Itoa(k) {
						t.Errorf("Get(%d) = %q", k, v)
					}
				}
			}
		})
	}
	wg.Wait()
	if c.order.Len() != len(c.items) || len(c.items) > 50 {
		t.Errorf("串列 %d 筆、map %d 筆，上限 50", c.order.Len(), len(c.items))
	}
}

func TestNewPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("New(0, …) 應該 panic")
		}
	}()
	New[string, int](0, time.Minute)
}

func ExampleCache_Keys() {
	c := New[string, int](3, 0)
	for _, item := range []string{"排骨飯", "雞腿飯", "魚排飯"} {
		c.Set(item, 0)
	}
	fmt.Println(c.Keys())
	c.Get("排骨飯")
	fmt.Println(c.Keys())
	c.Set("控肉飯", 95)
	fmt.Println(c.Keys())
	// Output:
	// [魚排飯 雞腿飯 排骨飯]
	// [排骨飯 魚排飯 雞腿飯]
	// [控肉飯 排骨飯 魚排飯]
}

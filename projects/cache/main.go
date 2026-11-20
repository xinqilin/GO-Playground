// 專題 7：TTL＋LRU 快取。
//
// 分店要查總店的價目表，每查一次要等 300ms。在中間放一個快取：最多記 3 種便當，
// 每筆 10 分鐘後過期；同一種便當同時有很多人查，只會問總店一次。
//
// 執行：go run ./projects/cache
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// prices 是總店的價目表（示範用的資料）。
var prices = map[string]int{"排骨飯": 100, "雞腿飯": 110, "魚排飯": 105, "控肉飯": 95, "素食餐": 90}

var errNoItem = errors.New("沒有這個品項")

// headOffice 模擬總店的價目表：每查一次要等 delay。
type headOffice struct {
	delay time.Duration
	mu    sync.Mutex
	calls int
}

func (h *headOffice) price(ctx context.Context, item string) (int, error) {
	h.mu.Lock()
	h.calls++
	h.mu.Unlock()
	select {
	case <-time.After(h.delay):
	case <-ctx.Done():
		return 0, ctx.Err()
	}
	p, ok := prices[item]
	if !ok {
		return 0, fmt.Errorf("%s：%w", item, errNoItem)
	}
	return p, nil
}

// demo 依序查一串便當，印出價格、是不是從快取拿到的、查完之後快取裡由新到舊的順序。
func demo(ctx context.Context, w io.Writer, c *Cache[string, int], load func(context.Context, string) (int, error)) {
	for _, item := range []string{"排骨飯", "雞腿飯", "排骨飯", "魚排飯", "控肉飯", "雞腿飯", "控肉飯", "炒泡麵"} {
		loads := c.Stats().Loads
		p, err := c.GetOrLoad(ctx, item, load)
		if err != nil {
			fmt.Fprintf(w, "%s：查不到（%v）\n", item, err)
			continue
		}
		from := "快取"
		if c.Stats().Loads > loads {
			from = "總店"
		}
		fmt.Fprintf(w, "%s %d 元　從%s　［%s］\n", item, p, from, strings.Join(c.Keys(), " "))
	}
	s := c.Stats()
	fmt.Fprintf(w, "找到 %d 次、沒找到 %d 次、問總店 %d 次、擠掉 %d 筆、過期 %d 筆\n",
		s.Hits, s.Misses, s.Loads, s.Evictions, s.Expirations)
}

func main() {
	h := &headOffice{delay: 300 * time.Millisecond}
	start := time.Now()
	demo(context.Background(), os.Stdout, New[string, int](3, 10*time.Minute), h.price)
	fmt.Printf("總共花了 %v（問總店 %d 次，每次 %v）\n", time.Since(start).Round(100*time.Millisecond), h.calls, h.delay)
}

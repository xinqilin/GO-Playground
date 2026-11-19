package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// fetchAll 為每一家分店開一個 goroutine 查庫存，但同時最多 limit 個在查。
// sem 是容量 limit 的 channel：放得進去才能開始，做完拿出來，就像有限的號碼牌。
func fetchAll(branches []string, limit int, fetch func(string) int) []int {
	sem := make(chan struct{}, limit)
	out := make([]int, len(branches))
	var wg sync.WaitGroup
	for i, b := range branches {
		wg.Go(func() {
			sem <- struct{}{}
			defer func() { <-sem }()
			out[i] = fetch(b)
		})
	}
	wg.Wait()
	return out
}

// gauge 記錄同時有幾個工作在跑，以及最多的時候是幾個。
type gauge struct {
	running, peak atomic.Int64
}

func (g *gauge) enter() {
	now := g.running.Add(1)
	for {
		p := g.peak.Load()
		if now <= p || g.peak.CompareAndSwap(p, now) {
			return
		}
	}
}

func (g *gauge) leave() { g.running.Add(-1) }

func semaphoreDemo() {
	stock := map[string]int{"東區": 12, "西區": 0, "南區": 7, "北區": 30, "中區": 5, "內湖": 18, "板橋": 9}
	branches := []string{"東區", "西區", "南區", "北區", "中區", "內湖", "板橋"}
	var g gauge
	got := fetchAll(branches, 2, func(b string) int {
		g.enter()
		defer g.leave()
		return stock[b]
	})
	fmt.Println(got)
	fmt.Printf("同時在查的分店沒有超過 2 家：%t\n", g.peak.Load() <= 2)
}

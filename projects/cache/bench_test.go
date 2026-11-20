package main

import (
	"strconv"
	"sync"
	"testing"
)

// keys 是 benchmark 共用的 1000 個 key，先建好，不把轉字串的時間算進去。
var keys = func() []string {
	k := make([]string, 1000)
	for i := range k {
		k[i] = "item" + strconv.Itoa(i)
	}
	return k
}()

func BenchmarkGet(b *testing.B) {
	c := New[string, int](len(keys), 0)
	for i, k := range keys {
		c.Set(k, i)
	}
	i := 0
	for b.Loop() {
		c.Get(keys[i%len(keys)])
		i++
	}
}

// BenchmarkSetEvict：每次 Set 都是新的 key，而且都要丟掉一筆舊的。
func BenchmarkSetEvict(b *testing.B) {
	c := New[string, int](100, 0)
	i := 0
	for b.Loop() {
		c.Set(keys[i%len(keys)], i)
		i++
	}
}

// BenchmarkGetParallel：很多 goroutine 同時 Get，全部搶同一把鎖。
func BenchmarkGetParallel(b *testing.B) {
	c := New[string, int](len(keys), 0)
	for i, k := range keys {
		c.Set(k, i)
	}
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			c.Get(keys[i%len(keys)])
			i++
		}
	})
}

// BenchmarkMapParallel 是對照組：只有 map 加一把鎖，沒有 LRU 和 TTL。
func BenchmarkMapParallel(b *testing.B) {
	var mu sync.Mutex
	m := map[string]int{}
	for i, k := range keys {
		m[k] = i
	}
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			mu.Lock()
			_ = m[keys[i%len(keys)]]
			mu.Unlock()
			i++
		}
	})
}

package main

import (
	"container/list"
	"context"
	"errors"
	"sync"
	"time"
)

// Cache 是有筆數上限（LRU：滿了丟掉最久沒用的）和存活時間（TTL）的快取，
// 可以同時給很多 goroutine 使用。用 New 建立。
type Cache[K comparable, V any] struct {
	mu      sync.Mutex
	size    int
	ttl     time.Duration
	order   *list.List          // 最近用過的在最前面；每個元素的 Value 是 *entry[K, V]
	items   map[K]*list.Element // key → order 裡的元素
	loading map[K]*call[V]      // 正在載入的 key
	stats   Stats
}

type entry[K comparable, V any] struct {
	key     K // 從串列尾端丟掉一筆時，要用它刪掉 map 裡的 key
	val     V
	expires time.Time // 零值表示不會過期
}

// call 是一次正在進行的載入。done 關閉之後，val 和 err 才能讀。
type call[V any] struct {
	done chan struct{}
	val  V
	err  error
}

// Stats 是快取的統計。
type Stats struct {
	Hits        int // 有找到
	Misses      int // 沒找到，包括過期的
	Loads       int // GetOrLoad 呼叫 load 的次數
	Evictions   int // 滿了被丟掉的筆數
	Expirations int // 過期被丟掉的筆數
}

// errPanicked 是 load 發生 panic 時，等待同一個 key 的其他人拿到的錯誤。
var errPanicked = errors.New("cache: load 發生 panic")

// New 建立最多存 size 筆、每筆存活 ttl 的快取。ttl 是 0 表示不會過期。size 必須大於 0。
func New[K comparable, V any](size int, ttl time.Duration) *Cache[K, V] {
	if size < 1 {
		panic("cache: size 必須大於 0")
	}
	return &Cache[K, V]{
		size:    size,
		ttl:     ttl,
		order:   list.New(),
		items:   map[K]*list.Element{},
		loading: map[K]*call[V]{},
	}
}

// Get 回傳 key 的值；沒有或已經過期時，第二個回傳值是 false。
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.getLocked(key)
}

// Set 存入 key 的值，存活時間從現在重新算。超過筆數上限時，丟掉最久沒用的那一筆。
func (c *Cache[K, V]) Set(key K, val V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.setLocked(key, val)
}

// Delete 刪掉 key。
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if el, ok := c.items[key]; ok {
		c.removeLocked(el)
	}
}

// Keys 由最近用過到最久沒用，回傳目前所有的 key（包括已經過期、還沒被丟掉的）。
func (c *Cache[K, V]) Keys() []K {
	c.mu.Lock()
	defer c.mu.Unlock()
	keys := make([]K, 0, c.order.Len())
	for el := c.order.Front(); el != nil; el = el.Next() {
		keys = append(keys, el.Value.(*entry[K, V]).key)
	}
	return keys
}

// Stats 回傳目前的統計。
func (c *Cache[K, V]) Stats() Stats {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.stats
}

func (c *Cache[K, V]) getLocked(key K) (V, bool) {
	el, ok := c.items[key]
	if !ok {
		c.stats.Misses++
		var zero V
		return zero, false
	}
	e := el.Value.(*entry[K, V])
	if !e.expires.IsZero() && !time.Now().Before(e.expires) {
		c.removeLocked(el) // 過期的資料用到時才丟掉
		c.stats.Expirations++
		c.stats.Misses++
		var zero V
		return zero, false
	}
	c.order.MoveToFront(el)
	c.stats.Hits++
	return e.val, true
}

func (c *Cache[K, V]) setLocked(key K, val V) {
	var expires time.Time
	if c.ttl > 0 {
		expires = time.Now().Add(c.ttl)
	}
	if el, ok := c.items[key]; ok {
		e := el.Value.(*entry[K, V])
		e.val, e.expires = val, expires
		c.order.MoveToFront(el)
		return
	}
	c.items[key] = c.order.PushFront(&entry[K, V]{key: key, val: val, expires: expires})
	if c.order.Len() > c.size {
		c.removeLocked(c.order.Back())
		c.stats.Evictions++
	}
}

func (c *Cache[K, V]) removeLocked(el *list.Element) {
	e := c.order.Remove(el).(*entry[K, V])
	delete(c.items, e.key)
}

// GetOrLoad 回傳 key 的值；沒有或已經過期時，呼叫 load 取得並存起來。load 回傳錯誤時不存。
// 同一個 key 同時有很多人沒找到，只有第一個人會呼叫 load，其他人等他的結果；
// 等待的人可以用自己的 ctx 放棄等待。
func (c *Cache[K, V]) GetOrLoad(ctx context.Context, key K, load func(context.Context, K) (V, error)) (V, error) {
	c.mu.Lock()
	if v, ok := c.getLocked(key); ok {
		c.mu.Unlock()
		return v, nil
	}
	if cl, ok := c.loading[key]; ok { // 已經有人在載入了：等他
		c.mu.Unlock()
		select {
		case <-cl.done:
			return cl.val, cl.err
		case <-ctx.Done():
			var zero V
			return zero, ctx.Err()
		}
	}
	cl := &call[V]{done: make(chan struct{}), err: errPanicked} // load 沒有正常回傳時，err 就留在 errPanicked
	c.loading[key] = cl
	c.stats.Loads++
	c.mu.Unlock() // load 可能很慢，不能拿著鎖等

	defer func() { // load panic 時也會執行，等待的人不會永遠卡住
		c.mu.Lock()
		delete(c.loading, key)
		if cl.err == nil {
			c.setLocked(key, cl.val)
		}
		c.mu.Unlock()
		close(cl.done)
	}()
	cl.val, cl.err = load(ctx, key)
	return cl.val, cl.err
}

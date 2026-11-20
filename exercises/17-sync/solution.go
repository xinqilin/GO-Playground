//go:build solution

package stall

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Jar 用一個 Mutex 同時保護 times 和 total。
// 用兩個 atomic.Int64 的話，每個數字各自正確，但 Count 可能讀到一個加了、另一個還沒加的狀態。
type Jar struct {
	mu    sync.Mutex
	times int
	total int
}

func (j *Jar) Put(amount int) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.times++
	j.total += amount
}

func (j *Jar) Count() (times, total int) {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.times, j.total
}

// Shutter 只有一個 bool，用 atomic.Bool 就夠了。
// CompareAndSwap(false, true) 只有在「目前是 false」時才會改成 true 並回傳 true，
// 同時呼叫的人裡只會有一個成功。
type Shutter struct {
	closed atomic.Bool
}

func (s *Shutter) Close() bool {
	return s.closed.CompareAndSwap(false, true)
}

func (s *Shutter) Closed() bool {
	return s.closed.Load()
}

// Reservations 的「查有沒有人訂」和「寫入」在同一段鎖裡。
type Reservations struct {
	mu     sync.Mutex
	tables map[int]string
}

func NewReservations() *Reservations {
	return &Reservations{tables: make(map[int]string)}
}

func (r *Reservations) Book(table int, name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tables[table]; ok {
		return ErrTaken
	}
	r.tables[table] = name
	return nil
}

func (r *Reservations) Who(table int) (string, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	name, ok := r.tables[table]
	return name, ok
}

// Current 用 atomic.Pointer：每次 Update 都做一份新的 Settings，用 CompareAndSwap 換上去。
// 換的時候如果發現別人已經先換過了（指標不是剛剛讀到的那個），就用新的設定重算一次。
// 用 Mutex 也對，這裡示範讀的人完全不用等鎖的寫法。
type Current struct {
	p atomic.Pointer[Settings]
}

func (c *Current) Get() Settings {
	if s := c.p.Load(); s != nil {
		return *s
	}
	return Settings{}
}

func (c *Current) Update(fn func(Settings) Settings) {
	for {
		old := c.p.Load()
		var cur Settings
		if old != nil {
			cur = *old
		}
		next := fn(cur)
		if c.p.CompareAndSwap(old, &next) {
			return
		}
	}
}

// OnceOK 呼叫 f 的整段都拿著鎖：同時呼叫的人排隊，同一時間只有一個 f 在執行；
// 前一個人成功之後，後面的人拿到鎖就直接回傳存好的值。
func OnceOK[T any](f func() (T, error)) func() (T, error) {
	var (
		mu   sync.Mutex
		done bool
		val  T
	)
	return func() (T, error) {
		mu.Lock()
		defer mu.Unlock()
		if done {
			return val, nil
		}
		v, err := f()
		if err != nil {
			return v, err
		}
		val, done = v, true
		return val, nil
	}
}

// Summary 原本在拿著 q.mu 的時候呼叫 q.Len()，Len 又要拿 q.mu。
// Go 的 Mutex 不能重複進入：同一個 goroutine 再 Lock 一次，會等自己 Unlock，永遠等不到。
// 已經拿著鎖的地方，直接讀欄位就好。
func (q *Queue) Summary() string {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.names) == 0 {
		return "目前不用等"
	}
	return fmt.Sprintf("%d 組候位，下一組：%s", len(q.names), q.names[0])
}

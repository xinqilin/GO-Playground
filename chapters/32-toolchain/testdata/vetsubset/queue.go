// Package vetsubset 的 Next 方法用了值 receiver，複製了整把鎖。
// go test 會通過，go vet 才抓得到。
package vetsubset

import "sync"

// Queue 是排隊叫號機。
type Queue struct {
	mu   sync.Mutex
	next int
}

// Next 回傳下一個號碼。
func (q Queue) Next() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.next++
	return q.next
}

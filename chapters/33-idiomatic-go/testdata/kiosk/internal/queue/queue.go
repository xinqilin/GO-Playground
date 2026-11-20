// Package queue 管理排隊的號碼。只有 example.com/kiosk 裡的程式可以 import 它。
package queue

// Queue 是一條隊伍。零值是空的隊伍。
type Queue struct {
	issued, called int
}

// Take 取一張號碼牌。
func (q *Queue) Take() int {
	q.issued++
	return q.issued
}

// Call 叫下一號；沒有人在等時 ok 是 false。
func (q *Queue) Call() (n int, ok bool) {
	if q.called == q.issued {
		return 0, false
	}
	q.called++
	return q.called, true
}

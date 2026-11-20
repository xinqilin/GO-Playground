// Package ticket 是排隊叫號機：取號、叫號，以及查詢還有幾個人在等。
//
// [Machine] 的零值就可以直接使用，也可以同時被多個 goroutine 使用：
//
//	var m ticket.Machine
//	n := m.Take()
//
// # 號碼
//
// 號碼從 1 開始依序發出，叫號的順序和取號的順序相同。叫過的號碼可以用 [Machine.Called] 查詢。
package ticket

import (
	"errors"
	"slices"
	"sync"
)

// ErrNoOneWaiting 表示目前沒有人在等，叫不到號。
var ErrNoOneWaiting = errors.New("ticket: 沒有人在等")

// Machine 是一台叫號機。零值是一台還沒有人取號的叫號機。Machine 開始使用後不能再複製。
type Machine struct {
	mu     sync.Mutex
	issued int
	called []int
}

// Take 取一張號碼牌，回傳它的號碼。
func (m *Machine) Take() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.issued++
	return m.issued
}

// Call 叫下一號並回傳號碼。沒有人在等時回傳 [ErrNoOneWaiting]。
func (m *Machine) Call() (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.called) == m.issued {
		return 0, ErrNoOneWaiting
	}
	n := len(m.called) + 1
	m.called = append(m.called, n)
	return n, nil
}

// Waiting 回傳還在等的人數。
func (m *Machine) Waiting() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.issued - len(m.called)
}

// Called 回傳叫過的號碼，依叫號的順序。回傳的是複本，修改它不會影響 m。
func (m *Machine) Called() []int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return slices.Clone(m.called)
}

// Remaining 回傳還在等的人數。
//
// Deprecated: 改用 [Machine.Waiting]。
func (m *Machine) Remaining() int {
	return m.Waiting()
}

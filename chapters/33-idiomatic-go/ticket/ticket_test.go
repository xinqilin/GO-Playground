package ticket_test

import (
	"errors"
	"fmt"
	"slices"
	"sync"
	"testing"

	"github.com/xinqilin/GO-Playground/chapters/33-idiomatic-go/ticket"
)

func ExampleMachine() {
	var m ticket.Machine
	m.Take()
	m.Take()
	n, err := m.Call()
	fmt.Println(n, err, m.Waiting())
	// Output: 1 <nil> 1
}

func TestCallEmpty(t *testing.T) {
	var m ticket.Machine
	if _, err := m.Call(); !errors.Is(err, ticket.ErrNoOneWaiting) {
		t.Errorf("空的叫號機 Call 的錯誤 = %v，想要 ErrNoOneWaiting", err)
	}
}

func TestCalledIsCopy(t *testing.T) {
	var m ticket.Machine
	m.Take()
	m.Call()
	got := m.Called()
	got[0] = 99
	if again := m.Called(); !slices.Equal(again, []int{1}) {
		t.Errorf("改了回傳的 slice 之後，Called() = %v，想要 [1]", again)
	}
}

func TestConcurrent(t *testing.T) {
	var m ticket.Machine
	var wg sync.WaitGroup
	for range 50 {
		wg.Go(func() { m.Take() })
	}
	wg.Wait()
	for range 50 {
		wg.Go(func() { m.Call() })
	}
	wg.Wait()
	if !slices.IsSorted(m.Called()) || len(m.Called()) != 50 || m.Waiting() != 0 {
		t.Errorf("Called() = %v，Waiting() = %d", m.Called(), m.Waiting())
	}
}

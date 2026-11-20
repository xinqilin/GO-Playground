package vetsubset

import "testing"

func TestNext(t *testing.T) {
	var q Queue
	if got := q.Next(); got != 1 {
		t.Errorf("第一個號碼 = %d，想要 1", got)
	}
}

package main

import "fmt"

// Set 是泛型型別：T 是元素的型別，要能當 map 的 key，所以限制是 comparable。
// 零值就能用：第一次 Add 時才建立 map。
type Set[T comparable] struct {
	m map[T]struct{}
}

func (s *Set[T]) Add(v T) {
	if s.m == nil {
		s.m = map[T]struct{}{}
	}
	s.m[v] = struct{}{}
}

func (s *Set[T]) Has(v T) bool {
	_, ok := s.m[v]
	return ok
}

func (s *Set[T]) Len() int {
	return len(s.m)
}

func setDemo() {
	var tables Set[int]
	tables.Add(3)
	tables.Add(5)
	tables.Add(3)
	fmt.Println(tables.Len(), tables.Has(5), tables.Has(7))

	var dishes Set[string]
	dishes.Add("排骨飯")
	fmt.Println(dishes.Has("排骨飯"), dishes.Len())
}

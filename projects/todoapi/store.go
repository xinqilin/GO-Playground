package main

import (
	"errors"
	"slices"
	"sync"
)

// Todo 是一筆待辦事項。
type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// ErrNotFound 表示找不到指定的待辦事項。
var ErrNotFound = errors.New("todo not found")

// Store 把待辦事項放在記憶體裡。所有方法都可以被多個 goroutine 同時呼叫。
type Store struct {
	mu     sync.Mutex
	todos  map[int]Todo
	nextID int
}

// NewStore 建立一個空的 Store。
func NewStore() *Store {
	return &Store{todos: make(map[int]Todo), nextID: 1}
}

// List 依 ID 由小到大回傳全部待辦事項；沒有資料時回傳空 slice（JSON 是 []，不是 null）。
func (s *Store) List() []Todo {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Todo, 0, len(s.todos))
	for _, t := range s.todos {
		out = append(out, t)
	}
	slices.SortFunc(out, func(a, b Todo) int { return a.ID - b.ID })
	return out
}

// Get 依 ID 取得一筆待辦事項。
func (s *Store) Get(id int) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.todos[id]
	if !ok {
		return Todo{}, ErrNotFound
	}
	return t, nil
}

// Create 新增一筆待辦事項並回傳它。
func (s *Store) Create(title string) Todo {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := Todo{ID: s.nextID, Title: title}
	s.todos[t.ID] = t
	s.nextID++
	return t
}

// SetDone 更新完成狀態。map 裡的值不能取地址，所以要取出、修改、再放回去。
func (s *Store) SetDone(id int, done bool) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.todos[id]
	if !ok {
		return Todo{}, ErrNotFound
	}
	t.Done = done
	s.todos[id] = t
	return t, nil
}

// Delete 刪除一筆待辦事項。
func (s *Store) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.todos[id]; !ok {
		return ErrNotFound
	}
	delete(s.todos, id)
	return nil
}

package main

import (
	"cmp"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sync"
)

// Store 把待辦事項放在記憶體裡；path 不是空字串時，每次修改都把整份資料寫進這個 JSON 檔。
// 所有方法都可以被多個 goroutine 同時呼叫。
type Store struct {
	mu     sync.Mutex
	path   string
	todos  map[int]Todo
	nextID int
}

// fileData 是存檔的格式。next_id 也要存：只存清單的話，刪掉編號最大的一筆再重新啟動，
// 下一筆新的待辦就會拿到被刪掉的那個編號。
type fileData struct {
	NextID int    `json:"next_id"`
	Todos  []Todo `json:"todos"`
}

// OpenStore 讀取 path 裡的資料；檔案還不存在時從空的開始。path 是空字串時只放在記憶體。
func OpenStore(path string) (*Store, error) {
	s := &Store{path: path, todos: map[int]Todo{}, nextID: 1}
	if path == "" {
		return s, nil
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	var fd fileData
	if err := json.Unmarshal(data, &fd); err != nil {
		return nil, fmt.Errorf("讀取 %s：%w", path, err)
	}
	s.nextID = max(fd.NextID, 1)
	for _, t := range fd.Todos {
		s.todos[t.ID] = t
		s.nextID = max(s.nextID, t.ID+1) // next_id 不見了或被改小時，也不能發出已經有人用的編號
	}
	return s, nil
}

// List 依 ID 由小到大回傳全部的待辦。
func (s *Store) List() []Todo {
	s.mu.Lock()
	defer s.mu.Unlock()
	return sortedTodos(s.todos)
}

// Get 依 ID 取得一筆待辦。
func (s *Store) Get(id int) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.todos[id]
	if !ok {
		return Todo{}, ErrNotFound
	}
	return t, nil
}

// Add 給 t 一個新的 ID，存起來之後回傳。
func (s *Store) Add(t Todo) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t.ID = s.nextID
	next := maps.Clone(s.todos)
	next[t.ID] = t
	if err := s.commit(next, s.nextID+1); err != nil {
		return Todo{}, err
	}
	return t, nil
}

// Update 在鎖裡把 id 那一筆交給 change 修改，再存回去。
// 讀出來、修改、寫回去是同一步，兩個同時進來的修改不會互相蓋掉。
func (s *Store) Update(id int, change func(*Todo)) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.todos[id]
	if !ok {
		return Todo{}, ErrNotFound
	}
	change(&t)
	next := maps.Clone(s.todos)
	next[id] = t
	if err := s.commit(next, s.nextID); err != nil {
		return Todo{}, err
	}
	return t, nil
}

// Delete 刪除一筆待辦。
func (s *Store) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.todos[id]; !ok {
		return ErrNotFound
	}
	next := maps.Clone(s.todos)
	delete(next, id)
	return s.commit(next, s.nextID)
}

// commit 先把新的狀態寫進檔案，成功了才換掉記憶體裡的資料。
// 寫檔失敗時兩邊都維持原樣，不會出現「API 說成功了，重新啟動後卻不見了」。呼叫的一方要拿著 s.mu。
func (s *Store) commit(todos map[int]Todo, nextID int) error {
	if s.path != "" {
		if err := saveAtomic(s.path, fileData{NextID: nextID, Todos: sortedTodos(todos)}); err != nil {
			return fmt.Errorf("存檔：%w", err)
		}
	}
	s.todos, s.nextID = todos, nextID
	return nil
}

func sortedTodos(todos map[int]Todo) []Todo {
	return slices.SortedFunc(maps.Values(todos), func(a, b Todo) int { return cmp.Compare(a.ID, b.ID) })
}

// saveAtomic 把 v 編成 JSON，先寫進同一個目錄裡的暫存檔，成功了才用 Rename 換掉 path。
// 中途出錯就刪掉暫存檔，原本的檔案完全沒有動到。
func saveAtomic(path string, v any) (err error) {
	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tmp.Close()
			os.Remove(tmp.Name())
		}
	}()
	if err = json.MarshalWrite(tmp, v, jsontext.WithIndent("  ")); err != nil {
		return err
	}
	if err = tmp.Sync(); err != nil { // 確定資料寫到磁碟上了，再換名字
		return err
	}
	if err = tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

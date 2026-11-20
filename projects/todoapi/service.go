package main

import (
	"fmt"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

// maxTitle 是標題最多幾個字（以 rune 計算，中文一個字算一個）。
const maxTitle = 100

// service 放這個服務的規則：標題怎麼整理、什麼樣的輸入合法、建立時間從哪裡來。
// 它不知道 HTTP，也不知道資料存在哪裡。
type service struct {
	store *Store
	now   func() time.Time // 測試時換成固定的時間
}

// NewTodo 是建立待辦時的輸入。
type NewTodo struct {
	Title string `json:"title"`
	Due   string `json:"due"`
}

// Patch 是修改時的輸入。nil 代表這個欄位不改：用指標才分得出「沒有傳」和「傳了零值」。
type Patch struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
	Due   *string `json:"due"` // 空字串代表清除期限
}

// Create 整理、檢查輸入，再交給 store 存起來。
func (s *service) Create(in NewTodo) (Todo, error) {
	title := strings.TrimSpace(in.Title)
	if err := s.check(&title, &in.Due); err != nil {
		return Todo{}, err
	}
	return s.store.Add(Todo{Title: title, Due: in.Due, CreatedAt: s.now()})
}

// Update 只檢查有傳進來的欄位：只想打勾的請求，不會因為期限已經過了而失敗。
func (s *service) Update(id int, p Patch) (Todo, error) {
	if p == (Patch{}) {
		return Todo{}, ValidationError{"body": "至少要改一個欄位"}
	}
	if p.Title != nil {
		title := strings.TrimSpace(*p.Title)
		p.Title = &title
	}
	if err := s.check(p.Title, p.Due); err != nil {
		return Todo{}, err
	}
	return s.store.Update(id, func(t *Todo) {
		if p.Title != nil {
			t.Title = *p.Title
		}
		if p.Done != nil {
			t.Done = *p.Done
		}
		if p.Due != nil {
			t.Due = *p.Due
		}
	})
}

// List 回傳全部的待辦；done 不是 nil 時只留下完成狀態相同的。
func (s *service) List(done *bool) []Todo {
	todos := s.store.List()
	if done != nil {
		todos = slices.DeleteFunc(todos, func(t Todo) bool { return t.Done != *done })
	}
	return todos
}

// Get 和 Delete 沒有規則，直接交給 store。
func (s *service) Get(id int) (Todo, error) { return s.store.Get(id) }

func (s *service) Delete(id int) error { return s.store.Delete(id) }

// check 檢查不是 nil 的欄位，把所有問題一次收集起來；沒有問題時回傳 nil。
func (s *service) check(title, due *string) error {
	problems := ValidationError{}
	if title != nil {
		switch n := utf8.RuneCountInString(*title); {
		case n == 0:
			problems["title"] = "不能是空的"
		case n > maxTitle:
			problems["title"] = fmt.Sprintf("最多 %d 個字，現在是 %d 個字", maxTitle, n)
		}
	}
	if due != nil && *due != "" {
		now := s.now()
		d, err := time.ParseInLocation(time.DateOnly, *due, now.Location())
		y, m, day := now.Date()
		switch {
		case err != nil:
			problems["due"] = "格式要是 2006-01-02"
		case d.Before(time.Date(y, m, day, 0, 0, 0, 0, now.Location())):
			problems["due"] = "不能早於今天"
		}
	}
	if len(problems) > 0 {
		return problems
	}
	return nil
}

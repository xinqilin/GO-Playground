package main

import (
	"errors"
	"maps"
	"strings"
	"sync"
	"testing"
	"time"
)

func newTestService(now func() time.Time) *service {
	store, err := OpenStore("")
	if err != nil {
		panic(err)
	}
	return &service{store: store, now: now}
}

func TestCreateValidation(t *testing.T) {
	tests := []struct {
		in   NewTodo
		want ValidationError // nil 代表合法
	}{
		{NewTodo{Title: "訂便當"}, nil},
		{NewTodo{Title: "  訂便當  ", Due: "2026-11-21"}, nil}, // 今天也可以
		{NewTodo{Title: strings.Repeat("飯", 100)}, nil},
		{NewTodo{Title: strings.Repeat("飯", 101)}, ValidationError{"title": "最多 100 個字，現在是 101 個字"}},
		{NewTodo{Title: " \t "}, ValidationError{"title": "不能是空的"}},
		{NewTodo{Title: "訂便當", Due: "2026-11-20"}, ValidationError{"due": "不能早於今天"}},
		{NewTodo{Title: "訂便當", Due: "2026/11/30"}, ValidationError{"due": "格式要是 2006-01-02"}},
		{NewTodo{Title: "訂便當", Due: "2026-02-30"}, ValidationError{"due": "格式要是 2006-01-02"}},
		{NewTodo{Title: "", Due: "明天"}, ValidationError{"title": "不能是空的", "due": "格式要是 2006-01-02"}},
	}
	for _, tt := range tests {
		got, err := newTestService(fixedNow).Create(tt.in)
		ve, _ := errors.AsType[ValidationError](err)
		if !maps.Equal(ve, tt.want) || (tt.want == nil && err != nil) {
			t.Errorf("Create(%+v) 的錯誤 = %v，想要 %v", tt.in, err, tt.want)
			continue
		}
		if err == nil && (got.Title != strings.TrimSpace(tt.in.Title) || !got.CreatedAt.Equal(fixedNow())) {
			t.Errorf("Create(%+v) = %+v", tt.in, got)
		}
	}
}

// TestUpdateChecksOnlyGivenFields：期限過了之後，打勾還是可以；但不能把期限改成過去的日期。
func TestUpdateChecksOnlyGivenFields(t *testing.T) {
	now := fixedNow()
	svc := newTestService(func() time.Time { return now })
	if _, err := svc.Create(NewTodo{Title: "訂便當", Due: "2026-11-22"}); err != nil {
		t.Fatal(err)
	}
	now = now.AddDate(0, 0, 3) // 三天後，期限已經過了

	done := true
	if _, err := svc.Update(1, Patch{Done: &done}); err != nil {
		t.Errorf("只打勾：%v", err)
	}
	due := "2026-11-23"
	if _, err := svc.Update(1, Patch{Due: &due}); err == nil {
		t.Error("期限改成過去的日期，應該失敗")
	}
	none := ""
	if got, err := svc.Update(1, Patch{Due: &none}); err != nil || got.Due != "" {
		t.Errorf("清除期限：%+v, %v", got, err)
	}
}

// TestConcurrentPatch：同一筆待辦，一半的請求改標題、一半的請求打勾。
// 讀出來、修改、寫回去都在同一把鎖裡，兩種修改都要留下來。
func TestConcurrentPatch(t *testing.T) {
	svc := newTestService(fixedNow)
	if _, err := svc.Create(NewTodo{Title: "訂便當"}); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := range 100 {
		wg.Go(func() {
			title, done := "訂排骨便當", true
			p := Patch{Title: &title}
			if i%2 == 0 {
				p = Patch{Done: &done}
			}
			if _, err := svc.Update(1, p); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	if got, _ := svc.Get(1); got.Title != "訂排骨便當" || !got.Done {
		t.Errorf("got %+v，標題和完成狀態都應該改到", got)
	}
}

func TestListFilter(t *testing.T) {
	svc := newTestService(fixedNow)
	for _, title := range []string{"訂便當", "繳電話費", "倒垃圾"} {
		svc.Create(NewTodo{Title: title})
	}
	done := true
	svc.Update(2, Patch{Done: &done})

	notDone := false
	for _, tt := range []struct {
		done *bool
		want int
	}{{nil, 3}, {&done, 1}, {&notDone, 2}} {
		if got := svc.List(tt.done); len(got) != tt.want {
			t.Errorf("List(%v) 有 %d 筆，想要 %d 筆", tt.done, len(got), tt.want)
		}
	}
	if got := svc.List(nil); len(got) != 3 {
		t.Errorf("篩選之後，原本的資料不應該少掉：%d 筆", len(got))
	}
}

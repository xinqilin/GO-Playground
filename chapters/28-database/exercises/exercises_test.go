package library

import (
	"database/sql"
	"errors"
	"path/filepath"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// newLib 給每個測試一個全新的資料庫檔案。
func newLib(t *testing.T) *Library {
	t.Helper()
	l, err := Open(t.Context(), filepath.Join(t.TempDir(), "library.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	return l
}

// exec 直接執行 SQL 準備測試資料，不依賴其他題目的答案。
func exec(t *testing.T, l *Library, query string, args ...any) {
	t.Helper()
	if _, err := l.db.ExecContext(t.Context(), query, args...); err != nil {
		t.Fatal(err)
	}
}

// seedBooks 放好 5 本書，id 是 1 到 5。
func seedBooks(t *testing.T, l *Library) {
	t.Helper()
	exec(t, l, `INSERT INTO books (title, author, isbn) VALUES
		('雨天的公車', '陳小雨', NULL),
		('巷口的貓', '陳小雨', '978-0-00-000001-1'),
		('夜市人生', '王大明', NULL),
		('Tea for Two', 'O''Neil', NULL),
		('溪邊的石頭', '林阿美', NULL)`)
}

func TestAddBook(t *testing.T) {
	l := newLib(t)
	ctx := t.Context()
	steps := []struct {
		title, author, isbn string
		wantID              int64
	}{
		{"巷口的貓", "陳小雨", "", 1},
		{"夜市人生", "王大明", "", 2}, // 第二本沒有 ISBN 的書：存 NULL 才不會違反 UNIQUE
		{"雨天的公車", "陳小雨", "978-0-00-000001-1", 3},
	}
	for _, s := range steps {
		id, err := l.AddBook(ctx, s.title, s.author, s.isbn)
		if err != nil || id != s.wantID {
			t.Errorf("AddBook(%q, isbn %q) = %d, %v，want %d, nil", s.title, s.isbn, id, err, s.wantID)
		}
	}
	if _, err := l.AddBook(ctx, "別的書", "某人", "978-0-00-000001-1"); err == nil {
		t.Error("重複的 ISBN 應該回傳錯誤")
	}

	var nulls int
	if err := l.db.QueryRowContext(ctx, `SELECT count(*) FROM books WHERE isbn IS NULL`).Scan(&nulls); err != nil {
		t.Fatal(err)
	}
	if nulls != 2 {
		t.Errorf("isbn 是 NULL 的書有 %d 本，want 2（沒有 ISBN 時要存 NULL，不是空字串）", nulls)
	}
}

func TestByAuthor(t *testing.T) {
	l := newLib(t)
	seedBooks(t, l)
	ctx := t.Context()

	got, err := l.ByAuthor(ctx, "陳小雨")
	want := []Book{
		{ID: 2, Title: "巷口的貓", Author: "陳小雨", ISBN: sql.Null[string]{V: "978-0-00-000001-1", Valid: true}},
		{ID: 1, Title: "雨天的公車", Author: "陳小雨"},
	}
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("陳小雨：\n  %+v, %v\nwant\n  %+v", got, err, want)
	}

	got, err = l.ByAuthor(ctx, "O'Neil")
	if err != nil || len(got) != 1 || got[0].Title != "Tea for Two" {
		t.Errorf("O'Neil：%+v, %v，want Tea for Two（作者名稱有單引號，要用 ? 帶入）", got, err)
	}

	got, err = l.ByAuthor(ctx, "沒有這個人")
	if err != nil || len(got) != 0 {
		t.Errorf("沒有這個人：%+v, %v，want 空的", got, err)
	}

	if n := l.db.Stats().InUse; n != 0 {
		t.Errorf("還有 %d 條連線被佔著：rows 沒有關", n)
	}
}

func TestReturn(t *testing.T) {
	l := newLib(t)
	seedBooks(t, l)
	ctx := t.Context()
	exec(t, l, `INSERT INTO loans (book_id, member, due, returned) VALUES
		(1, '阿明', '2026-11-01', 1),
		(1, '小華', '2026-12-01', 0)`)

	if err := l.Return(ctx, 1); err != nil {
		t.Errorf("Return(1) = %v，want nil", err)
	}
	var open int
	if err := l.db.QueryRowContext(ctx, `SELECT count(*) FROM loans WHERE returned = 0`).Scan(&open); err != nil {
		t.Fatal(err)
	}
	if open != 0 {
		t.Errorf("還有 %d 筆借閱沒還，want 0", open)
	}
	for _, id := range []int64{1, 2, 99} {
		if err := l.Return(ctx, id); !errors.Is(err, ErrNotOnLoan) {
			t.Errorf("Return(%d) = %v，want ErrNotOnLoan", id, err)
		}
	}
}

func TestBorrow(t *testing.T) {
	l := newLib(t)
	seedBooks(t, l)
	ctx := t.Context()
	due := time.Date(2026, 12, 4, 0, 0, 0, 0, time.UTC)

	steps := []struct {
		book   int64
		member string
		want   error
	}{
		{1, "阿明", nil},
		{1, "小華", ErrOnLoan},
		{99, "小華", ErrNotFound},
		{2, "阿明", nil},
		{3, "阿明", nil},
		{4, "阿明", ErrLimit},
	}
	for _, s := range steps {
		if err := l.Borrow(ctx, s.book, s.member, due); !errors.Is(err, s.want) {
			t.Errorf("%s 借 %d：%v，want %v", s.member, s.book, err, s.want)
		}
	}

	var loans int
	var firstDue string
	if err := l.db.QueryRowContext(ctx, `SELECT count(*), min(due) FROM loans`).Scan(&loans, &firstDue); err != nil {
		t.Fatal(err)
	}
	if loans != 3 || firstDue != "2026-12-04" {
		t.Errorf("借閱有 %d 筆、到期日 %q，want 3 筆、2026-12-04（失敗的不能留下資料）", loans, firstDue)
	}

	exec(t, l, `UPDATE loans SET returned = 1 WHERE book_id = 1`)
	if err := l.Borrow(ctx, 4, "阿明", due); err != nil {
		t.Errorf("還了一本之後再借：%v，want nil", err)
	}

	// 10 個人同時借第 5 本書：只能有 1 個人成功。
	var ok, onLoan atomic.Int64
	var wg sync.WaitGroup
	for i := range 10 {
		wg.Go(func() {
			err := l.Borrow(ctx, 5, string(rune('A'+i)), due)
			switch {
			case err == nil:
				ok.Add(1)
			case errors.Is(err, ErrOnLoan):
				onLoan.Add(1)
			default:
				t.Errorf("同時借：%v", err)
			}
		})
	}
	wg.Wait()
	if ok.Load() != 1 || onLoan.Load() != 9 {
		t.Errorf("同時借：成功 %d、已借出 %d，want 1、9", ok.Load(), onLoan.Load())
	}
}

func TestOverdue(t *testing.T) {
	l := newLib(t)
	seedBooks(t, l)
	exec(t, l, `INSERT INTO loans (book_id, member, due, returned) VALUES
		(3, '小華', '2026-11-18', 0),
		(2, '阿明', '2026-11-10', 0),
		(1, '阿明', '2026-11-01', 1),
		(4, '小華', '2026-11-20', 0),
		(5, '阿美', '2026-11-25', 0)`)

	today := time.Date(2026, 11, 20, 21, 30, 0, 0, time.UTC)
	got, err := l.Overdue(t.Context(), today)
	want := []OverdueLoan{{"巷口的貓", "阿明", 10}, {"夜市人生", "小華", 2}}
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("Overdue = %+v, %v\nwant %+v", got, err, want)
	}
}

func TestHasTitle(t *testing.T) {
	l := newLib(t)
	seedBooks(t, l)
	ctx := t.Context()

	tests := []struct {
		author, word string
		want         bool
	}{
		{"陳小雨", "貓", true},
		{"陳小雨", "狗", false},
		{"王大明", "貓", false},
	}
	for _, tt := range tests {
		got, err := l.HasTitle(ctx, tt.author, tt.word)
		if err != nil || got != tt.want {
			t.Errorf("HasTitle(%q, %q) = %v, %v，want %v", tt.author, tt.word, got, err, tt.want)
		}
		if n := l.db.Stats().InUse; n != 0 {
			t.Errorf("HasTitle(%q, %q) 之後還有 %d 條連線被佔著", tt.author, tt.word, n)
		}
	}
}

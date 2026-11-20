//go:build solution

package library

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

// AddBook 用 sql.Null[string] 帶入 isbn：Valid 是 false 時，driver 送出的是 NULL。
// UNIQUE 允許很多個 NULL，但空字串只能有一個。
func (l *Library) AddBook(ctx context.Context, title, author, isbn string) (int64, error) {
	res, err := l.db.ExecContext(ctx,
		`INSERT INTO books (title, author, isbn) VALUES (?, ?, ?)`,
		title, author, sql.Null[string]{V: isbn, Valid: isbn != ""})
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (l *Library) ByAuthor(ctx context.Context, author string) ([]Book, error) {
	rows, err := l.db.QueryContext(ctx,
		`SELECT id, title, author, isbn FROM books WHERE author = ? ORDER BY title, id`, author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []Book{}
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.ISBN); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, rows.Err()
}

// Return 不用先查再改：條件寫在 UPDATE 的 WHERE 裡，再看改到幾列。
func (l *Library) Return(ctx context.Context, bookID int64) error {
	res, err := l.db.ExecContext(ctx,
		`UPDATE loans SET returned = 1 WHERE book_id = ? AND returned = 0`, bookID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotOnLoan
	}
	return nil
}

// Borrow 的檢查和新增都在同一個交易裡。Open 設了 _txlock=immediate，交易一開始就拿到寫入鎖，
// 同時借同一本書的其他交易會排隊，輪到它們時已經看得到前一個人的借閱。
func (l *Library) Borrow(ctx context.Context, bookID int64, member string, due time.Time) error {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists, onLoan bool
	if err := tx.QueryRowContext(ctx,
		`SELECT EXISTS (SELECT 1 FROM books WHERE id = ?)`, bookID).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return ErrNotFound
	}
	if err := tx.QueryRowContext(ctx,
		`SELECT EXISTS (SELECT 1 FROM loans WHERE book_id = ? AND returned = 0)`, bookID).Scan(&onLoan); err != nil {
		return err
	}
	if onLoan {
		return ErrOnLoan
	}
	var holding int
	if err := tx.QueryRowContext(ctx,
		`SELECT count(*) FROM loans WHERE member = ? AND returned = 0`, member).Scan(&holding); err != nil {
		return err
	}
	if holding >= MaxLoans {
		return ErrLimit
	}

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO loans (book_id, member, due) VALUES (?, ?, ?)`,
		bookID, member, due.Format(time.DateOnly)); err != nil {
		return err
	}
	return tx.Commit()
}

// Overdue 讓資料庫做篩選和排序：日期存成 2006-01-02，字串的大小順序就是日期的先後。
func (l *Library) Overdue(ctx context.Context, today time.Time) ([]OverdueLoan, error) {
	day := today.Format(time.DateOnly)
	rows, err := l.db.QueryContext(ctx, `
		SELECT b.title, l.member, l.due
		FROM loans AS l JOIN books AS b ON b.id = l.book_id
		WHERE l.returned = 0 AND l.due < ?
		ORDER BY l.due, b.title`, day)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	start, err := time.Parse(time.DateOnly, day)
	if err != nil {
		return nil, err
	}
	var out []OverdueLoan
	for rows.Next() {
		var o OverdueLoan
		var due string
		if err := rows.Scan(&o.Title, &o.Member, &due); err != nil {
			return nil, err
		}
		d, err := time.Parse(time.DateOnly, due)
		if err != nil {
			return nil, err
		}
		o.DaysLate = int(start.Sub(d).Hours() / 24)
		out = append(out, o)
	}
	return out, rows.Err()
}

// HasTitle 用 defer rows.Close()：從迴圈中間 return 時也會關掉，連線才會還給連線池。
func (l *Library) HasTitle(ctx context.Context, author, word string) (bool, error) {
	rows, err := l.db.QueryContext(ctx, `SELECT title FROM books WHERE author = ?`, author)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return false, err
		}
		if strings.Contains(title, word) {
			return true, nil
		}
	}
	return false, rows.Err()
}

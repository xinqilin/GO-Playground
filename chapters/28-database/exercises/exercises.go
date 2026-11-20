//go:build !solution

package library

import (
	"context"
	"strings"
	"time"
)

// AddBook 新增一本書，回傳它的 id。isbn 是空字串時，資料庫裡存 NULL（不是空字串）。
// 同一個 ISBN 第二次新增時，資料庫會拒絕，回傳那個錯誤就好。
func (l *Library) AddBook(ctx context.Context, title, author, isbn string) (int64, error) {
	// TODO: 實作
	return 0, nil
}

// ByAuthor 回傳 author 寫的所有書，依書名排序，書名一樣時依 id。沒有的話回傳空的 slice 和 nil。
func (l *Library) ByAuthor(ctx context.Context, author string) ([]Book, error) {
	// TODO: 實作
	return nil, nil
}

// Return 還書：把 bookID 還沒還的那筆借閱標成已還。這本書沒有借出去（或不存在）時回傳 ErrNotOnLoan。
func (l *Library) Return(ctx context.Context, bookID int64) error {
	// TODO: 實作
	return nil
}

// Borrow 借書，在一個交易裡依序檢查：
//   - 書不存在：ErrNotFound
//   - 書已經借出去、還沒還：ErrOnLoan
//   - member 手上還沒還的書已經有 MaxLoans 本：ErrLimit
//
// 都通過才新增一筆借閱，due 存成 2006-01-02 的格式。失敗時資料庫不能有任何改變。
// 很多人同時借同一本書時，只能有一個人成功。
func (l *Library) Borrow(ctx context.Context, bookID int64, member string, due time.Time) error {
	// TODO: 實作
	return nil
}

// Overdue 回傳到 today 為止過期還沒還的借閱：到期日在 today 之前（不含當天）。
// 依到期日排序，同一天的依書名。DaysLate 是 today 的日期減掉到期日的天數。
func (l *Library) Overdue(ctx context.Context, today time.Time) ([]OverdueLoan, error) {
	// TODO: 實作
	return nil, nil
}

// HasTitle 回傳 author 有沒有一本書的書名含有 word。
//
// 找錯：呼叫幾次之後，其他的查詢都卡住不動了。
func (l *Library) HasTitle(ctx context.Context, author, word string) (bool, error) {
	rows, err := l.db.QueryContext(ctx, `SELECT title FROM books WHERE author = ?`, author)
	if err != nil {
		return false, err
	}
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return false, err
		}
		if strings.Contains(title, word) {
			return true, nil
		}
	}
	rows.Close()
	return false, rows.Err()
}

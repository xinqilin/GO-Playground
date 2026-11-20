// Package library 是第 28 章「資料庫」的練習題：社區圖書室的借還書系統。
//
// 這一章是獨立的 module，所以先進到它的目錄：
//
//	cd chapters/28-database
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/
//
// 只想跑某一題時加上 -run：
//
//	go test ./exercises/ -run TestBorrow
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/
package library

import (
	"context"
	"database/sql"
	"errors"

	_ "modernc.org/sqlite" // 註冊 "sqlite" driver
)

// 圖書室的錯誤，測試用 errors.Is 判斷。
var (
	ErrNotFound  = errors.New("找不到這本書")
	ErrOnLoan    = errors.New("這本書已經借出去了")
	ErrNotOnLoan = errors.New("這本書沒有借出去")
	ErrLimit     = errors.New("借的書已經到上限了")
)

// MaxLoans 是一個人同時最多可以借幾本書。
const MaxLoans = 3

// Book 是一本書。ISBN 可以是 NULL：社區捐的舊書常常沒有。
type Book struct {
	ID     int64
	Title  string
	Author string
	ISBN   sql.Null[string]
}

// OverdueLoan 是一筆過了到期日還沒還的借閱。
type OverdueLoan struct {
	Title    string
	Member   string
	DaysLate int
}

// Library 是圖書室的資料庫。
type Library struct {
	db *sql.DB
}

const schema = `
CREATE TABLE IF NOT EXISTS books (
    id     INTEGER PRIMARY KEY,
    title  TEXT NOT NULL,
    author TEXT NOT NULL,
    isbn   TEXT UNIQUE
);
CREATE TABLE IF NOT EXISTS loans (
    id       INTEGER PRIMARY KEY,
    book_id  INTEGER NOT NULL REFERENCES books (id),
    member   TEXT    NOT NULL,
    due      TEXT    NOT NULL,           -- 到期日，格式是 2006-01-02
    returned INTEGER NOT NULL DEFAULT 0  -- 0：還沒還，1：已經還了
);`

// Open 開啟 path 的資料庫並建立資料表。練習直接用它，不用修改。
func Open(ctx context.Context, path string) (*Library, error) {
	db, err := sql.Open("sqlite", path+"?_foreign_keys=1&_busy_timeout=5000&_journal_mode=WAL&_txlock=immediate")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(4)
	if _, err := db.ExecContext(ctx, schema); err != nil {
		db.Close()
		return nil, err
	}
	return &Library{db: db}, nil
}

// Close 關掉資料庫。
func (l *Library) Close() error {
	return l.db.Close()
}

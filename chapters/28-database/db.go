package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite" // 只為了它的 init：向 database/sql 註冊名叫 "sqlite" 的 driver
)

// openDB 開啟 path 的 SQLite 資料庫（檔案不存在就建立），設定連線池、確認連得上，再跑 migration。
func openDB(ctx context.Context, path string) (*sql.DB, error) {
	// ? 後面底線開頭的參數由 driver 處理：開啟外鍵檢查、被鎖住時最多等 5 秒、
	// 使用 WAL 模式（讀寫可以同時進行）、交易一開始就拿寫入鎖（第 6 節）。
	dsn := path + "?_foreign_keys=1&_busy_timeout=5000&_journal_mode=WAL&_txlock=immediate"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(8)
	db.SetMaxIdleConns(8)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("連不上資料庫：%w", err)
	}
	if err := migrate(ctx, db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

//go:embed migrations/*.sql
var migrations embed.FS

// migrate 依照檔名的順序，執行還沒跑過的 migration。
// 已經跑到第幾個，記在資料庫檔案裡的 PRAGMA user_version。
func migrate(ctx context.Context, db *sql.DB) error {
	names, err := fs.Glob(migrations, "migrations/*.sql")
	if err != nil {
		return err
	}
	var done int
	if err := db.QueryRowContext(ctx, "PRAGMA user_version").Scan(&done); err != nil {
		return err
	}
	if done > len(names) {
		return fmt.Errorf("資料庫的版本 %d 比程式認得的 %d 新，請更新程式", done, len(names))
	}
	for i := done; i < len(names); i++ {
		if err := applyMigration(ctx, db, names[i], i+1); err != nil {
			return fmt.Errorf("%s：%w", names[i], err)
		}
	}
	return nil
}

// applyMigration 在一個交易裡執行一個 migration 檔，並把 user_version 改成 version。
// 中途失敗時兩者一起復原，不會留下改到一半的資料表。
func applyMigration(ctx context.Context, db *sql.DB, name string, version int) error {
	stmts, err := migrations.ReadFile(name)
	if err != nil {
		return err
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // Commit 成功之後再 Rollback 不會有作用

	if _, err := tx.ExecContext(ctx, string(stmts)); err != nil {
		return err
	}
	// PRAGMA 不能用 ? 帶參數；version 是程式自己算出來的整數，不是外面來的輸入，所以可以直接放進字串。
	if _, err := tx.ExecContext(ctx, fmt.Sprintf("PRAGMA user_version = %d", version)); err != nil {
		return err
	}
	return tx.Commit()
}

// newDemoStore 在暫存目錄建立一個新的資料庫；seed 為 true 時先放好三道菜。
// 回傳的函式會關掉資料庫並刪掉暫存目錄。示範用，所以準備失敗就直接 panic。
func newDemoStore(seed bool) (*store, func()) {
	ctx := context.Background()
	dir, err := os.MkdirTemp("", "bento-*")
	if err != nil {
		panic(err)
	}
	db, err := openDB(ctx, filepath.Join(dir, "bento.db"))
	if err != nil {
		os.RemoveAll(dir)
		panic(err)
	}
	s := &store{db: db}
	if seed {
		for _, d := range []Dish{{Name: "滷肉飯", Price: 50, Stock: 20}, {Name: "排骨飯", Price: 90, Stock: 5}, {Name: "雞腿飯", Price: 100, Stock: 0}} {
			if _, err := s.addDish(ctx, d.Name, d.Price, d.Stock); err != nil {
				panic(err)
			}
		}
	}
	return s, func() {
		db.Close()
		os.RemoveAll(dir)
	}
}

// openDemo 開一個新的資料庫：migration 建好了兩個資料表。
func openDemo() {
	ctx := context.Background()
	dir, err := os.MkdirTemp("", "bento-*")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.RemoveAll(dir)

	db, err := openDB(ctx, filepath.Join(dir, "bento.db"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	var version int
	if err := db.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version); err != nil {
		fmt.Println(err)
		return
	}
	tables, err := (&store{db: db}).names(ctx, `SELECT name FROM sqlite_schema WHERE type = 'table' ORDER BY name`)
	fmt.Printf("migration 版本：%d\n", version)
	fmt.Printf("資料表：%v %v\n", tables, err)
	fmt.Printf("連線池上限：%d\n", db.Stats().MaxOpenConnections)
}

// migrateDemo 同一個資料庫檔案開兩次，migration 不會重跑；資料庫比程式新的時候拒絕開啟。
func migrateDemo() {
	ctx := context.Background()
	dir, err := os.MkdirTemp("", "bento-*")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "bento.db")

	for i := 1; i <= 2; i++ {
		db, err := openDB(ctx, path)
		if err != nil {
			fmt.Println(err)
			return
		}
		var version int
		db.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version)
		fmt.Printf("第 %d 次開啟：版本 %d\n", i, version)
		if i == 2 {
			// 假裝有一個更新版的程式，已經把這個資料庫升到第 3 版。
			db.ExecContext(ctx, "PRAGMA user_version = 3")
		}
		db.Close()
	}

	_, err = openDB(ctx, path)
	fmt.Printf("舊程式開新資料庫：%v\n", err)
}

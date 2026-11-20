package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// 店家的錯誤。呼叫的一方用 errors.Is 判斷，不需要知道底下是哪一種資料庫。
var (
	ErrNotFound  = errors.New("找不到")
	ErrDuplicate = errors.New("已經有了")
	ErrSoldOut   = errors.New("賣完了")
)

// Dish 是菜單上的一道菜。
type Dish struct {
	ID    int64
	Name  string
	Price int
	Stock int
}

// Order 是一筆訂單。Note 可以是 NULL：客人沒寫備註。
type Order struct {
	ID     int64
	DishID int64
	Qty    int
	Total  int
	Note   sql.Null[string]
}

// store 把所有的 SQL 集中在一個地方。*sql.DB 是連線池，可以被很多 goroutine 同時使用。
type store struct {
	db *sql.DB
}

// isUnique 判斷 err 是不是違反了 UNIQUE。錯誤的型別是 driver 自己的，只在 store 這一層出現。
func isUnique(err error) bool {
	e, ok := errors.AsType[*sqlite.Error](err)
	return ok && e.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE
}

// addDish 新增一道菜，回傳資料庫給它的 id。
func (s *store) addDish(ctx context.Context, name string, price, stock int) (int64, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO dishes (name, price, stock) VALUES (?, ?, ?)`, name, price, stock)
	if isUnique(err) {
		return 0, fmt.Errorf("%s %w", name, ErrDuplicate)
	}
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// restock 把 name 的庫存加上 n。
func (s *store) restock(ctx context.Context, name string, n int) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE dishes SET stock = stock + ? WHERE name = ?`, n, name)
	if err != nil {
		return err
	}
	changed, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if changed == 0 {
		return fmt.Errorf("%s：%w", name, ErrNotFound)
	}
	return nil
}

// findDish 用名稱找一道菜。
func (s *store) findDish(ctx context.Context, name string) (Dish, error) {
	var d Dish
	err := s.db.QueryRowContext(ctx,
		`SELECT id, name, price, stock FROM dishes WHERE name = ?`, name,
	).Scan(&d.ID, &d.Name, &d.Price, &d.Stock)
	if errors.Is(err, sql.ErrNoRows) {
		return Dish{}, fmt.Errorf("%s：%w", name, ErrNotFound)
	}
	return d, err
}

// listDishes 列出價格不超過 maxPrice 的菜，便宜的在前面。
func (s *store) listDishes(ctx context.Context, maxPrice int) ([]Dish, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, price, stock FROM dishes WHERE price <= ? ORDER BY price, id`, maxPrice)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dishes []Dish
	for rows.Next() {
		var d Dish
		if err := rows.Scan(&d.ID, &d.Name, &d.Price, &d.Stock); err != nil {
			return nil, err
		}
		dishes = append(dishes, d)
	}
	return dishes, rows.Err()
}

// names 執行查詢，回傳每一列第一個欄位的字串。
func (s *store) names(ctx context.Context, query string, args ...any) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		out = append(out, name)
	}
	return out, rows.Err()
}

// placeOrder 下單：檢查庫存、扣庫存、新增訂單。三件事在同一個交易裡，要嘛全部生效，要嘛全部沒發生。
func (s *store) placeOrder(ctx context.Context, dish string, qty int, note string) (Order, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return Order{}, err
	}
	defer tx.Rollback() // 中途 return 時復原；Commit 之後呼叫則什麼都不做

	var id int64
	var price, stock int
	err = tx.QueryRowContext(ctx,
		`SELECT id, price, stock FROM dishes WHERE name = ?`, dish).Scan(&id, &price, &stock)
	if errors.Is(err, sql.ErrNoRows) {
		return Order{}, fmt.Errorf("%s：%w", dish, ErrNotFound)
	}
	if err != nil {
		return Order{}, err
	}
	if stock < qty {
		return Order{}, fmt.Errorf("%s %w，剩 %d 份", dish, ErrSoldOut, stock)
	}

	if _, err := tx.ExecContext(ctx,
		`UPDATE dishes SET stock = stock - ? WHERE id = ?`, qty, id); err != nil {
		return Order{}, err
	}
	o := Order{DishID: id, Qty: qty, Total: price * qty, Note: sql.Null[string]{V: note, Valid: note != ""}}
	res, err := tx.ExecContext(ctx,
		`INSERT INTO orders (dish_id, qty, total, note) VALUES (?, ?, ?, ?)`, o.DishID, o.Qty, o.Total, o.Note)
	if err != nil {
		return Order{}, err
	}
	if o.ID, err = res.LastInsertId(); err != nil {
		return Order{}, err
	}
	return o, tx.Commit()
}

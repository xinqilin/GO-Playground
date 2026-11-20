package main

import (
	"database/sql"
	"errors"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func Example_openDemo() {
	openDemo()
	// Output:
	// migration 版本：2
	// 資料表：[dishes orders] <nil>
	// 連線池上限：8
}

func Example_execDemo() {
	execDemo()
	// Output:
	// 新增 滷肉飯 id = 1 <nil>
	// 新增 排骨飯 id = 2 <nil>
	// 新增 雞腿飯 id = 3 <nil>
	// 排骨飯 已經有了 true
	// constraint failed: CHECK constraint failed: price > 0 (275)
	// 補雞腿飯：<nil>
	// 補鰻魚飯：鰻魚飯：找不到
}

func Example_queryDemo() {
	queryDemo()
	// Output:
	// 1 滷肉飯 50 元，剩 20 份
	// 2 排骨飯 90 元，剩 5 份
	// {ID:3 Name:雞腿飯 Price:100 Stock:0} <nil>
	// 鰻魚飯：找不到 true
}

func Example_injectionDemo() {
	injectionDemo()
	// Output:
	// 拼字串 "排骨飯" → [排骨飯] <nil>
	// 用 ?   "排骨飯" → [排骨飯] <nil>
	// 拼字串 "x' OR '1'='1" → [排骨飯 滷肉飯 雞腿飯] <nil>
	// 用 ?   "x' OR '1'='1" → [] <nil>
	// 拼字串 "阿嬤's 滷肉" → [] SQL logic error: near "s": syntax error (1)
	// 用 ?   "阿嬤's 滷肉" → [] <nil>
	// 排序 stock → [滷肉飯 排骨飯 雞腿飯] <nil>
	// 排序 price; DROP TABLE dishes → [] 不能用 "price; DROP TABLE dishes" 排序
}

func Example_nullDemo() {
	nullDemo()
	// Output:
	// Scan 進 string：sql: Scan error on column index 0, name "note": converting NULL to string is unsupported
	// 訂單 1 沒有備註
	// 訂單 2 的備註：不要辣
	// 用 COALESCE：[（無） 不要辣] <nil>
}

func Example_txDemo() {
	txDemo()
	// Output:
	// {ID:1 DishID:2 Qty:2 Total:180 Note:{V:不要辣 Valid:true}} <nil>
	// 排骨飯 賣完了，剩 3 份 true
	// 鰻魚飯：找不到
	// 排骨飯剩 3 份，共 1 筆訂單
}

func Example_poolDemo() {
	poolDemo()
	// Output:
	// 下一個查詢：context deadline exceeded
	// 使用中 1 條，等過 1 次
	// 關掉 rows 之後：<nil>
	// 數到十億：context deadline exceeded true
}

func Example_migrateDemo() {
	migrateDemo()
	// Output:
	// 第 1 次開啟：版本 2
	// 第 2 次開啟：版本 2
	// 舊程式開新資料庫：資料庫的版本 3 比程式認得的 2 新，請更新程式
}

// newTestStore 給每個測試一個全新的資料庫檔案，測試結束時自動關掉、刪掉。
func newTestStore(t *testing.T) *store {
	t.Helper()
	db, err := openDB(t.Context(), filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return &store{db: db}
}

// TestPlaceOrderConcurrent：20 個人同時搶 5 份排骨飯，剛好 5 個人成功，庫存不會變成負的。
func TestPlaceOrderConcurrent(t *testing.T) {
	s := newTestStore(t)
	ctx := t.Context()
	if _, err := s.addDish(ctx, "排骨飯", 90, 5); err != nil {
		t.Fatal(err)
	}

	var ok, soldOut atomic.Int64
	var wg sync.WaitGroup
	for range 20 {
		wg.Go(func() {
			_, err := s.placeOrder(ctx, "排骨飯", 1, "")
			switch {
			case err == nil:
				ok.Add(1)
			case errors.Is(err, ErrSoldOut):
				soldOut.Add(1)
			default:
				t.Errorf("非預期的錯誤：%v", err)
			}
		})
	}
	wg.Wait()

	d, err := s.findDish(ctx, "排骨飯")
	if err != nil {
		t.Fatal(err)
	}
	var orders int
	if err := s.db.QueryRowContext(ctx, `SELECT count(*) FROM orders`).Scan(&orders); err != nil {
		t.Fatal(err)
	}
	if ok.Load() != 5 || soldOut.Load() != 15 || d.Stock != 0 || orders != 5 {
		t.Errorf("成功 %d、賣完 %d、剩 %d 份、%d 筆訂單；want 5、15、0、5",
			ok.Load(), soldOut.Load(), d.Stock, orders)
	}
}

// TestMemoryPerConnection：":memory:" 的資料庫只屬於一條連線，
// 連線池裡的另一條連線看到的是另一個空的資料庫。
func TestMemoryPerConnection(t *testing.T) {
	ctx := t.Context()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	c1, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c1.Close()
	c2, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c2.Close()

	if _, err := c1.ExecContext(ctx, `CREATE TABLE dishes (name TEXT)`); err != nil {
		t.Fatal(err)
	}
	_, err = c2.ExecContext(ctx, `INSERT INTO dishes VALUES ('排骨飯')`)
	if err == nil || !strings.Contains(err.Error(), "no such table") {
		t.Errorf("第二條連線：err = %v，want no such table", err)
	}
}

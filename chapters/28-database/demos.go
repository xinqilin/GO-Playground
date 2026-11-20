package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// execDemo 新增、修改資料，以及資料庫擋下來的錯誤。
func execDemo() {
	ctx := context.Background()
	s, done := newDemoStore(false)
	defer done()

	for _, d := range []Dish{{Name: "滷肉飯", Price: 50, Stock: 20}, {Name: "排骨飯", Price: 90, Stock: 5}, {Name: "雞腿飯", Price: 100, Stock: 0}} {
		id, err := s.addDish(ctx, d.Name, d.Price, d.Stock)
		fmt.Println("新增", d.Name, "id =", id, err)
	}

	_, err := s.addDish(ctx, "排骨飯", 95, 3)
	fmt.Println(err, errors.Is(err, ErrDuplicate))
	_, err = s.addDish(ctx, "空氣便當", 0, 1)
	fmt.Println(err)

	fmt.Printf("補雞腿飯：%v\n", s.restock(ctx, "雞腿飯", 10))
	fmt.Printf("補鰻魚飯：%v\n", s.restock(ctx, "鰻魚飯", 10))
}

// queryDemo 讀很多列、讀一列，以及一列都沒有。
func queryDemo() {
	ctx := context.Background()
	s, done := newDemoStore(true)
	defer done()

	dishes, err := s.listDishes(ctx, 90)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, d := range dishes {
		fmt.Printf("%d %s %d 元，剩 %d 份\n", d.ID, d.Name, d.Price, d.Stock)
	}

	d, err := s.findDish(ctx, "雞腿飯")
	fmt.Printf("%+v %v\n", d, err)
	_, err = s.findDish(ctx, "鰻魚飯")
	fmt.Println(err, errors.Is(err, ErrNotFound))
}

// searchUnsafe 用字串相加組出 SQL。不要這樣寫：name 的內容會變成 SQL 的一部分。
func (s *store) searchUnsafe(ctx context.Context, name string) ([]string, error) {
	return s.names(ctx, "SELECT name FROM dishes WHERE name = '"+name+"'")
}

// search 用 ? 帶入 name：它永遠只是一個值，不會變成 SQL。
func (s *store) search(ctx context.Context, name string) ([]string, error) {
	return s.names(ctx, "SELECT name FROM dishes WHERE name = ?", name)
}

// sortColumns 是允許排序的方式。欄位名稱不能用 ? 帶入，只能從固定的清單裡挑。
var sortColumns = map[string]string{
	"price": "price, id",
	"name":  "name",
	"stock": "stock DESC, id",
}

// dishNamesBy 依照 sort 指定的方式列出菜名；sort 通常來自網址的查詢字串。
func (s *store) dishNamesBy(ctx context.Context, sort string) ([]string, error) {
	order, ok := sortColumns[sort]
	if !ok {
		return nil, fmt.Errorf("不能用 %q 排序", sort)
	}
	return s.names(ctx, "SELECT name FROM dishes ORDER BY "+order)
}

func injectionDemo() {
	ctx := context.Background()
	s, done := newDemoStore(true)
	defer done()

	for _, input := range []string{"排骨飯", "x' OR '1'='1", "阿嬤's 滷肉"} {
		got, err := s.searchUnsafe(ctx, input)
		fmt.Printf("拼字串 %q → %v %v\n", input, got, err)
		got, err = s.search(ctx, input)
		fmt.Printf("用 ?   %q → %v %v\n", input, got, err)
	}

	for _, sort := range []string{"stock", "price; DROP TABLE dishes"} {
		got, err := s.dishNamesBy(ctx, sort)
		fmt.Println("排序", sort, "→", got, err)
	}
}

// nullDemo 讀出可能是 NULL 的欄位。
func nullDemo() {
	ctx := context.Background()
	s, done := newDemoStore(true)
	defer done()

	_, err := s.db.ExecContext(ctx, `INSERT INTO orders (dish_id, qty, total, note)
		VALUES (1, 1, 50, NULL), (2, 2, 180, '不要辣')`)
	if err != nil {
		fmt.Println(err)
		return
	}

	var note string
	err = s.db.QueryRowContext(ctx, `SELECT note FROM orders WHERE id = 1`).Scan(&note)
	fmt.Printf("Scan 進 string：%v\n", err)

	for id := 1; id <= 2; id++ {
		var o Order
		err := s.db.QueryRowContext(ctx,
			`SELECT id, note FROM orders WHERE id = ?`, id).Scan(&o.ID, &o.Note)
		if err != nil {
			fmt.Println(err)
			return
		}
		if o.Note.Valid {
			fmt.Printf("訂單 %d 的備註：%s\n", o.ID, o.Note.V)
		} else {
			fmt.Printf("訂單 %d 沒有備註\n", o.ID)
		}
	}

	notes, err := s.names(ctx, `SELECT COALESCE(note, '（無）') FROM orders ORDER BY id`)
	fmt.Printf("用 COALESCE：%v %v\n", notes, err)
}

// txDemo 下單成功、賣完、找不到：失敗的時候，庫存和訂單都不會被改到一半。
func txDemo() {
	ctx := context.Background()
	s, done := newDemoStore(true)
	defer done()

	o, err := s.placeOrder(ctx, "排骨飯", 2, "不要辣")
	fmt.Printf("%+v %v\n", o, err)
	_, err = s.placeOrder(ctx, "排骨飯", 5, "")
	fmt.Println(err, errors.Is(err, ErrSoldOut))
	_, err = s.placeOrder(ctx, "鰻魚飯", 1, "")
	fmt.Println(err)

	d, _ := s.findDish(ctx, "排骨飯")
	var orders int
	s.db.QueryRowContext(ctx, `SELECT count(*) FROM orders`).Scan(&orders)
	fmt.Printf("排骨飯剩 %d 份，共 %d 筆訂單\n", d.Stock, orders)
}

// poolDemo 忘了關 rows 會佔住連線；context 的期限讓等待和慢查詢都能停下來。
func poolDemo() {
	ctx := context.Background()
	s, done := newDemoStore(true)
	defer done()
	s.db.SetMaxOpenConns(1) // 示範用：連線池只有 1 條連線

	rows, err := s.db.QueryContext(ctx, `SELECT name FROM dishes`)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 忘了 rows.Close()：唯一的連線一直被佔著。

	waitCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	_, err = s.findDish(waitCtx, "排骨飯")
	st := s.db.Stats()
	fmt.Printf("下一個查詢：%v\n", err)
	fmt.Printf("使用中 %d 條，等過 %d 次\n", st.InUse, st.WaitCount)

	rows.Close()
	_, err = s.findDish(ctx, "排骨飯")
	fmt.Printf("關掉 rows 之後：%v\n", err)

	slowCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	var n int
	err = s.db.QueryRowContext(slowCtx, `WITH RECURSIVE c(x) AS (
		SELECT 1 UNION ALL SELECT x + 1 FROM c WHERE x < 1000000000
	) SELECT count(*) FROM c`).Scan(&n)
	fmt.Printf("數到十億：%v %v\n", err, errors.Is(err, context.DeadlineExceeded))
}

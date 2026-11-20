-- 訂單：dish_id 一定要對到一道存在的菜；備註可以是 NULL。
CREATE TABLE orders (
    id      INTEGER PRIMARY KEY,
    dish_id INTEGER NOT NULL REFERENCES dishes (id),
    qty     INTEGER NOT NULL CHECK (qty > 0),
    total   INTEGER NOT NULL,
    note    TEXT
);

CREATE INDEX orders_dish_id ON orders (dish_id);

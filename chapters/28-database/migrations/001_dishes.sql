-- 菜單：每道菜的名稱不能重複，價格要大於 0，庫存不能是負的。
CREATE TABLE dishes (
    id    INTEGER PRIMARY KEY,
    name  TEXT    NOT NULL UNIQUE,
    price INTEGER NOT NULL CHECK (price > 0),
    stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0)
);

-- Основная таблица заказов
CREATE TABLE IF NOT EXISTS orders (
    order_uid          VARCHAR(255) PRIMARY KEY,
    track_number       VARCHAR(255) NOT NULL,
    entry              VARCHAR(255) NOT NULL,
    locale             VARCHAR(10) NOT NULL,
    internal_signature VARCHAR(255) DEFAULT '',
    customer_id        VARCHAR(255) NOT NULL,
    delivery_service   VARCHAR(255) NOT NULL,
    shardkey           VARCHAR(255) NOT NULL,
    sm_id              INTEGER NOT NULL,
    date_created       TIMESTAMPTZ NOT NULL,
    oof_shard          VARCHAR(255) NOT NULL
);


-- Таблица доставки
CREATE TABLE IF NOT EXISTS deliveries (
    order_uid VARCHAR(255) PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
    name      VARCHAR(255) NOT NULL,
    phone     VARCHAR(255) NOT NULL,
    zip       VARCHAR(255) NOT NULL,
    city      VARCHAR(255) NOT NULL,
    address   VARCHAR(255) NOT NULL,
    region    VARCHAR(255) NOT NULL,
    email     VARCHAR(255) NOT NULL
);

-- Таблица платежей
CREATE TABLE IF NOT EXISTS payments (
    order_uid     VARCHAR(255) PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
    transaction  VARCHAR(255) NOT NULL,
    request_id    VARCHAR(255) DEFAULT '',
    currency     VARCHAR(3) NOT NULL,
    provider     VARCHAR(255) NOT NULL,
    amount       INTEGER NOT NULL CHECK (amount >= 0),
    payment_dt    BIGINT NOT NULL,
    bank         VARCHAR(255) NOT NULL,
    delivery_cost INTEGER NOT NULL CHECK (delivery_cost >= 0),
    goods_total   INTEGER NOT NULL CHECK (goods_total >= 0),
    custom_fee    INTEGER DEFAULT 0 CHECK (custom_fee >= 0)
);

-- Таблица товаров
CREATE TABLE IF NOT EXISTS items (
    order_uid    VARCHAR(255) NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id      BIGINT NOT NULL,
    track_number VARCHAR(255) NOT NULL,
    price        INTEGER NOT NULL CHECK (price >= 0),
    rid          VARCHAR(255) NOT NULL,
    name         VARCHAR(255) NOT NULL,
    sale         INTEGER NOT NULL CHECK (sale >= 0 AND sale <= 100),
    size         VARCHAR(255) DEFAULT '',
    total_price  INTEGER NOT NULL CHECK (total_price >= 0),
    nm_id        BIGINT NOT NULL,
    brand        VARCHAR(255) NOT NULL,
    status       INTEGER NOT NULL,
    PRIMARY KEY (order_uid, chrt_id)
);

-- Индексы для таблицы заказов
CREATE INDEX IF NOT EXISTS idx_orders_track_number ON orders(track_number);
CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders(customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_date_created ON orders(date_created);

-- Индексы для таблицы товаров
CREATE INDEX IF NOT EXISTS idx_items_order_uid ON items(order_uid);
CREATE INDEX IF NOT EXISTS idx_items_chrt_id ON items(chrt_id);
CREATE INDEX IF NOT EXISTS idx_items_nm_id ON items(nm_id);

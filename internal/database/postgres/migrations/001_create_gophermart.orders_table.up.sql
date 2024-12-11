CREATE TABLE IF NOT EXISTS gophermart.orders (
    order_id TEXT NOT NULL UNIQUE,
    order_sum DECIMAL,
    processed_at TIMESTAMPTZ,
    user_id INTEGER,
    PRIMARY KEY ( order_id ),
    FOREIGN KEY ( user_id ) REFERENCES gophermart.users ( user_id )
);
CREATE TABLE IF NOT EXISTS orders (
    order_id TEXT NOT NULL UNIQUE,
    sum DECIMAL,
    processed_at TIMESTAMPTZ,
    user_id INTEGER,
    PRIMARY KEY ( order_id ),
    FOREIGN KEY ( user_id ) REFERENCES users ( user_id )
)
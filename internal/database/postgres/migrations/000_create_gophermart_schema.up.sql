;CREATE SCHEMA IF NOT EXISTS gophermart;
CREATE TABLE IF NOT EXISTS gophermart.system (
    var_id TEXT NOT NULL UNIQUE,
    var_value TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS gophermart.users (
    user_id SERIAL,
    user_name TEXT NOT NULL UNIQUE,
    password_hash BYTEA NOT NULL,
    password_salt BYTEA NOT NULL,
    create_user_timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ( user_id ) 
);
CREATE TABLE IF NOT EXISTS gophermart.user_sessions (
    session_id SERIAL NOT NULL,
    user_id INTEGER,
    token TEXT NOT NULL,
    PRIMARY KEY ( session_id ),
    FOREIGN KEY ( user_id ) REFERENCES gophermart.users ( user_id ) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS gophermart.orders (
    order_id TEXT NOT NULL UNIQUE,
    order_sum DECIMAL,
    processed_at TIMESTAMPTZ,
    user_id INTEGER,
    PRIMARY KEY ( order_id ),
    FOREIGN KEY ( user_id ) REFERENCES gophermart.users ( user_id )
);
CREATE TABLE IF NOT EXISTS gophermart.balances (
    user_id INTEGER,
	current_balance DECIMAL,
    withdrawn DECIMAL,
    FOREIGN KEY ( user_id ) REFERENCES gophermart.users ( user_id )
);
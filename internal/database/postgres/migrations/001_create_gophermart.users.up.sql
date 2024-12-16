CREATE TABLE IF NOT EXISTS gophermart.users (
    user_id SERIAL,
    user_name TEXT NOT NULL UNIQUE,
    pasword_hash BYTEA NOT NULL,
    pasword_salt BYTEA NOT NULL,
    create_user_timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ( user_id ) 
);

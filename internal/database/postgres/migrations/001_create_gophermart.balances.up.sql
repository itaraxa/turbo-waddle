CREATE TABLE IF NOT EXISTS gophermart.balances (
    user_id INTEGER,
	current_balance DECIMAL,
    withdrawn DECIMAL,
    FOREIGN KEY ( user_id ) REFERENCES gophermart.users ( user_id )
);
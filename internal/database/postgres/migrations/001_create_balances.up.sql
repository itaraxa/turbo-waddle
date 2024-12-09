CREATE TABLE IF NOT EXISTS balances (
    current DECIMAL;
    withdrawn DECIMAL;
    REFERENCES users ( user_id )
);
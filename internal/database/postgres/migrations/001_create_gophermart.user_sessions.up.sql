CREATE TABLE IF NOT EXISTS gophermart.user_sessions (
    session_id SERIAL NOT NULL,
    user_id INTEGER,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    start_timestamp TIMESTAMPTZ,
    end_timestamp TIMESTAMPTZ,
    PRIMARY KEY ( session_id ),
    FOREIGN KEY ( user_id ) REFERENCES gophermart.users ( user_id ) ON DELETE CASCADE
);
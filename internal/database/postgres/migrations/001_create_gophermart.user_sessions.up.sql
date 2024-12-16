CREATE TABLE IF NOT EXISTS gophermart.user_sessions (
    session_id SERIAL NOT NULL,
    user_id INTEGER,
    token TEXT NOT NULL,
    PRIMARY KEY ( session_id ),
    FOREIGN KEY ( user_id ) REFERENCES gophermart.users ( user_id ) ON DELETE CASCADE
);
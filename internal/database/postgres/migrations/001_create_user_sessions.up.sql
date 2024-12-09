CREATE TABLE IF NOT EXISTS user_sessions (
    session_id SERIAL NOT NULL,
    user_id INTEGER,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    start_timestamp TIMESTAMPTZ,
    end_timestamp TIMESTAMPTZ,
    PRIMARY KEY ( session_id ),
    FOREIGN KEY ( user_id ) REFERENCES users ( user_id ) ON DELETE CASCADE
);

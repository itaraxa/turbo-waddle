package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

/*
PostgresRepository is the struct for wrapping PostgreSQL storage
*/
type PostgresRepository struct {
	db *sql.DB
	mu sync.Mutex
}

/*
NewPostgresRepository creates instance of PostgresRepository

Args:

	ctx context.Context
	databaseURL: string for connection to databse, example: "postgres://username:password@localhost:5432/database_name"

Returns:

	dbStorager
	error
*/
func NewPostgresRepository(ctx context.Context, databaseURL string) (*PostgresRepository, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	ctxWithTimeout, cancelWithTimeout := context.WithTimeout(ctx, 5*time.Second)
	defer cancelWithTimeout()

	err = prepareTablesContext(ctxWithTimeout, db)
	if err != nil {
		return nil, fmt.Errorf("cannot create tables in database storage: %w", err)
	}

	return &PostgresRepository{db: db}, nil
}

/*
PingContext check connection to db

Args:

	ctx context.Context

Returns:

	error: nil or an error that occurred while processing the ping db
*/
func (pr *PostgresRepository) PingContext(ctx context.Context) error {
	if err := pr.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

/*
Close closes the open database connection

Returns:

	error: nil or an error that occurred while closing connection
*/
func (pr *PostgresRepository) Close() error {
	if err := pr.db.Close(); err != nil {
		return err
	}
	return nil
}

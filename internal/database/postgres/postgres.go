package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/itaraxa/turbo-waddle/internal/log"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresRepository struct {
	DB *sql.DB
}

/*
NewPostgresRepository creates instance of PostgresRepository

Args:

	ctx context.Context
	l log.Logger
	databaseURL: string for connection to databse, example: "postgres://username:password@localhost:5432/database_name"

Returns:

	db *sql.DB
	err error
*/
func NewPostgresRepository(ctx context.Context, l log.Logger, databaseURL string) (pr *PostgresRepository, err error) {
	l.Info("Open connection to database", "connection string", databaseURL)
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		l.Error("Openning connection to database error", "error", err)
		return nil, ErrOpenConnection
	}

	// check connetion to storage
	if err := db.PingContext(ctx); err != nil {
		l.Error("Checking connection to database error", "error", err)
		return nil, errors.Join(err, ErrOpenConnection)
	}

	ctxWithTimeout, cancelWithTimeout := context.WithTimeout(ctx, 5*time.Second)
	defer cancelWithTimeout()

	l.Info("Start database migrations")
	err = prepareTablesContext(ctxWithTimeout, db)
	if err != nil {
		l.Error("Database migration error", "error", err)
		return nil, ErrMigration
	}
	l.Info("Database is ready")

	return &PostgresRepository{DB: db}, nil
}

/*
PingContext check connection to db

Args:

	ctx context.Context

Returns:

	error: nil or an error that occurred while processing the ping db
*/
func (pr *PostgresRepository) PingContext(ctx context.Context) error {
	if err := pr.DB.PingContext(ctx); err != nil {
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
	if err := pr.DB.Close(); err != nil {
		return err
	}
	return nil
}

/*
NewTransaction init transaction and create function that apply or rollback changes
*/
func NewTransaction(ctx context.Context, txOpts *sql.TxOptions, db *sql.DB) (*sql.Tx, func(tx *sql.Tx), error) {
	tx, err := db.BeginTx(ctx, txOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot start transaction: %w", err)
	}
	txFinish := func(tx *sql.Tx) {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}
	return tx, txFinish, nil
}

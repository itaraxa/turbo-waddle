package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/itaraxa/turbo-waddle/internal/log"
)

func (pr *PostgresRepository) CheckOrderInDB(ctx context.Context, l log.Logger, tx *sql.Tx, order string) (login string, err error) {

	err = tx.QueryRowContext(ctx, "SELECT user_id FROM gophermart.orders WHERE order_id = $1;", order).Scan(&login)
	return
}

func (pr *PostgresRepository) AddOrder(ctx context.Context, l log.Logger, tx *sql.Tx, order string, login string) (err error) {
	_, err = tx.ExecContext(ctx, "INSERT INTO gophermart.orders (order_id, user_id, order_status, processed_at) VALUES ($1, (SELECT user_id FROM gophermart.users WHERE user_name = $2), $3, $4)",
		order, login, `NEW`, time.Now())
	return
}

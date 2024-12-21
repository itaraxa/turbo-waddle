package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/itaraxa/turbo-waddle/internal/log"
	"github.com/itaraxa/turbo-waddle/internal/models"
	"github.com/shopspring/decimal"
)

/*
order status
*/
const (
	ORDER_NEW        = `NEW`
	ORDER_REGISTERED = `REGISTERED`
	ORDER_PROCESSING = `PROCESSING`
	ORDER_INVALID    = `INVALID`
	ORDER_PROCESSED  = `PROCESSED`
)

func (pr *PostgresRepository) CheckOrderInDB(ctx context.Context, l log.Logger, tx *sql.Tx, order string) (login string, err error) {

	err = tx.QueryRowContext(ctx, "SELECT user_id FROM gophermart.orders WHERE order_id = $1;", order).Scan(&login)
	return
}

/*
AddOrder executes SQL-request for adding order into gophermart.orders table

Args:

	ctx context.Context
	l log.Logger
	tx *sql.Tx
	order string
	login string

Returns:

	err error
*/
func (pr *PostgresRepository) AddOrder(ctx context.Context, l log.Logger, tx *sql.Tx, order string, login string) (err error) {
	l.Debug("adding order into postgres db", "login", login, "order", order)
	startTime := time.Now()
	_, err = tx.ExecContext(ctx, "INSERT INTO gophermart.orders (order_id, user_id, order_status, processed_at) VALUES ($1, (SELECT user_id FROM gophermart.users WHERE user_name = $2), $3, $4)",
		order, login, ORDER_NEW, time.Now())
	if err != nil {
		l.Error("adding order into postgres db error", "error", err)
	}
	l.Debug("adding order into postgres db completed", "duration", time.Since(startTime))
	return
}

type OrderStatus struct {
	Order, Status string
}

/*
GetNotProcessedOrders selects nonprocessed orders from postgres db

Args:

	ctx context.Context
	l log.Logger

Returns:

	orders []OrderStatus
	err error
*/
func (pr *PostgresRepository) GetNotProcessedOrders(ctx context.Context, l log.Logger) (orders []OrderStatus, err error) {
	l.Debug("getting not processed orders from postgres db")
	startTime := time.Now()

	rows, err := pr.DB.QueryContext(ctx, "SELECT order_id, order_status FROM gophermart.orders WHERE order_status != $1;", ORDER_PROCESSED)
	if err != nil {
		l.Error("getting not processed orders error", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order, status string

		if err = rows.Scan(&order, &status); err != nil {
			l.Error("getting not processed orders error", "error", err)
			return nil, err
		}

		orders = append(orders, OrderStatus{Order: order, Status: status})
	}
	if err = rows.Err(); err != nil {
		l.Error("rows iteration error", "error", err)
		return nil, err
	}

	l.Debug("getting not processed orders from postgres db completed", "duration", time.Since(startTime))
	return
}

/*
UpdateOrder updates order information: status and accrual sun in storage

Args:

	ctx context.Context
	l log.Logger
	order string
	status string
	accrual decimal.Decimal

Returns:

	err error
*/
func (pr *PostgresRepository) UpdateOrder(ctx context.Context, l log.Logger, order string, status string, accrual decimal.Decimal) (err error) {
	l.Debug("updating order in postgres db")
	startTIme := time.Now()

	_, err = pr.DB.ExecContext(ctx, "UPDATE gophermart.orders SET order_status = $2, order_sum = $3 WHERE order_id = $1", order, status, accrual)
	if err != nil {
		l.Error("updating order in postgres db error", "error", err)
		return
	}
	l.Debug("updating order in postgres db complited", "duration", time.Since(startTIme))
	return
}

/*
GetOrders gets list of user orders

Args:

	ctx context.Context
	l log.Logger
	login string

Returns:

	orders []models.Order
	err error
*/
func (pr *PostgresRepository) GetOrders(ctx context.Context, l log.Logger, login string) (orders []models.Order, err error) {
	l.Debug("getting orders from postgres db", "login", login)
	startTime := time.Now()

	ctxWithTimeout, cancel := context.WithTimeout(ctx, sqlQueryTimeout)
	defer cancel()

	query := `
		SELECT order_id, order_status, order_sum, processed_at 
		FROM gophermart.orders 
		WHERE user_id = (SELECT user_id FROM gophermart.users WHERE user_name = $1) 
		ORDER BY processed_at DESC;
	`
	rows, err := pr.DB.QueryContext(ctxWithTimeout, query, login)
	if err != nil {
		l.Error("getting orders from postgres db error", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order, status string
		var accrual decimal.Decimal
		var processedTime time.Time

		if err = rows.Scan(&order, &status, &accrual, &processedTime); err != nil {
			l.Error("getting orders from postgres db error", "error", err)
			return nil, err
		}

		orders = append(orders, models.Order{
			Number:     order,
			Status:     status,
			Accrual:    accrual,
			UploadedAt: processedTime,
		})
	}

	if err = rows.Err(); err != nil {
		l.Error("rows iteration error", "error", err)
		return nil, err
	}

	l.Debug("getting orders from postgres db completed", "duration", time.Since(startTime), "login", login, "order number", len(orders))
	return orders, nil
}

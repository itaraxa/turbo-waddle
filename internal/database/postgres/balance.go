package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	e "github.com/itaraxa/turbo-waddle/internal/errors"
	"github.com/itaraxa/turbo-waddle/internal/log"
	"github.com/itaraxa/turbo-waddle/internal/models"
	"github.com/shopspring/decimal"
)

/*
GetBalance gets current balance and total withdrawn from database

Args:

	ctx context.Context
	l log.Logger
	login string

Returns:

	balance models.Balance
	err error
*/
func (pr *PostgresRepository) GetBalance(ctx context.Context, l log.Logger, login string) (balance models.Balance, err error) {
	l.Debug("getting balance from postgres db", "login", login)
	startTime := time.Now()

	ctxWithTimeout, cancel := context.WithTimeout(ctx, sqlQueryTimeout)
	defer cancel()

	query := `
		SELECT current_balance, withdrawn
		FROM gophermart.balances
		WHERE user_id = (SELECT user_id FROM gophermart.users WHERE user_name = $1);
	`
	var curBallance, withdrawn decimal.Decimal
	err = pr.DB.QueryRowContext(ctxWithTimeout, query, login).Scan(&curBallance, &withdrawn)
	if err != nil && err != sql.ErrNoRows {
		l.Error("getting balance from postgres db error", "error", err)
		return
	}
	if errors.Is(err, sql.ErrNoRows) {
		l.Error("getting balance from postgres db error", "error", e.ErrNoData)
		err = errors.Join(err, e.ErrNoData)
		return
	}
	balance.Current = curBallance
	balance.Withdrawn = withdrawn

	l.Debug("getting balance from postgres db completed", "duration", time.Since(startTime), "login", login, "balance", balance.String())
	return
}

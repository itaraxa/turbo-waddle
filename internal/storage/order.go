package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/itaraxa/turbo-waddle/internal/database/postgres"
	e "github.com/itaraxa/turbo-waddle/internal/errors"
	"github.com/itaraxa/turbo-waddle/internal/log"
)

/*
LoadOrder loads order number in database.
Plan:
 0. check user in db
 1. check order in db
 2. add order into db

Args:

	ctx context.Context
	l log.Logger
	login string
	order string

Returns:

	err error
*/
func (s *Storage) LoadOrder(ctx context.Context, l log.Logger, login string, order string) (err error) {
	tx, txFinish, err := postgres.NewTransaction(ctx, nil, s.DB)
	if err != nil {
		l.Error("init transaction error", "error", err)
		return errors.Join(postgres.ErrInitTransaction, err)
	}
	defer txFinish(tx)

	exist, _, err := s.PostgresRepository.CheckUser(ctx, l, tx, login)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		l.Error("checking user in database error", "login", login, "error", err)
		err = errors.Join(err, e.ErrInternalServerError)
		return
	}
	if !exist {
		l.Error("unknown user", "login", login)
		err = e.ErrUserNotFound
		return
	}

	loginInDB, err := s.PostgresRepository.CheckOrderInDB(ctx, l, tx, order)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		// database error
		l.Error("checking order in database error", "order", order, "error", err)
		err = errors.Join(err, e.ErrInternalServerError)
		return
	}
	if errors.Is(err, sql.ErrNoRows) {
		// order does not stored in db -> 'NEW' order
		err = s.PostgresRepository.AddOrder(ctx, l, tx, order, login)
		if err != nil {
			l.Error("adding order in database error", "order", order, "error", err)
			err = errors.Join(err, e.ErrInternalServerError)
			return
		}
	}

	if loginInDB != login {
		// order has been loaded by other user
		l.Error("order loaded by other user", "login", login, "login in db", loginInDB, "order", order)
		err = e.ErrOrderAlreadyUploadedOtherUser
		return
	} else {
		// order has been loaded user
		l.Info("order already uploaded by user", "login", login)
		err = e.ErrOrderAlreadyUploaded
		return
	}
}

package postgres

import (
	"context"
	"database/sql"
	"errors"

	e "github.com/itaraxa/turbo-waddle/internal/errors"
	"github.com/itaraxa/turbo-waddle/internal/log"
)

/*
AddUser adds new user into gophermart.users
>Tx start
 1. Chek if not exist
 2. Add user
 3. Check adding

>Tx finish

Args:

	ctx context.Context
	l log.Logger
	login string
	hash [32]byte
	salt []byte

Returns:

	err error
*/
func (pr *PostgresRepository) AddUser(ctx context.Context, l log.Logger, login string, hash []byte, salt []byte) (err error) {
	tx, txFinish, err := NewTransaction(ctx, nil, pr.DB)
	if err != nil {
		l.Error("init transaction error", "error", err)
		return errors.Join(ErrInitTransaction, err)
	}
	defer txFinish(tx)

	l.Info("check existing user in database", "login", login)
	exist, user_id, err := pr.CheckUser(ctx, l, tx, login)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		l.Error("getting data from database error", "error", err)
		return errors.Join(err, ErrAddUserQueryToDB)
	}
	if exist {
		l.Info("user already exists", "login", login, "user_id", user_id)
		return errors.Join(e.ErrLoginIsAlreadyUsed)
	}
	l.Info("add user into database", "login", login)
	_, err = tx.ExecContext(ctx, "INSERT INTO gophermart.users (user_name, password_hash, password_salt) VALUES ($1, $2, $3)", login, hash, salt)
	if err != nil {
		l.Info("adding user error", "error", err)
		return errors.Join(e.ErrInternalServerError, err)
	}
	err = tx.QueryRowContext(ctx, "SELECT user_id FROM gophermart.users WHERE user_name = $1;", login).Scan(&user_id)
	if err != nil {
		l.Error("checking user in database error", "error", err)
		return errors.Join(err, ErrAddUserQueryToDB)
	}

	l.Info("user added", "login", login, "user_id", user_id)

	return
}

/*
AddSession saves autorisation token into user_sessions table

Args:

	ctx context.Context
	l log.Logger
	login string
	token string

Returns:

	err error
*/
func (pr *PostgresRepository) AddSession(ctx context.Context, l log.Logger, login string, token string) (err error) {
	_, err = pr.DB.ExecContext(ctx, "INSERT INTO gophermart.user_sessions (user_id, token) VALUES ((SELECT user_id FROM gophermart.users WHERE user_name = $1), $2)", login, token)
	return
}

/*
GetUserHash returns password_salt and password_hash from users table

Args:

	ctx context.Context
	l log.Logger
	login string

Returns:

	salt, hash []byte
	err error
*/
func (pr *PostgresRepository) GetUserHash(ctx context.Context, l log.Logger, login string) (salt, hash []byte, err error) {
	var user_id sql.NullInt64
	err = pr.DB.QueryRowContext(ctx, "SELECT user_id, password_hash, password_salt FROM gophermart.users WHERE user_name = $1;", login).Scan(&user_id, &hash, &salt)
	if err != nil && err != sql.ErrNoRows {
		l.Error("error executing SQL-query", "error", err)
		err = errors.Join(err, e.ErrInternalServerError)
	}

	if !user_id.Valid {
		l.Error("user not found", "login", login)
		err = errors.Join(e.ErrUserNotFound)
		return
	}

	return
}

func (pr *PostgresRepository) CheckUser(ctx context.Context, l log.Logger, tx *sql.Tx, login string) (exist bool, id int64, err error) {
	l.Info("check existing user in database", "login", login)
	var user_id sql.NullInt64
	err = tx.QueryRowContext(ctx, "SELECT user_id FROM gophermart.users WHERE user_name = $1;", login).Scan(&user_id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		l.Error("getting data from database error", "error", err)
		return false, 0, errors.Join(err, ErrAddUserQueryToDB)
	}
	return true, user_id.Int64, nil
}

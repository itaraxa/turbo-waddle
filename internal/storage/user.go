package storage

import (
	"context"
	"encoding/hex"

	"github.com/itaraxa/turbo-waddle/internal/log"
)

/*
AddNewUser adds new user into storage and add new session

Args:

	ctx context.Context
	l log.Logger
	login string
	hash []byte
	salt []byte
	token string

Returns:

	err error
*/
func (s *Storage) AddNewUser(ctx context.Context, l log.Logger, login string, hash []byte, salt []byte, token string) (err error) {
	l.Info("add new user", "login", login, "token", token)
	err = s.PostgresRepository.AddUser(ctx, l, login, hash, salt)
	if err != nil {
		l.Error("adding new user error", "error", err)
		return
	}
	l.Info("user added", "login", login)

	err = s.PostgresRepository.AddSession(ctx, l, login, token)
	if err != nil {
		l.Error("adding new user error", "error", err)
		return
	}

	return
}

/*
GetUserHash returns salt and password hash from storage

Args:

	ctx context.Context
	l log.Logger
	login string

Returns:

	salt, hash []byte
	err error
*/
func (s *Storage) GetUserHash(ctx context.Context, l log.Logger, login string) (salt, hash []byte, err error) {
	l.Debug("getting password hash and salt of existing user", "login", login)
	salt, hash, err = s.PostgresRepository.GetUserHash(ctx, l, login)
	if err != nil {
		l.Error("getting password hash and salt error", "error", err)
		return
	}
	l.Debug("password hash and salt are getted", "login", login, "salt", hex.EncodeToString(salt), "hash", hex.EncodeToString(hash))
	return
}

/*
AddSession saves token for registered user

Args:

	ctx context.Context
	l log.Logger
	login string
	token string

Returns:

	err error
*/
func (s *Storage) AddSession(ctx context.Context, l log.Logger, login string, token string) (err error) {
	l.Debug("adding new session with token", "login", login, "token", token)
	err = s.PostgresRepository.AddSession(ctx, l, login, token)
	l.Debug("session added", "login", login, "token", token)
	return
}

func (s *Storage) LoginUser(ctx context.Context, l log.Logger, login string, password string, token string) (err error) {

	return
}

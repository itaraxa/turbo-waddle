package storage

import (
	"context"

	"github.com/itaraxa/turbo-waddle/internal/log"
)

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

func (s *Storage) LoginUser(ctx context.Context, l log.Logger, login string, password string, token string) (err error) {

	return
}

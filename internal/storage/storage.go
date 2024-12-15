package storage

import (
	"context"
	"errors"

	"github.com/itaraxa/turbo-waddle/internal/database/postgres"
	"github.com/itaraxa/turbo-waddle/internal/log"
)

type Storager interface {
	UserStorager
}

type UserStorager interface {
	AddNewUser(ctx context.Context, l log.Logger, login string, password string, token string) (err error)
	LoginUser(ctx context.Context, l log.Logger, login string, password string, token string) (err error)
}

type Storage struct {
	*postgres.PostgresRepository
	dsn string
}

func NewStorage(ctx context.Context, l log.Logger, dsn string) (*Storage, error) {
	l.Info("Start creating new storage", "database source name", dsn)
	pr, err := postgres.NewPostgresRepository(ctx, l, dsn)
	if err != err {
		l.Error("Creating new storage error", "error", err)
		return nil, errors.Join(ErrCreateNewStorage, err)
	}
	l.Info("Strorrage created", "database source name", dsn)

	return &Storage{
		pr,
		dsn,
	}, nil
}
